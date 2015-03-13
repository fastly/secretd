package model

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"io"
)

type GenericMessage interface {
	Action() string
	SetAction(string)
}

type GenericMessageJSON struct {
	Action string `json:"action"`
}

type AuthorizationMessage struct {
	Act       string `json:"action"`
	Principal string `json:"principal"`
}

func (m AuthorizationMessage) Action() string {
	return m.Act
}

func (m AuthorizationMessage) SetAction(s string) {
	m.Act = s
}

type SecretGetMessage struct {
	Act string   `json:"action"`
	Key []string `json:"key"`
}

func (m SecretGetMessage) Action() string {
	return m.Act
}

func (m SecretGetMessage) SetAction(s string) {
	m.Act = s
}

func GetMessage(r io.Reader) (m GenericMessage, err error) {
	var rawmessage json.RawMessage
	var message GenericMessageJSON
	dec := json.NewDecoder(r)
	err = dec.Decode(&rawmessage)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawmessage, &message)
	if err != nil {
		return nil, err
	}
	switch message.Action {
	case "authorize":
		m = new(AuthorizationMessage)
	case "secret.get":
		m = new(SecretGetMessage)
	default:
		panic("Unknown message type")
	}
	err = json.Unmarshal(rawmessage, &m)
	spew.Dump(rawmessage, message)
	return m, err
}

type GenericResponse interface {
}

type GenericResponseJSON struct {
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

func SendReply(w io.Writer, response GenericResponse) (err error) {
	enc := json.NewEncoder(w)
	spew.Dump(response)
	err = enc.Encode(response)
	if err != nil {
		return err
	}
	return nil
}

func SendReplySimpleStatus(w io.Writer, status string) (err error) {
	resp := GenericResponseJSON{Status: status}
	return SendReply(w, resp)
}

func SendReplySimpleOK(w io.Writer) (err error) {
	resp := GenericResponseJSON{Status: "ok"}
	spew.Dump(resp)
	return SendReply(w, resp)
}

func SendReplySimpleError(w io.Writer, reason string) (err error) {
	resp := GenericResponseJSON{Status: "error", Reason: reason}
	return SendReply(w, resp)
}
