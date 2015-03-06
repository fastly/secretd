
CREATE TABLE secrets (
       secret_id serial PRIMARY KEY,
       parent int REFERENCES secrets,
       key text NOT NULL,
       value bytea
);

CREATE TABLE acl_types (
       acl_type_id serial PRIMARY KEY,
       name text
);

CREATE TABLE acls (
       acl_id serial PRIMARY KEY,
       secret_id int REFERENCES secrets ON DELETE CASCADE,
       group_id int REFERENCES groups ON DELETE CASCADE,
       acl_type_id int REFERENCES acl_types
);

CREATE TABLE groups (
       group_id serial PRIMARY KEY,
       name text
);

CREATE TABLE principals (
       principal_id serial PRIMARY KEY,
       name text,
       ssh_key text, -- XXX: add constraint on ok characters
       provisioned boolean -- needed? Key off ssh_key?
);

CREATE TABLE group_membership (
       group_id int REFERENCES groups ON DELETE CASCADE,
       principal_id int REFERENCES principals ON DELETE CASCADE,
       PRIMARY_KEY(group_id, principal_id)
);
