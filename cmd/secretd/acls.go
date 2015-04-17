package main

import (
	"database/sql"
	"errors"
	"github.com/davecgh/go-spew/spew"
	"log"
)

func CheckAclNonHierarchical(db *sql.DB, principal, aclType string) (err error) {
	rows, err := db.Query("SELECT * FROM acl_non_hierarchical WHERE principal = $1 AND acl_type = $2", principal, aclType)
	if err != nil {
		// XXX: better logging
		log.Fatal(err)
		return errors.New("Permission denied")
	}
	defer rows.Close()
	if !rows.Next() {
		return errors.New("Permission denied")
	}
	return nil
}

func aclGet(db *sql.DB, principal string, key []string) (groups map[string][]string, err error) {
	// XXX: handle non-hierarchical ACLs too
	k := dbArrayToString(key)
	rows, err := db.Query("SELECT 1 FROM acl_tree WHERE principal = $1 AND acl_type = 'manage' AND path = $2::text[]", principal, k)
	if err != nil {
		log.Fatal(err)
		return groups, err
	}
	defer rows.Close()
	if !rows.Next() {
		// XXX: add actual error
		return groups, errors.New("Not found or permission denied")
	}
	rows, err = db.Query("SELECT grp, acl_type FROM acl_group_tree WHERE path = $1::text[]", k)
	if err != nil {
		log.Fatal(err)
		return groups, err
	}
	defer rows.Close()

	groups = make(map[string][]string)
	for rows.Next() {
		var group string
		var aclType string
		err = rows.Scan(&group, &aclType)
		if err != nil {
			log.Fatal(err)
		}
		g, ok := groups[group]
		if !ok {
			g = make([]string, 0)
			groups[group] = g
		}
		groups[group] = append(g, aclType)
	}
	return groups, err
}

func aclSet(db *sql.DB, principal string, key []string, group string, permissions []string) (err error) {
	// XXX: handle non-hierarchical ACLs too

	var secretId int
	k := dbArrayToString(key)
	spew.Dump(key, k)
	rows, err := db.Query("SELECT secret_id FROM acl_tree WHERE principal = $1 AND acl_type = 'manage' AND path = $2::text[]", principal, k)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		// XXX: add actual error
		return errors.New("Not found or permission denied")
	}

	err = rows.Scan(&secretId)
	if err != nil {
		log.Fatal(err)
	}

	validPermissions := map[string]int{}
	rows, err = db.Query("SELECT name, acl_type_id FROM acl_types")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var n string
		var i int
		err = rows.Scan(&n, &i)
		if err != nil {
			log.Fatal(err)
		}
		validPermissions[n] = i
	}

	for _, permission := range permissions {
		// XXX: make sure the permission is valid, else return error, before making any changes

		rows, err := db.Query("SELECT 1 FROM acls WHERE secret_id = $1 AND group_id = (SELECT group_id FROM groups WHERE name = $2) AND acl_type_id = $3", secretId, group, validPermissions[permission])
		if err != nil {
			log.Fatal(err)
			return err
		}
		defer rows.Close()
		if rows.Next() {
			continue
		}

		_, err = db.Exec("INSERT INTO acls(secret_id, group_id, acl_type_id) SELECT $1, (SELECT group_id FROM groups WHERE name = $2), $3)", secretId, group, validPermissions[permission])
		if err != nil {
			log.Fatal(err)
			return err
		}
		defer rows.Close()
	}
	return err
}
