# secretd

The goal of secretd is to store secrets that hosts need access to in a
safe and secure manner.  Examples of such secrets could be API keys,
RSA private keys, or HTTP basic auth passwords.  Hosts should only
have access to the secrets they need, not all secrets.

## design considerations/goals

secretd should be:

- secure
- safe
- easy to use and easy onboard new users
- auditable

## design

Each environment has a server (`secret-1`) that stores the
secrets. Authentication is via SSH using the host SSH key of the host
requesting secrets. Using an `authorized_keys` file, this is then
mapped to a CLI client which again talks to the daemon over a unix
socket, passing the host name as well as the query. Authorization is
done by the daemon which fetches secrets from a database, decrypts
them and returns them over the unix socket, the CLI utility reads it
back and returns it to the host.

This is secure by employing defense in depth and reusing well-known
authentication methods we use elsewhere.  In addition, limiting the
exposure through well-defined interfaces which are auditable and
automatically testable helps improve security.  The use of a single
host for all secrets further reduces the attack surfaces and makes it
so we only have to protect a single host, rather than a model where
secrets are stored in multiple places.

Secrets are protected by ACLs, which ensures principals only get
access to the secrets they should. A missing ACL means no access, so
it is fail safe.

Onboarding new users and hosts is done via an enrolment key.  This
reduces the amount of software needed on the client side to an SSH
client, which will generally be available already.

Since access to secrets is only through the daemon, it is easy to add
comprehensive logging which ensures we have good audit trails for when
secrets were added, updated, accessed and removed and by whom.

## socket protocol

The client connects to the server over a UNIX socket at `/run/secretd/secretd.sock`.
Messages are JSON formatted hashes.  Each include an action and zero
or more arguments.

### enrol

    {
		"action": "enrol",
		"principal": "cache-lcy1120",
		"key": "ssh-rsa …"
    }

Adds the user to the list of users.  If the user already exists,
returns an error.

### authorization

    {
		"action": "authorize",
		"principal": "cache-lcy1120"
    }

Return value:

    {
		"action": "authorize",
		"status": "ok"
    }

Authorizes the connection for the user.  Can only be used once (per
connection) and must be used at the start of the connection.

### storing secrets

    {
		"action": "secret.put",
		"key": ["a", "b", "c"],
		"value": "s3kr1t"
    }

Return value:

    {
		"action": "secret.put",
		"status": "ok"
    }

Stores `s3kr1t` under the key given by key.  Any intermediary nodes in
the tree are created.  Any ACLs need to be explicitly applied, by
default only the inherited ACLs apply to the node.

### retrieving secrets

    {
		"action": "secret.get",
		"key": ["a", "b", "c"],
    }

Return value:

    {
		"status": "ok",
		"action": "secret.get",
		"value": "foo"
    }

Returns the secret stored under the given key.  The string may be
escaped (according to JSON's escape rules).

### discovering structure

    {
		"action": "secret.list",
		"key": ["a", "b"],
    }

Returns a list of the key's immediate children keys.

### ACL management

#### updating ACLs

    {
		"action": "acl.set",
		"key": ["a", "b", "c"],
		"group": "ops",
		"permissions": [ "read" ]
    }

Replaces any existing ACLs for the particular group.  An empty
permissions set will delete the ACL for the group.

### retrieving ACLs

    {
		"action": "acl.get",
		"key": ["a", "b", "c"],
    }

Return value:

    {
		"groups": {
			"ops": [ "read", "write" ]
		},
    }

### group management

#### create group

    {
		"action": "group.create",
		"group": "ops"
    }

#### add member to group

    {
		"action": "group.member_add",
		"group": "ops",
		"member": "tfheen"
    }

#### remove member from group

    {
		"action": "group.member_remove",
		"group": "ops",
		"member": "tfheen"
    }

#### retrieve group members

    {
		"action": "group.member_list",
		"group": "ops"
    }

#### delete group

    {
		"action": "group.delete",
		"group": "ops"
    }

#### list groups

    {
		"action": "group.list",
    }

Return value:

    {
		"action": "group.list",
		"groups": [
			"ops",
			"caches"
		]
    }

Lists all groups

### principal management

#### create principal

    {
		"action": "principal.create",
		"principal": "foo",
		"key": "ssh-rsa …",
		"provisioned": true
    }

#### delete principal

    {
		"action": "principal.delete",
		"group": "foo"
    }

#### list principals

    {
		"action": "principal.list",
    }

Return value:

    {
		"action": "principal.list",
		"principals": [
			"foo",
			"bar"
		]
    }

Lists all principals.

## ACL primitives

ACLs are additive and positive, there is no way to grant A access to
/a/b, but not /a/b/c.

### read

The permission to read the value part of a secret

### write

The ability to update the value of a secret

### manage

The ability to give other principals rights to read, write or discover
a secret

### discover

The ability to find a secret and discover any leaf nodes.

### group_manage

The ability to create, update and remove groups.  The key is in this
case null.  For the initial implementation, this is a binary flag.

### principal_manage

The ability to create, update and remove principals.  The key is in
this case null.  For the initial implementation, this is a binary
flag.

### enrol

The ability to enrol new hosts. The key is in this case null. This is
a binary flag.

## built-in/magic groups

### all

all is the group that is magically populated by all principals.  It
can be granted permissions as any other group.

XXX: hard to implement.  Needed?

## what do we need to build?

✓ database schema
✓ socket protocol
✓ define ACLs

- logging
- cli utility
- server
✓ key enrolment
  ✓ generate authorized_keys file


## XXX Things to figure out

- naming of the client: What kerberos calls a principal, but the name
  principal confuses people.  Suggestions so far include client
  entity, entity, client, identity, grantee, and principal

- defense in depth if the authorized_keys generator is buggy.

- lock down ssh options (from, forcecommand)

- return messages, including errors

- audit logs, log to syslog or in db?

- message for setting encryption key

- user list/create/delete/show
