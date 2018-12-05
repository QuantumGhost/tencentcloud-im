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

type MultiAccountImportResponse struct {
	IMResponse
	FailAccounts []string `json:"FailAccounts"`
}

func (c *Client) MultiAccountImport(ctx context.Context, identifiers []string) *MultiAccountImportResponse {
	type multiAccountImportRequest struct {
		Accounts []string `json:"Accounts"`
	}
	req := c.newRequest(ctx, Service_IM_OPEN_LOGIN_SVC, Command_MULTIACCOUNT_IMPORT)
	result := &MultiAccountImportResponse{}
	_, err := req.SetBody(multiAccountImportRequest{Accounts: identifiers}).SetResult(result).
		Post(tencentCloudIMAPIEndpoint)
	if err != nil {
		result.internal = errors.Wrap(err, "error while multiaccount import")
	}
	return result
}

func (c *Client) Kick(ctx context.Context, identifier string) *IMResponse {
	type kickRequest struct {
		Identifier string `json:"Identifier"`
	}
	req := c.newRequest(ctx, Service_IM_OPEN_LOGIN_SVC, Command_KICK)
	result := &IMResponse{}
	_, err := req.SetBody(kickRequest{Identifier: identifier}).SetResult(result).
		Post(tencentCloudIMAPIEndpoint)
	if err != nil {
		result.internal = errors.Wrap(err, "error while kick user")
	}
	return result
}
