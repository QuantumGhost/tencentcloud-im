package tencentcloud_im

import (
	"context"
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
	Service_IM_OPEN_LOGIN_SVC = "im_open_login_svc"
	Service_OPEN_IM           = "openim"
)

const (
	Command_ACCOUNT_IMPORT      = "account_import"
	Command_MULTIACCOUNT_IMPORT = "multiaccount_import"
	Command_SEND_MSG            = "sendmsg"
)

type IMResponse struct {
	ActionStatus string `json:"ActionStatus"`
	ErrorInfo    string `json:"ErrorInfo"`
	ErrorCode    int    `json:"ErrorCode"`
	internal     error  `json:"-"`
}

func (r *IMResponse) Error() string {
	if r.internal != nil {
		return r.internal.Error()
	}
	return r.ErrorInfo
}

func (r *IMResponse) Succeed() bool {
	return r.ErrorCode == 0 && r.internal == nil
}

func GetErrorCode(err error) (int, bool) {
	if r, ok := err.(*IMResponse); ok {
		return r.ErrorCode, true
	}
	return 0, false
}

func GetInternalError(err error) (error, bool) {
	if r, ok := err.(*IMResponse); ok {
		return r.internal, true
	}
	return nil, false
}

func (r *IMResponse) Cause() error {
	return r.internal
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
		SetPreRequestHook(c.preRequestHook)
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
