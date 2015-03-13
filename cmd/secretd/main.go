package main

import (
	_ "database/sql"
	"github.com/davecgh/go-spew/spew"
	"github.com/fastly/secretd/model"
	_ "github.com/lib/pq"
	"log"
	"net"
)

func secretServer(c net.Conn) {
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
		model.SendReplySimpleOK(c)
	default:
		model.SendReplySimpleError("Missing authorization message")
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

		go secretServer(fd)
	}
}

/*func main() {
	db, err := sql.Open("postgres", "user=secretd dbname=secrets host=/run/postgresql sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT 5*5")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var i int
		if err := rows.Scan(&i); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d\n", i)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
*/
