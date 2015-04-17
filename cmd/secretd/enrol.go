package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func generateAuthorizedKeys(db *sql.DB) (lines []string, err error) {
	rows, err := db.Query("SELECT name, ssh_key FROM principals WHERE provisioned = true AND ssh_key IS NOT NULL")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name, ssh_key string
		err = rows.Scan(&name, &ssh_key)
		if err != nil {
			log.Fatal(err)
		}
		forceCommand := fmt.Sprintf("/usr/bin/secret-shell --principal %s --ssh", name)
		line := fmt.Sprintf("no-agent-forwarding,no-pty,no-X11-forwarding,no-user-rc,no-port-forwarding,command=\"%s\" %s %s\n", forceCommand, ssh_key, name)
		lines = append(lines, line)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

func enrol(db *sql.DB, principal string, newPrincipal, key string) (err error) {
	if err = CheckAclNonHierarchical(db, principal, "enrol"); err != nil {
		return
	}

	rows, err := db.Query("SELECT provisioned FROM principals WHERE name = $1", newPrincipal)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer rows.Close()
	if rows.Next() {
		var provisioned bool
		err = rows.Scan(&provisioned)
		if err != nil {
			log.Fatal(err)
		}
		if provisioned {
			return errors.New("User already provisioned")
		}
		_, err = db.Exec("UPDATE principals SET ssh_key = $2, provisioned = true WHERE name = $1", newPrincipal, key)
		if err != nil {
			log.Fatal(err)
			return err
		}
	} else {
		// insert fresh entry
		_, err = db.Exec("INSERT INTO principals(name, ssh_key, provisioned) VALUES ($1, $2, true)", newPrincipal, key)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	return err
}
