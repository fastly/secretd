package message

type GenericMessage interface {
}

type GenericMessageJSON struct {
	Action string `json:"action"`
}

type AuthorizationMessage struct {
	Action    string `json:"action"`
	Principal string `json:"principal"`
}

func NewAuthorizationMessage(principal string) AuthorizationMessage {
	m := AuthorizationMessage{Action: "authorize", Principal: principal}
	return m
}

type AuthorizationReplyMessage struct {
	Action string `json:"action"`
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

func NewAuthorizationReplyMessage(status string, reason string) AuthorizationReplyMessage {
	m := AuthorizationReplyMessage{Action: "authorize", Status: status, Reason: reason}
	return m
}

type SecretGetMessage struct {
	Action string   `json:"action"`
	Key    []string `json:"key"`
}

func NewSecretGetMessage(key []string) SecretGetMessage {
	m := SecretGetMessage{Action: "secret.get", Key: key}
	return m
}

type SecretGetReplyMessage struct {
	Action string `json:"action"`
	Status string `json:"status"`
	Value  string `json:"value"`
	Reason string `json:"reason,omitempty"`
}

func NewSecretGetReplyMessage(status string) SecretGetReplyMessage {
	m := SecretGetReplyMessage{Action: "secret.get", Status: status}
	return m
}

type SecretPutMessage struct {
	Action string   `json:"action"`
	Key    []string `json:"key"`
	Value  string   `json:"value"`
}

func NewSecretPutMessage(key []string, secret string) SecretPutMessage {
	m := SecretPutMessage{Action: "secret.put", Key: key, Value: secret}
	return m
}

type SecretPutReplyMessage struct {
	Action string `json:"action"`
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

func NewSecretPutReplyMessage(status string) SecretPutReplyMessage {
	m := SecretPutReplyMessage{Action: "secret.put", Status: status}
	return m
}

type GenericReply interface {
}

type GenericReplyJSON struct {
	Status string `json:"status"`
	Action string `json:"action"`
	Reason string `json:"reason,omitempty"`
	Value  string `json:"value,omitempty"`
}
