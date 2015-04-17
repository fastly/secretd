package client

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	message "github.com/fastly/secretd/model/message"
	"io"
)

func SendMessage(w io.Writer, m message.GenericMessage) (err error) {
	enc := json.NewEncoder(w)
	err = enc.Encode(m)
	if err != nil {
		return err
	}
	return nil
}

func GetMessage(r io.Reader) (m message.GenericMessage, err error) {
	var rawmessage json.RawMessage
	var gm message.GenericMessageJSON

	if err = json.NewDecoder(r).Decode(&rawmessage); err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawmessage, &gm)
	if err != nil {
		return nil, err
	}
	switch gm.Action {
	case "authorize":
		m = new(message.AuthorizationReplyMessage)
	case "secret.get":
		m = new(message.SecretGetReplyMessage)
	case "secret.put":
		m = new(message.SecretPutReplyMessage)
	case "secret.list":
		m = new(message.SecretListReplyMessage)
	case "group.list":
		m = new(message.GroupListReplyMessage)
	case "group.create":
		m = new(message.GroupCreateReplyMessage)
	case "group.delete":
		m = new(message.GroupDeleteReplyMessage)
	case "group.member_list":
		m = new(message.GroupMemberListReplyMessage)
	case "group.member_add":
		m = new(message.GroupMemberAddReplyMessage)
	case "group.member_remove":
		m = new(message.GroupMemberRemoveReplyMessage)
	case "acl.get":
		m = new(message.AclGetReplyMessage)
	case "acl.set":
		m = new(message.AclSetReplyMessage)
	case "enrol":
		m = new(message.EnrolReplyMessage)
	case "principal.list":
		m = new(message.PrincipalListReplyMessage)
	case "principal.create":
		m = new(message.PrincipalCreateReplyMessage)
	case "principal.delete":
		m = new(message.PrincipalDeleteReplyMessage)
	default:
		spew.Dump(rawmessage, gm)
		panic("Unknown message type")
	}
	err = json.Unmarshal(rawmessage, m)
	return m, err
}
