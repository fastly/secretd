package main

import (
	"database/sql"
	"github.com/davecgh/go-spew/spew"
	"github.com/fastly/secretd/model"
	_ "github.com/lib/pq"
	"log"
	"net"
)

func secretServer(c net.Conn, db *sql.DB) {
	/* state machine layout:
	   - enrol-then-terminate OR
	   - authorize
	   loop {
	     read operation
	     verify ACL
	     handle operation
	   }
	   exit
	*/
	var principal string
	/* Do authorization first */
	authMessage, err := model.GetMessage(c)
	spew.Dump(authMessage, err)
	if err != nil {
		return
	}
	/* State machine */
	spew.Dump(authMessage)
	switch m := authMessage.(type) {
	case *model.AuthorizationMessage:
		principal = m.Principal
		rows, err := db.Query("SELECT name FROM principals WHERE name = $1 AND provisioned = true", principal)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		if !rows.Next() {
			model.SendReplySimpleError(c, "No such principal")
			return
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		model.SendReplySimpleOK(c)
	default:
		model.SendReplySimpleError(c, "Missing authorization message")
		return
	}

	// Authorized as $principal, start next step of state machine
	spew.Dump(principal)
	for {
		message, err := model.GetMessage(c)
		if err != nil {
			/* XXX: log */
			log.Printf("got %s, exiting loop\n", err)
			return
		}
		switch m := message.(type) {
		case *model.AuthorizationMessage:
			model.SendReplySimpleError(c, "Unexpected authorization message")
			return
		case *model.SecretGetMessage:
			spew.Dump(m)
		default:
			panic("Unknown message:")
			spew.Dump(m)
		}
	}
}

func main() {
	// XXX: make socket location configurable
	l, err := net.Listen("unix", "/tmp/secretd.sock")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	for {
		fd, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}

		/* XXX: make connection string configurable */
		db, err := sql.Open("postgres", "user=secretd dbname=secrets host=/run/postgresql sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		go secretServer(fd, db)
	}
}
