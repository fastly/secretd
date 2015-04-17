CREATE OR REPLACE FUNCTION path_create_missing_elements(path text[]) RETURNS integer AS $$
DECLARE
current_parent_id integer := NULL;
old_parent_id integer := NULL;
a_key text;
BEGIN
    FOREACH a_key IN ARRAY path LOOP
        RAISE NOTICE 'current_parent_id is %', current_parent_id;
        IF current_parent_id IS NULL THEN
            SELECT secret_id INTO current_parent_id FROM secret_tree WHERE parent IS NULL AND key = a_key;
        ELSE
	    SELECT secret_id INTO current_parent_id FROM secret_tree WHERE parent = current_parent_id AND key = a_key;
        END IF;
        IF NOT FOUND THEN
            RAISE NOTICE 'current_parent_id is %, inserting', old_parent_id;
            INSERT INTO secrets(parent, key) VALUES (old_parent_id, a_key) RETURNING secret_id INTO current_parent_id;
        END IF;
        old_parent_id := current_parent_id;
    END LOOP;
    RETURN current_parent_id;
END;
$$ LANGUAGE plpgsql;
