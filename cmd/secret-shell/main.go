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
			panic("Type conversion failed")
		}
		if m.Status != "ok" {
			println(m.Reason)
			os.Exit(1)
		}
		println("Group created")
	case "group.delete":
		client.SendMessage(c, message.NewGroupDeleteMessage(args[0]))
		m, err = client.GetMessage(c)
		if err != nil {
			panic(err)
		}
		m, ok := m.(*message.GroupDeleteReplyMessage)
		if !ok {
			panic("Type conversion failed")
		}
		if m.Status != "ok" {
			println(m.Reason)
			os.Exit(1)
		}
		println("Group deleted")
	case "group.member_list":
		client.SendMessage(c, message.NewGroupMemberListMessage(args[0]))
		m, err = client.GetMessage(c)
		if err != nil {
			panic(err)
		}
		m, ok := m.(*message.GroupMemberListReplyMessage)
		if !ok {
			panic("Type conversion failed")
		}
		if m.Status != "ok" {
			println(m.Reason)
			os.Exit(1)
		}
		for _, key := range m.Members {
			println(key)
		}
	case "group.member_add":
		// XXX: check arguments and give useful error messages
		client.SendMessage(c, message.NewGroupMemberAddMessage(args[0], args[1]))
		m, err = client.GetMessage(c)
		if err != nil {
			panic(err)
		}
		m, ok := m.(*message.GroupMemberAddReplyMessage)
		if !ok {
			panic("Type conversion failed")
		}
		if m.Status != "ok" {
			println(m.Reason)
			os.Exit(1)
		}
		println("Member added")
	case "group.member_remove":
		// XXX: check arguments and give useful error messages
		client.SendMessage(c, message.NewGroupMemberRemoveMessage(args[0], args[1]))
		m, err = client.GetMessage(c)
		if err != nil {
			panic(err)
		}
		m, ok := m.(*message.GroupMemberRemoveReplyMessage)
		if !ok {
			panic("Type conversion failed")
		}
		if m.Status != "ok" {
			println(m.Reason)
			os.Exit(1)
		}
		println("Member removed")

	case "principal.list":
		client.SendMessage(c, message.NewPrincipalListMessage())
		m, err = client.GetMessage(c)
		if err != nil {
			panic(err)
		}
		m, ok := m.(*message.PrincipalListReplyMessage)
		if !ok {
			panic("Type conversion failed")
		}
		if m.Status != "ok" {
			println(m.Reason)
			os.Exit(1)
		}
		for _, key := range m.Principals {
			println(key)
		}
	case "principal.create":
		newPrincipal := args[0]
		key := args[1]
		provisioned := false
		if runFromSSH { // Will have been split earlier, need to join type and key
			key = fmt.Sprintf("%s %s", args[1], args[2])
			if len(args) > 3 {
				provisioned, err = strconv.ParseBool(args[3])
				if err != nil {
					panic(err)
				}
			}
		} else {
			if len(args) > 2 {
				provisioned, err = strconv.ParseBool(args[2])
				if err != nil {
					panic(err)
				}
			}
		}

		client.SendMessage(c, message.NewPrincipalCreateMessage(newPrincipal, key, provisioned))
		m, err = client.GetMessage(c)
		if err != nil {
			panic(err)
		}
		m, ok := m.(*message.PrincipalCreateReplyMessage)
		if !ok {
			panic("Type conversion failed")
		}
		if m.Status != "ok" {
			println(m.Reason)
			os.Exit(1)
		}
		println("Principal created")
	case "principal.delete":
		client.SendMessage(c, message.NewPrincipalDeleteMessage(args[0]))
		m, err = client.GetMessage(c)
		if err != nil {
			panic(err)
		}
		m, ok := m.(*message.PrincipalDeleteReplyMessage)
		if !ok {
			panic("Type conversion failed")
		}
		if m.Status != "ok" {
			println(m.Reason)
			os.Exit(1)
		}
		println("Principal deleted")
	case "acl.get":
		client.SendMessage(c, message.NewAclGetMessage(args))
		m, err = client.GetMessage(c)
		if err != nil {
			panic(err)
		}
		m, ok := m.(*message.AclGetReplyMessage)
		if !ok {
			panic("Type conversion failed")
		}
		if m.Status != "ok" {
			println(m.Reason)
			os.Exit(1)
		}
		for k, v := range m.Groups {
			fmt.Printf("%v:\t\t%v\n", k, v)
		}

	case "acl.set":
		// We need to be slightly clever with arguments here, order is:
		// --action acl.set $group "read,write" a b c
		group := args[0]
		permissions := strings.Split(args[1], ",")
		for i, p := range permissions {
			permissions[i] = strings.TrimSpace(p)
		}
		msg := message.NewAclSetMessage(args[2:], group, permissions)
		spew.Dump(msg)
		client.SendMessage(c, msg)
		m, err = client.GetMessage(c)
		if err != nil {
			panic(err)
		}
		m, ok := m.(*message.AclSetReplyMessage)
		if !ok {
			spew.Dump(m, ok)
			panic("Type conversion failed")
		}
		if m.Status != "ok" {
			println(m.Reason)
			os.Exit(1)
		}
		println("Permissions set")
	case "enrol":
		newPrincipal := args[0]
		key := args[1]
		if runFromSSH { // Will have been split earlier, need to join type and key
			key = fmt.Sprintf("%s %s", args[1], args[2])
		}
		msg := message.NewEnrolMessage(newPrincipal, key)
		client.SendMessage(c, msg)
		m, err = client.GetMessage(c)
		if err != nil {
			panic(err)
		}
		m, ok := m.(*message.EnrolReplyMessage)
		if !ok {
			panic("Type conversion failed")
		}
		if m.Status != "ok" {
			println(m.Reason)
			os.Exit(1)
		}
		println("Permissions set")
		spew.Dump(m, ok)
	default:
		log.Fatal("Unknown action: ", action)
	}
}
