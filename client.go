package tencentcloud_im

import (
	"context"
	"fmt"
	"gopkg.in/resty.v1"
	"math/rand"
	"strconv"
	"time"
)

const (
	tencentCloudIMBaseUrl     = "https://console.tim.qq.com/v4"
	tencentCloudIMAPIEndpoint = tencentCloudIMBaseUrl + "/{serviceName}/{command}"
)

const (
	Service_IM_OPEN_LOGIN_SVC   = "im_open_login_svc"
	Service_OPEN_IM             = "openim"
	Service_PROFILE             = "profile"
	Service_GROUP_OPEN_HTTP_SVC = "group_open_http_svc"
)

const (
	// im_open_login_svc
	Command_ACCOUNT_IMPORT      = "account_import"
	Command_MULTIACCOUNT_IMPORT = "multiaccount_import"
	Command_KICK                = "kick"
	// openim
	Command_SEND_MSG       = "sendmsg"
	Command_BATCH_SEND_MSG = "batchsendmsg"
	Command_QUERY_STATE    = "querystate"
	// profile
	Command_PORTRAIT_SET = "portrait_set"
	// group_open_http_svc
	Command_CREATE_GROUP          = "create_group"
	Command_ADD_GROUP_MEMBER      = "add_group_member"
	Command_GROUP_MSG_GET_SIMPLE  = "group_msg_get_simple"
	Command_GET_GROUP_MEMBER_INFO = "get_group_member_info"
	Command_DELETE_GROUP_MEMBER   = "delete_group_member"
)

type IMResponse struct {
	ActionStatus string `json:"ActionStatus"`
	ErrorInfo    string `json:"ErrorInfo"`
	ErrorCode    int    `json:"errorCode"`
	internal     error
	imErr        IMError
}

type IMError interface {
	error
	ErrorInfo() string
	ErrorCode() int
}

type imError struct {
	actionStatus string
	errorInfo    string
	errorCode    int
}

func (e *imError) ErrorInfo() string {
	return e.errorInfo
}

func (e *imError) ErrorCode() int {
	return e.errorCode
}

func (e *imError) Error() string {
	return fmt.Sprintf("Error from tencentcloud, desc %s, code %d", e.errorInfo, e.errorCode)
}

func (r *IMResponse) Error() error {
	if r.internal != nil {
		return r.internal
	}
	if r.ErrorCode == 0 {
		return nil
	}
	if r.imErr != nil {
		return r.imErr
	}
	r.imErr = &imError{errorCode: r.ErrorCode, errorInfo: r.ErrorInfo, actionStatus: r.ActionStatus}
	return r.imErr
}

func ToIMError(err error) (IMError, bool) {
	r, ok := err.(*imError)
	return r, ok
}

type Client struct {
	sdkAppId        int
	adminIdentifier string
	userSig         string
	client          *resty.Client
}

func (c *Client) preRequestHook(_ *resty.Client, req *resty.Request) error {
	req.SetQueryParam("random", strconv.Itoa(rand.Intn(1<<31-1)))
	return nil
}

func (c *Client) SetClient(client *resty.Client) {
	c.client = client
	client.SetQueryParam("sdkappid", strconv.Itoa(c.sdkAppId)).
		SetQueryParam("identifier", c.adminIdentifier).
		SetQueryParam("usersig", c.userSig).
		SetQueryParam("contenttype", "json").
		SetHeader("Content-Type", "application/json").
		SetRESTMode().
		OnBeforeRequest(c.preRequestHook)
}

func NewClient(appId int, identifier string, urlsig string) *Client {
	restyClient := resty.New().SetTimeout(2 * time.Second)
	client := &Client{sdkAppId: appId, adminIdentifier: identifier, userSig: urlsig}
	client.SetClient(restyClient)
	return client
}

func (c *Client) newRequest(ctx context.Context, serviceName string, command string) *resty.Request {
	req := c.client.NewRequest().SetContext(ctx)
	req.SetPathParams(map[string]string{"serviceName": serviceName, "command": command})
	return req
}
