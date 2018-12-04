package tencentcloud_im

import (
	"context"
	"github.com/pkg/errors"
)

type AccountImportRequest struct {
	Identifier string `json:"Identifier"`
	Nick       string `json:"Nick"`
	FaceURL    string `json:"FaceUrl"`
	Type       int    `json:"Type"`
}

func (c *Client) AccountImport(ctx context.Context, request AccountImportRequest) *IMResponse {
	req := c.newRequest(ctx, Service_IM_OPEN_LOGIN_SVC, Command_ACCOUNT_IMPORT)
	result := &IMResponse{}
	_, err := req.SetBody(request).SetResult(result).Post(tencentCloudIMAPIEndpoint)
	if err != nil {
		result.internal = errors.Wrap(err, "error while account import")
	}
	return result
}

func (c *Client) MultiAccountImport(identifiers []string) *IMResponse {
	panic("not implemented")
}
