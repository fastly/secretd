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
