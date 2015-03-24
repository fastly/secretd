package message

type GroupListMessage struct {
	Action string `json:"action"`
}

func NewGroupListMessage() GroupListMessage {
	m := GroupListMessage{Action: "group.list"}
	return m
}

type GroupListReplyMessage struct {
	Action string   `json:"action"`
	Status string   `json:"status"`
	Reason string   `json:"reason,omitempty"`
	Groups []string `json:"groups"`
}

func NewGroupListReplyMessage(status string, groups []string) GroupListReplyMessage {
	m := GroupListReplyMessage{Action: "group.list", Status: status, Groups: groups}
	return m
}

type GroupCreateMessage struct {
	Action string `json:"action"`
	Group  string `json:"group"`
}

func NewGroupCreateMessage(group string) GroupCreateMessage {
	m := GroupCreateMessage{Action: "group.create", Group: group}
	return m
}

type GroupCreateReplyMessage struct {
	Action string   `json:"action"`
	Status string   `json:"status"`
	Reason string   `json:"reason,omitempty"`
	Groups []string `json:"groups"`
}

func NewGroupCreateReplyMessage(status string) GroupCreateReplyMessage {
	m := GroupCreateReplyMessage{Action: "group.create", Status: status}
	return m
}
