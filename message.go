package tencentcloud_im

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"math/rand"
	"time"
)

type SyncSetting int
type MsgType string

const (
	SyncToFromAccunt SyncSetting = 1
	DonotSyncToFromAccount = 2
)

const (
	MsgType_TIMTextElem MsgType = "TIMTextElem"
	MsgType_TIMFaceElem = "TIMFaceElem"
	MsgType_TIMLocationElem = "TIMLocationElem"
	MsgType_TIMCustomElem = "TIMCustomElem"
)

type SendMsgOpt interface {
	ApplyToSendMsgRequest(*SendMsgRequest)
}

type BatchSendMsgOpt interface {
	ApplyToBatchSendMsgRequest(*BatchSendMsgRequest)
}

type SendMsgRequest struct {
	// 1：把消息同步到 From_Account 在线终端和漫游上；
	// 2：消息不同步至 From_Account；若不填写默认情况下会将消息存 From_Account 漫游
	SyncOtherMachine SyncSetting    `json:"SyncOtherMachine,omitempty"`
	// 消息发送方帐号。（用于指定发送消息方帐号）
	FromAccount      string `json:"From_Account,omitempty"`
	// 必填，消息接收方帐号。
	ToAccount        string `json:"To_Account"`
	// 消息离线保存时长（秒数），最长为 7 天（604800s）。若消息只发在线用户，不想保存离线，则该字段填 0。若不填，则默认保存 7 天
	// NOTE(QuantumGhost): 这里为了避免 *int，默认值设置为 604800，如果设置为 0 就会发送 0
	MsgLifeTime      int    `json:"MsgLifeTime"`
	// 必填，消息随机数,由随机函数产生。（用作消息去重）（
	// NOTE(QuantumGhost): 这个字段会在初始化的时候自懂添加
	MsgRandom        int    `json:"MsgRandom"`
	// 消息时间戳，UNIX 时间戳。
	MsgTimeStamp     int    `json:"MsgTimeStamp,omitempty"`
	// 必填，消息内容，具体格式请参考 [消息格式描述](http://bit.ly/2ARGSfV)。
	// （注意，一条消息可包括多种消息元素，MsgBody 为 Array 类型）
	MsgBody          []MsgBody `json:"MsgBody"`
	// 离线推送信息配置，具体可参考 [消息格式描述](http://bit.ly/2AMK9NQ)。
	OfflinePushInfo  *OfflinePushInfo  `json:"OfflinePushInfo,omitempty"`
}

type AndroidOfflinePushInfo struct {
	// 离线推送声音文件路径。
	Sound string `json:"Sound"`
}

type APNSOfflinePushInfo struct {
	// 离线推送声音文件路径。
	Sound     string `json:"Sound"`
	// 这个字段缺省或者为 0 表示需要计数，为 1 表示本条消息不需要计数，即右上角图标数字不增加
	BadgeMode int    `json:"BadgeMode"`
	// 该字段用于标识apns推送的标题，若填写则会覆盖最上层Title
	Title     string `json:"Title"`
	// 该字段用于标识apns推送的子标题
	SubTitle  string `json:"SubTitle"`
	// 该字段用于标识apns携带的图片地址，当客户端拿到该字段时，可以通过下载图片资源的方式将图片展示在弹窗上
	Image     string `json:"Image"`
}

type OfflinePushInfo struct {
	// 0表示推送，1表示不离线推送。
	PushFlag    int    `json:"PushFlag"`
	// 离线推送标题。该字段为ios和android共用
	Title       string `json:"Title"`
	// 离线推送内容 。
	Desc        string `json:"Desc"`
	// 离线推送透传内容。
	Ext         string `json:"Ext"`
	//
	AndroidInfo AndroidOfflinePushInfo `json:"AndroidInfo"`
	ApnsInfo APNSOfflinePushInfo `json:"ApnsInfo"`
}

func (o *OfflinePushInfo) ApplyToBatchSendMsgRequest(req *BatchSendMsgRequest) {
	req.OfflinePushInfo = o
}

func (o *OfflinePushInfo) ApplyToSendMsgRequest(req *SendMsgRequest) {
	req.OfflinePushInfo = o
}

// marker interface
type MsgBody interface {
	Type() MsgType
}

type msgBody struct {
	// TIM消息对象类型，目前支持的消息对象包括： TIMTextElem(文本消息),TIMFaceElem(表情消息)，
	// TIMLocationElem(位置消息)， TIMCustomElem(自定义消息)。
	MsgType    MsgType `json:"MsgType"`
	MsgContent interface{}
}

type customMsgContent struct {
	Data  string `json:"Data"`
	Desc  string `json:"Desc"`
	Ext   string `json:"Ext"`
	Sound string `json:"Sound"`
}

type textMsgContent struct {
	Text string `json:"Text"`
}

func (t *msgBody) Type() MsgType {
	return t.MsgType
}

func (c *Client) SendMsg(ctx context.Context, from string, to string, bodies []MsgBody, opts ...SendMsgOpt) *IMResponse {
	req := c.newRequest(ctx, Service_OPEN_IM, Command_SEND_MSG)
	payload := newSendMsgRequest()
	payload.FromAccount = from
	payload.ToAccount = to
	payload.MsgBody = bodies
	for _, opt := range opts {
		opt.ApplyToSendMsgRequest(&payload)
	}
	result := &IMResponse{}
	resp, err := req.SetResult(result).SetBody(&payload).Post(tencentCloudIMAPIEndpoint)
	fmt.Println(resp)
	if err != nil {
		result.internal = errors.Wrap(err, "error while sendmsg")
	}
	return result
}

func (c *Client) BatchSendMsg(ctx context.Context, from string, to []string, bodies []MsgBody, opts ...BatchSendMsgOpt) *IMResponse {
	req := c.newRequest(ctx, Service_OPEN_IM, Command_SEND_MSG)
	payload := newBatchSendMsgRequest()
	payload.FromAccount = from
	payload.ToAccount = to
	payload.MsgBody = bodies
	for _, opt := range opts {
		opt.ApplyToBatchSendMsgRequest(&payload)
	}
	result := &IMResponse{}
	_, err := req.SetResult(result).SetBody(&payload).Post(tencentCloudIMAPIEndpoint)
	if err != nil {
		result.internal = errors.Wrap(err, "error while sendmsg")
	}
	return result
}

func newSendMsgRequest() SendMsgRequest {
	return SendMsgRequest{
		SyncOtherMachine: SyncToFromAccunt,
		MsgLifeTime: 604800,
		MsgRandom: int(rand.Int31()),
		MsgTimeStamp: int(time.Now().Unix()),
	}
}

func NewCustomMsgBody(desc, ext, data, sound string) MsgBody {
	return &msgBody{
		MsgType: MsgType_TIMCustomElem,
		MsgContent: customMsgContent{Data: data, Desc: desc, Ext: ext, Sound: sound},
	}
}

func NewTextMsgBody(body string) MsgBody {
	return &msgBody{
		MsgType: MsgType_TIMTextElem,
		MsgContent: textMsgContent{Text: body},
	}
}

type BatchSendMsgRequest struct {
	// 1：把消息同步到 From_Account 在线终端和漫游上；
	// 2：消息不同步至 From_Account；若不填写默认情况下会将消息存 From_Account 漫游
	SyncOtherMachine SyncSetting    `json:"SyncOtherMachine,omitempty"`
	// 消息发送方帐号。（用于指定发送消息方帐号）
	FromAccount      string `json:"From_Account,omitempty"`
	// 必填，消息接收方帐号。
	ToAccount        []string `json:"To_Account"`
	// 必填，消息随机数,由随机函数产生。（用作消息去重）（
	// NOTE(QuantumGhost): 这个字段会在初始化的时候自懂添加
	MsgRandom        int    `json:"MsgRandom"`
	// 必填，消息内容，具体格式请参考 [消息格式描述](http://bit.ly/2ARGSfV)。
	// （注意，一条消息可包括多种消息元素，MsgBody 为 Array 类型）
	MsgBody          []MsgBody `json:"MsgBody"`
	// 离线推送信息配置，具体可参考 [消息格式描述](http://bit.ly/2AMK9NQ)。
	OfflinePushInfo  *OfflinePushInfo  `json:"OfflinePushInfo,omitempty"`
}

func newBatchSendMsgRequest() BatchSendMsgRequest {
	return BatchSendMsgRequest{
		SyncOtherMachine: SyncToFromAccunt,
		MsgRandom: int(rand.Int31()),
	}
}
