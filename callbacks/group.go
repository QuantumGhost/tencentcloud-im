package callbacks

import (
	"github.com/leapthinking/tencentcloud-im/consts"
	"github.com/leapthinking/tencentcloud-im/types"
)

type CallbackBeforeInviteJoinGroupPayload struct {
	CallbackCommand    string                `json:"CallbackCommand"`
	GroupID            string                `json:"GroupId"`
	Type               string                `json:"Type"`
	OperatorAccount    string                `json:"Operator_Account"`
	DestinationMembers []types.MinimalMember `json:"DestinationMembers"`
}

type CallbackBeforeInviteJoinGroupResponse struct {
	CallbackCommonResponse
	RefusedMembers_Account []string `json:"RefusedMembers_Account"`
}

type CallbackAfterNewMemberJoinPayload struct {
	CallbackCommand string                `json:"CallbackCommand"`
	GroupID         string                `json:"GroupId"`
	Type            string                `json:"Type"`
	JoinType        string                `json:"JoinType"`
	OperatorAccount string                `json:"Operator_Account"`
	NewMemberList   []types.MinimalMember `json:"NewMemberList"`
}

type CallbackAfterMemberExitPayload struct {
	CallbackCommand string                `json:"CallbackCommand"`
	GroupID         string                `json:"GroupId"`
	Type            string                `json:"Type"`
	ExitType        string                `json:"ExitType"`
	OperatorAccount string                `json:"Operator_Account"`
	ExitMemberList  []types.MinimalMember `json:"ExitMemberList"`
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

type CallbackBeforeCreateGroupPayload struct {
	CallbackCommand string                `json:"CallbackCommand"`
	OperatorAccount string                `json:"Operator_Account"`
	OwnerAccount    string                `json:"Owner_Account"`
	Type            consts.GroupType      `json:"Type"`
	Name            string                `json:"Name"`
	CreatedNum      int                   `json:"CreatedNum"`
	MemberList      []types.MinimalMember `json:"MemberList"`
}
