package callbacks

type DestinationMemberItem struct {
	MemberAccount string `json:"Member_Account"`
}

type CallbackBeforeInviteJoinGroupPayload struct {
	CallbackCommand    string                  `json:"CallbackCommand"`
	GroupID            string                  `json:"GroupId"`
	Type               string                  `json:"Type"`
	OperatorAccount    string                  `json:"Operator_Account"`
	DestinationMembers []DestinationMemberItem `json:"DestinationMembers"`
}

type CallbackBeforeInviteJoinGroupResponse struct {
	CallbackCommonResponse
	RefusedMembers_Account []string `json:"RefusedMembers_Account"`
}

type NewMemberListItem struct {
	MemberAccount string `json:"Member_Account"`
}

type CallbackAfterNewMemberJoinPayload struct {
	CallbackCommand string              `json:"CallbackCommand"`
	GroupID         string              `json:"GroupId"`
	Type            string              `json:"Type"`
	JoinType        string              `json:"JoinType"`
	OperatorAccount string              `json:"Operator_Account"`
	NewMemberList   []NewMemberListItem `json:"NewMemberList"`
}

type ExitMemberListItem struct {
	MemberAccount string `json:"Member_Account"`
}

type CallbackAfterMemberExitPayload struct {
	CallbackCommand string               `json:"CallbackCommand"`
	GroupID         string               `json:"GroupId"`
	Type            string               `json:"Type"`
	ExitType        string               `json:"ExitType"`
	OperatorAccount string               `json:"Operator_Account"`
	ExitMemberList  []ExitMemberListItem `json:"ExitMemberList"`
}

type UserDefinedDataListItem struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type CallbackAfterGroupInfoChangedPayload struct {
	CallbackCommand     string                    `json:"CallbackCommand"`
	GroupID             string                    `json:"GroupId"`
	Type                string                    `json:"Type"`
	OperatorAccount     string                    `json:"Operator_Account"`
	Name                string                    `json:"Name"`
	Introduction        string                    `json:"Introduction"`
	Notification        string                    `json:"Notification"`
	FaceURL             string                    `json:"FaceUrl"`
	UserDefinedDataList []UserDefinedDataListItem `json:"UserDefinedDataList"`
}
