package main

import (
	"database/sql"
	"log"
)

func listPrincipals(db *sql.DB, principal string) (principals []string, err error) {
	if err = CheckAclNonHierarchical(db, principal, "principal_manage"); err != nil {
		return
	}
	rows, err := db.Query("SELECT name FROM principals")
	if err != nil {
		log.Fatal(err)
		return
	}
	for rows.Next() {
		var principal string
		err = rows.Scan(&principal)
		if err != nil {
			log.Fatal(err)
		}
		principals = append(principals, principal)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

func createPrincipal(db *sql.DB, principal, newPrincipal, SSHKey string, provisioned bool) (err error) {
	if err = CheckAclNonHierarchical(db, principal, "principal_manage"); err != nil {
		return
	}

	_, err = db.Exec("INSERT INTO principals(name, ssh_key, provisioned) VALUES ($1, $2, $3)", newPrincipal, SSHKey, provisioned)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return err
}

func deletePrincipal(db *sql.DB, principal, deletePrincipal string) (err error) {
	if err = CheckAclNonHierarchical(db, principal, "principal_manage"); err != nil {
		return
	}

	_, err = db.Exec("DELETE FROM principals WHERE name = $1", deletePrincipal)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return err
}
