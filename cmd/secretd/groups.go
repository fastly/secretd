package main

import (
	"database/sql"
	"errors"
	"log"
)

func listGroups(db *sql.DB, principal string) (groups []string, err error) {
	rows, err := db.Query("SELECT * FROM acl_non_hierarchical WHERE principal = $1 AND acl_type = 'group_manage'", principal)
	if err != nil {
		log.Fatal(err)
		return []string{}, err
	}
	defer rows.Close()
	if !rows.Next() {
		return []string{}, errors.New("Permission denied")
	}

	rows, err = db.Query("SELECT name FROM groups")
	if err != nil {
		log.Fatal(err)
		return []string{}, err
	}
	for rows.Next() {
		var group string
		err = rows.Scan(&group)
		if err != nil {
			log.Fatal(err)
		}
		groups = append(groups, group)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return groups, err
}
