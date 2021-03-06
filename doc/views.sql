CREATE OR REPLACE VIEW secret_tree AS WITH RECURSIVE
search_graph(secret_id, parent, path, key, value, path_id) AS
(
SELECT  secret_id, parent, ARRAY[key], key, value, ARRAY[secret_id]
FROM    secrets
WHERE parent IS NULL
UNION ALL
SELECT  s.secret_id, s.parent, sg.path|| s.key, s.key, s.value, sg.path_id || s.secret_id
FROM    search_graph sg
JOIN    secrets s
ON      s.parent = sg.secret_id
)
SELECT  *
FROM    search_graph;

GRANT SELECT on secret_tree TO secretd;

CREATE OR REPLACE VIEW acl_tree AS
WITH RECURSIVE
acl_graph(secret_id, path, principal, acl_type) AS
(
SELECT  s.secret_id, ARRAY[s.key], p.name, acl_types.name
FROM    secrets s
JOIN    acls a
ON      a.secret_id = s.secret_id
JOIN    acl_types
ON      acl_types.acl_type_id = a.acl_type_id
JOIN    group_membership g
ON      a.group_id = g.group_id
JOIN    principals p
ON      p.principal_id = g.principal_id
WHERE parent IS NULL
UNION ALL
SELECT DISTINCT s.secret_id, ag.path || s.key, p.name, acl_types.name
FROM    acl_graph ag
JOIN    secrets s
ON      s.parent = ag.secret_id
JOIN    acls a
ON      a.secret_id = s.secret_id
JOIN    acl_types
ON      acl_types.acl_type_id = a.acl_type_id
JOIN    group_membership g
ON      a.group_id = g.group_id
JOIN    principals p
ON      p.principal_id = g.principal_id
)
SELECT  a.principal, a.acl_type, s.*
FROM    acl_graph a, secret_tree s
WHERE arraycontains(s.path, a.path);

GRANT SELECT on acl_tree TO secretd;

CREATE OR REPLACE VIEW acl_non_hierarchical AS
SELECT  p.name AS principal, acl_types.name AS acl_type
FROM    acls a
JOIN    acl_types
ON      acl_types.acl_type_id = a.acl_type_id
JOIN    group_membership g
ON      a.group_id = g.group_id
JOIN    principals p
ON      p.principal_id = g.principal_id
WHERE a.secret_id IS NULL;

GRANT SELECT on acl_non_hierarchical TO secretd;

CREATE OR REPLACE VIEW acl_group_tree AS
SELECT  g.name AS grp, acl_types.name AS acl_type, s.*
FROM    acls a
JOIN    acl_types USING (acl_type_id)
JOIN    groups g USING (group_id)
JOIN    secret_tree s ON arraycontains(s.path_id, ARRAY[a.secret_id]);

GRANT SELECT on acl_group_tree TO secretd;
