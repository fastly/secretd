package message

type AclGetMessage struct {
	Action string   `json:"action"`
	Key    []string `json:"key"`
}

func NewAclGetMessage(key []string) AclGetMessage {
	m := AclGetMessage{Action: "acl.get", Key: key}
	return m
}

type AclGetReplyMessage struct {
	Action string              `json:"action"`
	Status string              `json:"status"`
	Reason string              `json:"reason,omitempty"`
	Groups map[string][]string `json:"groups"`
}

func NewAclGetReplyMessage(status string, groups map[string][]string) AclGetReplyMessage {
	m := AclGetReplyMessage{Action: "acl.get", Status: status, Groups: groups}
	return m
}

type AclSetMessage struct {
	Action      string   `json:"action"`
	Key         []string `json:"key"`
	Group       string   `json:"group"`
	Permissions []string `json:"permissions"`
}

func NewAclSetMessage(key []string, group string, permissions []string) AclSetMessage {
	m := AclSetMessage{Action: "acl.set", Key: key, Group: group, Permissions: permissions}
	return m
}

type AclSetReplyMessage struct {
	Action string `json:"action"`
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

func NewAclSetReplyMessage(status string) AclSetReplyMessage {
	m := AclSetReplyMessage{Action: "acl.set", Status: status}
	return m
}
