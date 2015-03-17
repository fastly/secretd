CREATE OR REPLACE VIEW secret_tree AS WITH RECURSIVE
search_graph(secret_id, path, value) AS
(
SELECT  secret_id, ARRAY[key], value
FROM    secrets
WHERE parent IS NULL
UNION ALL
SELECT  s.secret_id, sg.path|| s.key, s.value
FROM    search_graph sg
JOIN    secrets s
ON      s.parent = sg.secret_id
)
SELECT  *
FROM    search_graph;

GRANT SELECT on secret_tree TO secretd;
