package tencentcloud_im

import (
	"context"
	"math/rand"

	"github.com/leapthinking/tencentcloud-im/consts"
	"github.com/leapthinking/tencentcloud-im/types"
)

type GroupAppDefinedDataItem struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type AppMemberDefinedDataItem struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type CreateGroupMemberItem struct {
	Member_Account       string                     `json:"Member_Account"`
	Role                 consts.GroupRole           `json:"Role,omitempty"`
	AppMemberDefinedData []AppMemberDefinedDataItem `json:"AppMemberDefinedData,omitempty"`
}

type CreateGroupRequest struct {
	// required
	Type consts.GroupType `json:"Type"`
	Name string           `json:"Name,omitempty"`
	// optional
	Owner_Account   string                    `json:"Owner_Account,omitempty"`
	GroupId         string                    `json:"GroupId,omitempty"`
	Introduction    string                    `json:"Introduction,omitempty"`
	Notification    string                    `json:"Notification,omitempty"`
	FaceUrl         string                    `json:"FaceUrl,omitempty"`
	MaxMemberCount  int                       `json:"MaxMemberCount,omitempty"`
	ApplyJoinOption consts.ApplyJoinOption    `json:"ApplyJoinOption,omitempty"`
	AppDefinedData  []GroupAppDefinedDataItem `json:"AppDefinedData,omitempty"`
	MemberList      []CreateGroupMemberItem   `json:"MemberList,omitempty"`
}

type CreateGroupResponse struct {
	IMResponse
	GroupId string `json:"GroupId"`
}

type CreateGroupOpt interface {
	ApplyToCreateGroupRequest(req *CreateGroupRequest)
}

func (c *Client) CreateGroup(ctx context.Context, name string, groupType consts.GroupType, opts ...CreateGroupOpt) *CreateGroupResponse {
	req := c.newRequest(ctx, Service_GROUP_OPEN_HTTP_SVC, Command_CREATE_GROUP)
	result := &CreateGroupResponse{}
	payload := CreateGroupRequest{Name: name, Type: groupType}
	for _, opt := range opts {
		opt.ApplyToCreateGroupRequest(&payload)
	}
	_, err := req.SetBody(payload).SetResult(result).Post(tencentCloudIMAPIEndpoint)
	if err != nil {
		result.internal = err
		return result
	}
	return result
}

type AddGroupMemberRequest struct {
	GroupId    string                `json:"GroupId"`
	Silence    int                   `json:"Silence,omitempty"`
	MemberList []types.MinimalMember `json:"MemberList"`
}

type AddGroupMemberOpt interface {
	ApplyToAddGroupMemberRequest(req *AddGroupMemberRequest)
}

type AddGroupMemberResultItem struct {
	Member_Account string `json:"Member_Account"`
	Result         int    `json:"Result"`
}

type AddGroupMemberResponse struct {
	IMResponse
	MemberList []AddGroupMemberResultItem `json:"MemberList"`
}

func (c *Client) AddGroupMember(ctx context.Context, groupId string, im_ids []string, opts ...AddGroupMemberOpt) *AddGroupMemberResponse {
	req := c.newRequest(ctx, Service_GROUP_OPEN_HTTP_SVC, Command_ADD_GROUP_MEMBER)
	members := make([]types.MinimalMember, len(im_ids))
	for index, im_id := range im_ids {
		members[index].MemberAccount = im_id
	}
	payload := AddGroupMemberRequest{GroupId: groupId, MemberList: members}
	for _, opt := range opts {
		opt.ApplyToAddGroupMemberRequest(&payload)
	}
	result := &AddGroupMemberResponse{}
	_, err := req.SetBody(payload).SetResult(&result).Post(tencentCloudIMAPIEndpoint)
	if err != nil {
		result.internal = err
		return result
	}
	return result
}

type GroupMsgGetSimpleRequest struct {
	GroupId      string `json:"GroupId"`
	ReqMsgNumber int    `json:"ReqMsgNumber"`
	ReqMsgSeq    *int   `json:"ReqMsgSeq,omitempty"`
}

type GroupMsgBody struct {
	MsgContent GroupMsgContent `json:"MsgContent"`
	MsgType    string          `json:"MsgType"`
}

type GroupMsgContent struct {
	Data string `json:"Data"`
	Desc string `json:"Desc"`
	Ext  string `json:"Ext"`
}

type GroupMessageItem struct {
	FromAccount  string         `json:"From_Account"`
	IsPlaceMsg   int            `json:"IsPlaceMsg"`
	MsgBody      []GroupMsgBody `json:"MsgBody"`
	MsgRandom    int            `json:"MsgRandom"`
	MsgSeq       int            `json:"MsgSeq"`
	MsgTimeStamp int            `json:"MsgTimeStamp"`
}

type GroupMsgGetSimpleResponse struct {
	IMResponse
	IsFinished int                `json:"IsFinished"`
	RspMsgList []GroupMessageItem `json:"RspMsgList"`
}

func (c *Client) GroupMsgGetSimple(ctx context.Context, groupId string, msgNumber int, seq *int) *GroupMsgGetSimpleResponse {
	req := c.newRequest(ctx, Service_GROUP_OPEN_HTTP_SVC, Command_GROUP_MSG_GET_SIMPLE)
	payload := GroupMsgGetSimpleRequest{GroupId: groupId, ReqMsgNumber: msgNumber, ReqMsgSeq: seq}
	result := &GroupMsgGetSimpleResponse{}
	_, err := req.SetBody(payload).SetResult(result).Post(tencentCloudIMAPIEndpoint)
	if err != nil {
		result.internal = err
		return result
	}
	return result
}

type GetGroupMemberInfoRequest struct {
	GroupId                          string   `json:"GroupId"`
	MemberInfoFilter                 []string `json:"MemberInfoFilter,omitempty"`
	MemberRoleFilter                 []string `json:"MemberRoleFilter,omitempty"`
	AppDefinedDataFilter_GroupMember []string `json:"AppDefinedDataFilter_GroupMember,omitempty"`
	Limit                            int      `json:"Limit,omitempty"`
	Offset                           int      `json:"Offset,omitempty"`
}

type GetGroupMemberInfoOpt interface {
	ApplyToGetGroupMemberInfoRequest(req *GetGroupMemberInfoRequest)
}

type GroupMemberItem struct {
	MemberAccount        string                    `json:"Member_Account"`
	Role                 string                    `json:"Role"`
	JoinTime             int                       `json:"JoinTime"`
	MsgSeq               int                       `json:"MsgSeq"`
	MsgFlag              string                    `json:"MsgFlag"`
	LastSendMsgTime      int                       `json:"LastSendMsgTime"`
	ShutUpUntil          int                       `json:"ShutUpUntil"`
	AppMemberDefinedData []GroupAppDefinedDataItem `json:"AppMemberDefinedData"`
}

type GetGroupMemberInfoResponse struct {
	IMResponse
	MemberNum  int               `json:"MemberNum"`
	MemberList []GroupMemberItem `json:"MemberList"`
}

func (c *Client) GetGroupMemberInfo(ctx context.Context, groupId string, opts ...GetGroupMemberInfoOpt) *GetGroupMemberInfoResponse {
	req := c.newRequest(ctx, Service_GROUP_OPEN_HTTP_SVC, Command_GET_GROUP_MEMBER_INFO)
	payload := GetGroupMemberInfoRequest{GroupId: groupId}
	for _, opt := range opts {
		opt.ApplyToGetGroupMemberInfoRequest(&payload)
	}
	result := &GetGroupMemberInfoResponse{}
	_, err := req.SetResult(result).SetBody(payload).Post(tencentCloudIMAPIEndpoint)
	if err != nil {
		result.internal = err
		return result
	}
	return result
}

type DeleteGroupMemberRequest struct {
	GroupID            string   `json:"GroupId"`
	Silence            int64    `json:"Silence"`
	Reason             string   `json:"Reason"`
	MemberToDelAccount []string `json:"MemberToDel_Account"`
}

type DeleteGroupMemberOpt interface {
	ApplyToDeleteGroupMemberRequest(request *DeleteGroupMemberRequest)
}

func (c *Client) DeleteGroupMember(ctx context.Context, groupId string, imIds []string, opts ...DeleteGroupMemberOpt) *IMResponse {
	req := c.newRequest(ctx, Service_GROUP_OPEN_HTTP_SVC, Command_DELETE_GROUP_MEMBER)
	payload := DeleteGroupMemberRequest{GroupID: groupId, MemberToDelAccount: imIds}
	for _, opt := range opts {
		opt.ApplyToDeleteGroupMemberRequest(&payload)
	}
	result := &IMResponse{}
	_, err := req.SetResult(result).SetBody(payload).Post(tencentCloudIMAPIEndpoint)
	if err != nil {
		result.internal = err
		return result
	}
	return result
}

type ModifyGroupMemberInfoRequest struct {
	GroupID              string                    `json:"GroupId"`
	MemberAccount        string                    `json:"Member_Account"`
	Role                 consts.GroupRole          `json:"Role,omitempty"`
	MsgFlag              string                    `json:"MsgFlag,omitempty"`
	NameCard             string                    `json:"NameCard,omitempty"`
	AppMemberDefinedData []GroupAppDefinedDataItem `json:"AppMemberDefinedData,omitempty"`
}

type ModifyGroupMemberInfoOpt interface {
	ApplyToModifyGroupMemberInfoRequest(request *ModifyGroupMemberInfoRequest)
}

func (c *Client) ModifyGroupMemberInfo(ctx context.Context, groupId string, imId string, opts ...ModifyGroupMemberInfoOpt) *IMResponse {
	req := c.newRequest(ctx, Service_GROUP_OPEN_HTTP_SVC, Command_MODIFY_GROUP_MEMBER_INFO)
	payload := ModifyGroupMemberInfoRequest{
		GroupID: groupId, MemberAccount: imId,
	}
	for _, opt := range opts {
		opt.ApplyToModifyGroupMemberInfoRequest(&payload)
	}
	result := &IMResponse{}
	c.sendRequest(req, payload, result)
	return result
}

type ModifyGroupBaseInfoRequest struct {
	GroupID         string                    `json:"GroupId"`
	Name            string                    `json:"Name,omitempty"`
	Introduction    string                    `json:"Introduction,omitempty"`
	Notification    string                    `json:"Notification,omitempty"`
	FaceURL         string                    `json:"FaceUrl,omitempty"`
	MaxMemberNum    int64                     `json:"MaxMemberNum,omitempty"`
	ApplyJoinOption consts.ApplyJoinOption    `json:"ApplyJoinOption,omitempty"`
	ShutUpAllMember string                    `json:"ShutUpAllMember,omitempty"`
	AppDefinedData  []GroupAppDefinedDataItem `json:"AppDefinedData,omitempty"`
}

type ModifyGroupBaseInfoOpt interface {
	ApplyModifyGroupBaseInfoRequest(req *ModifyGroupBaseInfoRequest)
}

func (c *Client) ModifyGroupBaseInfo(ctx context.Context, groupId string, opts ...ModifyGroupBaseInfoOpt) *IMResponse {
	req := c.newRequest(ctx, Service_GROUP_OPEN_HTTP_SVC, Command_MODIFY_GROUP_BASE_INFO)
	payload := ModifyGroupBaseInfoRequest{GroupID: groupId}
	for _, opt := range opts {
		opt.ApplyModifyGroupBaseInfoRequest(&payload)
	}
	result := &IMResponse{}
	c.sendRequest(req, payload, result)
	return result
}

type ChangeGroupOwnerRequest struct {
	GroupID         string `json:"GroupId"`
	NewOwnerAccount string `json:"NewOwner_Account"`
}

func (c *Client) ChangeGroupOwner(ctx context.Context, groupId string, imId string) *IMResponse {
	req := c.newRequest(ctx, Service_GROUP_OPEN_HTTP_SVC, Command_MODIFY_GROUP_BASE_INFO)
	payload := ChangeGroupOwnerRequest{GroupID: groupId, NewOwnerAccount: imId}
	result := &IMResponse{}
	c.sendRequest(req, payload, result)
	return result
}

type SendGroupMsgRequest struct {
	GroupID               string           `json:"GroupId"`
	Random                int64            `json:"Random"`
	MsgPriority           GroupMsgPriority `json:"MsgPriority,omitempty"`
	MsgBody               []MsgBody        `json:"MsgBody"`
	FromAccount           string           `json:"From_Account,omitempty"`
	OfflinePushInfo       *OfflinePushInfo `json:"OfflinePushInfo,omitempty"`
	ForbidCallbackControl []string         `json:"ForbidCallbackControl,omitempty"`
}

type SendGroupMsgRequestOpt interface {
	ApplyToSendGroupMsgRequest(req *SendGroupMsgRequest)
}

type sendGroupMsgRequestOptFunc func(req *SendGroupMsgRequest)

func (f sendGroupMsgRequestOptFunc) ApplyToSendGroupMsgRequest(req *SendGroupMsgRequest) {
	f(req)
}

const (
	forbidBeforeSendMsgCallback = "ForbidBeforeSendMsgCallback"
	forbidAfterSendMsgCallback  = "ForbidAfterSendMsgCallback"
)

func ForbidBeforeSendMsgCallback() SendGroupMsgRequestOpt {
	return sendGroupMsgRequestOptFunc(func(req *SendGroupMsgRequest) {
		for _, item := range req.ForbidCallbackControl {
			if item == forbidBeforeSendMsgCallback {
				return
			}
		}
		req.ForbidCallbackControl = append(req.ForbidCallbackControl, forbidBeforeSendMsgCallback)
	})
}

func ForbidAfterSendMsgCallback() SendGroupMsgRequestOpt {
	return sendGroupMsgRequestOptFunc(func(req *SendGroupMsgRequest) {
		for _, item := range req.ForbidCallbackControl {
			if item == forbidAfterSendMsgCallback {
				return
			}
		}
		req.ForbidCallbackControl = append(req.ForbidCallbackControl, forbidAfterSendMsgCallback)
	})
}

func SetFromAccount(from string) sendGroupMsgRequestOptFunc {
	return sendGroupMsgRequestOptFunc(func(req *SendGroupMsgRequest) {
		req.FromAccount = from
	})
}

type GroupMsgPriority string

const (
	GroupMsgPriorityHigh   GroupMsgPriority = "High"
	GroupMsgPriorityNormal                  = "Normal"
	GroupMsgPriorityLow                     = "Low"
	GroupMsgPriorityLowest                  = "Lowest"
)

type SendGroupMsgResponse struct {
	IMResponse
	MsgTime int64 `json:"MsgTime"`
	MsgSeq  int64 `json:"MsgSeq"`
}

func (c *Client) SendGroupMsg(ctx context.Context, groupId string, msgBody []MsgBody, opts ...SendGroupMsgRequestOpt) *SendGroupMsgResponse {
	req := c.newRequest(ctx, Service_GROUP_OPEN_HTTP_SVC, Command_SEND_GROUP_MSG)
	payload := SendGroupMsgRequest{
		GroupID: groupId, MsgBody: msgBody,
		Random: int64(rand.Intn(1<<31 - 1)),
	}
	for _, opt := range opts {
		opt.ApplyToSendGroupMsgRequest(&payload)
	}
	result := &SendGroupMsgResponse{}
	c.sendRequest(req, payload, result)
	return result
}
