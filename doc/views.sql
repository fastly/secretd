CREATE OR REPLACE VIEW secret_tree AS WITH RECURSIVE
search_graph(secret_id, parent, path, key, value) AS
(
SELECT  secret_id, parent, ARRAY[key], key, value
FROM    secrets
WHERE parent IS NULL
UNION ALL
SELECT  s.secret_id, s.parent, sg.path|| s.key, s.key, s.value
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
SELECT  secrets.secret_id, ARRAY[secrets.key], NULL, NULL
FROM    secrets
WHERE parent IS NULL
UNION ALL
SELECT  s.secret_id, ag.path || s.key, p.name, acl_types.name
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
WHERE principal IS NOT NULL AND arraycontains(s.path, a.path);

GRANT SELECT on acl_tree TO secretd;
