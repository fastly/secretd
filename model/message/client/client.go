package client

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"io"
	message "github.com/fastly/secretd/model/message"
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
	default:
		spew.Dump(rawmessage, gm)
		panic("Unknown message type")
	}
	err = json.Unmarshal(rawmessage, m)
	return m, err
}
