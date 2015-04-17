CREATE TABLE secrets (
       secret_id serial PRIMARY KEY,
       parent int REFERENCES secrets,
       key text NOT NULL,
       value bytea
);
GRANT SELECT, INSERT, UPDATE, DELETE ON secrets to secretd;

CREATE TABLE acl_types (
       acl_type_id serial PRIMARY KEY,
       name text
);
GRANT SELECT, INSERT, UPDATE, DELETE ON acl_types to secretd;

CREATE TABLE groups (
       group_id serial PRIMARY KEY,
       name text
);
GRANT SELECT, INSERT, UPDATE, DELETE ON groups to secretd;

CREATE TABLE acls (
       acl_id serial PRIMARY KEY,
       secret_id int REFERENCES secrets ON DELETE CASCADE,
       group_id int NOT NULL REFERENCES groups ON DELETE CASCADE,
       acl_type_id int NOT NULL REFERENCES acl_types,
       UNIQUE (secret_id, group_id, acl_type_id)
);
GRANT SELECT, INSERT, UPDATE, DELETE ON acls to secretd;

CREATE TABLE principals (
       principal_id serial PRIMARY KEY,
       name text,  -- XXX: add constraint on ok characters
       ssh_key text, -- XXX: add constraint on ok characters
       provisioned boolean -- needed? Key off ssh_key?
);
GRANT SELECT, INSERT, UPDATE, DELETE ON principals to secretd;

CREATE TABLE group_membership (
       group_id int NOT NULL REFERENCES groups ON DELETE CASCADE,
       principal_id int NOT NULL REFERENCES principals ON DELETE CASCADE,
       PRIMARY KEY(group_id, principal_id)
);
GRANT SELECT, INSERT, UPDATE, DELETE ON group_membership to secretd;

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public to secretd;
