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

	_, err = db.Exec("DELETE FROM groups WHERE name = $1", group)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return err
}

func groupMemberList(db *sql.DB, principal, group string) (members []string, err error) {
	if err = CheckAclNonHierarchical(db, principal, "group_manage"); err != nil {
		return
	}

	rows, err := db.Query("SELECT principals.name FROM groups JOIN group_membership USING (group_id) JOIN principals USING (principal_id) WHERE groups.name = $1", group)
	if err != nil {
		log.Fatal(err)
		return []string{}, err
	}
	for rows.Next() {
		var member string
		err = rows.Scan(&member)
		if err != nil {
			log.Fatal(err)
		}
		members = append(members, member)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return members, err
}

func groupMemberAdd(db *sql.DB, authPrincipal, group, principal string) (err error) {
	if err = CheckAclNonHierarchical(db, principal, "group_manage"); err != nil {
		return
	}

	_, err = db.Exec("INSERT INTO group_membership(group_id, principal_id) VALUES ((SELECT group_id FROM groups WHERE name = $1), (SELECT principal_id FROM principals WHERE name = $2))", group, principal)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return err
}

func groupMemberRemove(db *sql.DB, authPrincipal, group, principal string) (err error) {
	if err = CheckAclNonHierarchical(db, principal, "group_manage"); err != nil {
		return
	}

	_, err = db.Exec("DELETE FROM group_membership WHERE group_id = (SELECT group_id FROM groups WHERE name = $1) AND principal_id = (SELECT principal_id FROM principals WHERE name = $2)", group, principal)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return err
}
