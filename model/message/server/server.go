package server

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	message "github.com/fastly/secretd/model/message"
	"io"
)

func GetMessage(r io.Reader) (ret message.GenericMessage, err error) {
	var rawmessage json.RawMessage
	var m message.GenericMessageJSON

	if err = json.NewDecoder(r).Decode(&rawmessage); err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawmessage, &m)
	if err != nil {
		return nil, err
	}
	switch m.Action {
	case "authorize":
		ret = new(message.AuthorizationMessage)
	case "secret.get":
		ret = new(message.SecretGetMessage)
	case "secret.put":
		ret = new(message.SecretPutMessage)
	case "secret.list":
		ret = new(message.SecretListMessage)
	case "group.list":
		ret = new(message.GroupListMessage)
	case "group.create":
		ret = new(message.GroupCreateMessage)
	default:
		// XXX: handle this more gracefully
		panic("Unknown message type")
	}
	err = json.Unmarshal(rawmessage, &ret)
	spew.Dump(rawmessage, m)
	return ret, err
}

func SendReply(w io.Writer, reply message.GenericReply) (err error) {
	enc := json.NewEncoder(w)
	spew.Dump(reply)
	err = enc.Encode(reply)
	if err != nil {
		return err
	}
	return nil
}

func SendReplySimpleStatus(w io.Writer, status string) (err error) {
	reply := message.GenericReplyJSON{Status: status}
	return SendReply(w, reply)
}

func SendReplySimpleOK(w io.Writer) (err error) {
	reply := message.GenericReplyJSON{Status: "ok"}
	spew.Dump(reply)
	return SendReply(w, reply)
}

func SendReplySimpleError(w io.Writer, reason string) (err error) {
	reply := message.GenericReplyJSON{Status: "error", Reason: reason}
	return SendReply(w, reply)
}
