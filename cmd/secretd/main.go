package main

import (
	"database/sql"
	"github.com/davecgh/go-spew/spew"
	model "github.com/fastly/secretd/model/message/server"
	message "github.com/fastly/secretd/model/message"
	_ "github.com/lib/pq"
	"log"
	"net"
	"strings"
	"errors"
)

func getSecret(db *sql.DB, key []string) (secret string, err error) {
	// XXX: the pq driver should just be taught how to do arrays..
	k := "{" + strings.Join(key, ",") + "}"
	rows, err := db.Query("SELECT value FROM secret_tree WHERE path = $1::text[]", k)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer rows.Close()
	if !rows.Next() {
		// XXX: add actual error
		return "", errors.New("Not found")
	}
	rows.Scan(&secret)
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return secret, err
}

func putSecret(db *sql.DB, key []string, secret string) (err error) {
	// XXX: the pq driver should just be taught how to do arrays..
	var id uint64;
	k := "{" + strings.Join(key, ",") + "}"
	rows, err := db.Query("SELECT path_create_missing_elements from path_create_missing_elements($1::text[])", k)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		// XXX: add actual error
		return errors.New("Error inserting path")
	}
	rows.Scan(&id)
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	rows, err = db.Query("UPDATE secrets SET value = $1 WHERE secret_id = $2", secret, id)
	if err != nil {
		log.Fatal(err)
		return err
	}
	rows.Close()
	return err
}


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
	case *message.AuthorizationMessage:
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
		reply := message.AuthorizationReplyMessage{Action: "authorize", Status: "ok"}
		model.SendReply(c, reply)
	default:
		model.SendReplySimpleError(c, "Missing authorization message")
		return
	}

	// Authorized as $principal, start next step of state machine
	spew.Dump(principal)
	for {
		gm, err := model.GetMessage(c)
		if err != nil {
			/* XXX: log */
			log.Printf("got %s, exiting loop\n", err)
			return
		}
		switch m := gm.(type) {
		case *message.AuthorizationMessage:
			model.SendReplySimpleError(c, "Unexpected authorization message")
			return
		case *message.SecretGetMessage:
			secret, err := getSecret(db, m.Key)
			if err != nil {
				/* XXX: secret not found */
				log.Printf("secret not found?\n", err)
				reply := message.SecretGetReplyMessage{Action: "secret.get", Status: "error", Reason: err.Error()}
				model.SendReply(c, reply)
				continue
			}
			resp := message.GenericReplyJSON{Status: "ok", Action: "secret.get", Value: secret}
			err = model.SendReply(c, resp)
		case *message.SecretPutMessage:
			// XXX: check ACL
			err := putSecret(db, m.Key, m.Value)
			if err != nil {
				/* XXX: secret not found */
				log.Printf("something went wrong: %s\n", err)
				reply := message.SecretPutReplyMessage{Action: "secret.put", Status: "error", Reason: err.Error()}
				model.SendReply(c, reply)
				continue
			}
			resp := message.GenericReplyJSON{Status: "ok", Action: "secret.put"}
			err = model.SendReply(c, resp)
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
