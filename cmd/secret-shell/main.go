package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/fastly/secretd/model/message/client"
	"github.com/fastly/secretd/model/message"
	"log"
	"net"
	"flag"
	"os"
	"bufio"
)

var principal string
var action string 

var flagvar int
func init() {
	flag.StringVar(&principal, "principal", "", "principal to authorize as")
	flag.StringVar(&action, "action", "", "action")
}

func main() {
	flag.Parse()

	if principal == "" || action == "" {
		flag.Usage()
		return
	}
	
	// XXX: make socket location configurable
	c, err := net.Dial("unix", "/tmp/secretd.sock")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	/* Authorize */
	// XXX: grab from command line flag
	authorizationMessage := message.NewAuthorizationMessage(principal)
	client.SendMessage(c, authorizationMessage)
	m, err := client.GetMessage(c)
	if err != nil {
		panic(err)
	}
	if m.(*message.AuthorizationReplyMessage).Status != "ok" {
		panic(m)
	}
	// XXX: Don't for loop, do a switch on command. for loop would be for reply.
	switch action {
	case "secret.get":
		client.SendMessage(c, message.NewSecretGetMessage(flag.Args()))
		m, err = client.GetMessage(c)
		if err != nil {
			panic(err)
		}
		m, ok := m.(*message.SecretGetReplyMessage)
		if !ok {
			spew.Dump(m, ok)
			panic("Type conversion failed")
		}
		if m.Status != "ok" {
			println(m.Reason)
			os.Exit(1)
		}
		println(m.Value)

	case "secret.put":
		secret, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			panic(err)
		}
		secret = secret[:len(secret)-1]
		client.SendMessage(c, message.NewSecretPutMessage(flag.Args(), secret))
		m, err = client.GetMessage(c)
		if err != nil {
			panic(err)
		}
		m, ok := m.(*message.SecretPutReplyMessage)
		if !ok {
			spew.Dump(m, ok)
			panic("Type conversion failed")
		}
		if m.Status != "ok" {
			println(m.Reason)
			os.Exit(1)
		}
		println("Secret updated")
	case "secret.list":
		client.SendMessage(c, message.NewSecretListMessage(flag.Args()))
		m, err = client.GetMessage(c)
		if err != nil {
			panic(err)
		}
		m, ok := m.(*message.SecretListReplyMessage)
		if !ok {
			spew.Dump(m, ok)
			panic("Type conversion failed")
		}
		if m.Status != "ok" {
			println(m.Reason)
			os.Exit(1)
		}
		for _,key := range m.Keys {
			println(key)
		}
	default:
		log.Fatal("Unknown action %s", action)
	}
}
