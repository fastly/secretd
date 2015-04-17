package main

import (
	"database/sql"
	"log"
)

func listGroups(db *sql.DB, principal string) (groups []string, err error) {
	if err = CheckAclNonHierarchical(db, principal, "group_manage"); err != nil {
		return
	}

	rows, err := db.Query("SELECT name FROM groups")
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

func createGroup(db *sql.DB, principal, group string) (err error) {
	if err = CheckAclNonHierarchical(db, principal, "group_manage"); err != nil {
		return
	}

	_, err = db.Exec("INSERT INTO groups(name) VALUES ($1)", group)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return err
}

func deleteGroup(db *sql.DB, principal, group string) (err error) {
	if err = CheckAclNonHierarchical(db, principal, "group_manage"); err != nil {
		return
	}

	_, err = db.Exec("INSERT INTO groups(name) VALUES ($1)", group)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return err
}
