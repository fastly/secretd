package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/fastly/secretd/model/message"
	"github.com/fastly/secretd/model/message/client"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var principal string
var action string
var runFromSSH bool

var flagvar int

func init() {
	flag.StringVar(&principal, "principal", "", "principal to authorize as")
	flag.StringVar(&action, "action", "", "action")
	flag.BoolVar(&runFromSSH, "ssh", false, "look to SSH_ORIGINAL_COMMAND for action")
}

// parseOriginalCommand, unsurprisingly parses the
// SSH_ORIGINAL_COMMAND setting. It does that via a regex, since SSH
// provides absolutely no useful help for parsing it ourselves.  This
// means some crazy key names won't work.
//
// XXX: tests
func parseOriginalCommand() (action string, args []string, err error) {
	r, err := regexp.Compile(`secret-shell --action ([\w.-]+)\s+((?:[\w-,+/=]+\s*)+)?$`)
	if err != nil {
		panic(err)
	}
	m := r.FindStringSubmatch(os.Getenv("SSH_ORIGINAL_COMMAND"))
	if m == nil {
		return action, args, errors.New("Malformed command")
	}
	action = m[1]
	args = strings.Split(m[2], " ")
	return
}

func main() {
	flag.Parse()

	if principal == "" || (action == "" && !runFromSSH) {
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
	authorizationMessage := message.NewAuthorizationMessage(principal)
	client.SendMessage(c, authorizationMessage)
	m, err := client.GetMessage(c)
	if err != nil {
		panic(err)
	}
	if m.(*message.AuthorizationReplyMessage).Status != "ok" {
		panic(m)
	}
	args := flag.Args()
	if runFromSSH {
		// Parse action from
		action, args, err = parseOriginalCommand()
		if err != nil {
			log.Fatal(err)
		}
	}

	// XXX: Don't for loop, do a switch on command. for loop would be for reply.
	switch action {
	case "secret.get":
		client.SendMessage(c, message.NewSecretGetMessage(args))
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
		client.SendMessage(c, message.NewSecretPutMessage(args, secret))
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
		client.SendMessage(c, message.NewSecretListMessage(args))
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
		for _, key := range m.Keys {
			println(key)
		}
	case "group.list":
		client.SendMessage(c, message.NewGroupListMessage())
		m, err = client.GetMessage(c)
		if err != nil {
			panic(err)
		}
		m, ok := m.(*message.GroupListReplyMessage)
		if !ok {
			spew.Dump(m, ok)
			panic("Type conversion failed")
		}
		if m.Status != "ok" {
			println(m.Reason)
			os.Exit(1)
		}
		for _, key := range m.Groups {
			println(key)
		}
	case "group.create":
		client.SendMessage(c, message.NewGroupCreateMessage(args[0]))
		m, err = client.GetMessage(c)
		if err != nil {
			panic(err)
		}
		m, ok := m.(*message.GroupCreateReplyMessage)
		if !ok {
			spew.Dump(m, ok)
			panic("Type conversion failed")
		}
		if m.Status != "ok" {
			println(m.Reason)
			os.Exit(1)
		}
		println("Group created")
	default:
		log.Fatal("Unknown action %s", action)
	}
}
