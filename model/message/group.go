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

type GroupDeleteMessage struct {
	Action string `json:"action"`
	Group  string `json:"group"`
}

func NewGroupDeleteMessage(group string) GroupDeleteMessage {
	m := GroupDeleteMessage{Action: "group.delete", Group: group}
	return m
}

type GroupDeleteReplyMessage struct {
	Action string `json:"action"`
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

func NewGroupDeleteReplyMessage(status string) GroupDeleteReplyMessage {
	m := GroupDeleteReplyMessage{Action: "group.delete", Status: status}
	return m
}

type GroupMemberListMessage struct {
	Action string `json:"action"`
	Group  string `json:"group"`
}

func NewGroupMemberListMessage(group string) GroupMemberListMessage {
	m := GroupMemberListMessage{Action: "group.member_list", Group: group}
	return m
}

type GroupMemberListReplyMessage struct {
	Action  string   `json:"action"`
	Status  string   `json:"status"`
	Reason  string   `json:"reason,omitempty"`
	Members []string `json:"members"`
}

func NewGroupMemberListReplyMessage(status string, members []string) GroupMemberListReplyMessage {
	m := GroupMemberListReplyMessage{Action: "group.member_list", Status: status, Members: members}
	return m
}

type GroupMemberAddMessage struct {
	Action    string `json:"action"`
	Group     string `json:"group"`
	Principal string `json:"principal"`
}

func NewGroupMemberAddMessage(group, principal string) GroupMemberAddMessage {
	m := GroupMemberAddMessage{Action: "group.member_add", Group: group, Principal: principal}
	return m
}

type GroupMemberAddReplyMessage struct {
	Action string `json:"action"`
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

func NewGroupMemberAddReplyMessage(status string) GroupMemberAddReplyMessage {
	m := GroupMemberAddReplyMessage{Action: "group.member_add", Status: status}
	return m
}

type GroupMemberRemoveMessage struct {
	Action    string `json:"action"`
	Group     string `json:"group"`
	Principal string `json:"principal"`
}

func NewGroupMemberRemoveMessage(group, principal string) GroupMemberRemoveMessage {
	m := GroupMemberRemoveMessage{Action: "group.member_remove", Group: group, Principal: principal}
	return m
}

type GroupMemberRemoveReplyMessage struct {
	Action string `json:"action"`
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

func NewGroupMemberRemoveReplyMessage(status string) GroupMemberRemoveReplyMessage {
	m := GroupMemberRemoveReplyMessage{Action: "group.member_remove", Status: status}
	return m
}
