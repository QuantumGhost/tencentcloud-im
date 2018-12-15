package tencentcloud_im

import (
	"context"
	"github.com/pkg/errors"
)

type ProfileTag string

const (
	ProfileTag_Tag_Profile_IM_Nick            ProfileTag = "Tag_Profile_IM_Nick"
	ProfileTag_Tag_Profile_IM_Gender                     = "Tag_Profile_IM_Gender"
	ProfileTag_Tag_Profile_IM_BirthDay                   = "Tag_Profile_IM_BirthDay"
	ProfileTag_Tag_Profile_IM_Location                   = "Tag_Profile_IM_Location"
	ProfileTag_Tag_Profile_IM_SelfSignature              = "Tag_Profile_IM_SelfSignature"
	ProfileTag_Tag_Profile_IM_AllowType                  = "Tag_Profile_IM_AllowType"
	ProfileTag_Tag_Profile_IM_Language                   = "Tag_Profile_IM_Language"
	ProfileTag_Tag_Profile_IM_Image                      = "Tag_Profile_IM_Image"
	ProfileTag_Tag_Profile_IM_MsgSettings                = "Tag_Profile_IM_MsgSettings"
	ProfileTag_Tag_Profile_IM_AdminForbidType            = "Tag_Profile_IM_AdminForbidType"
	ProfileTag_Tag_Profile_IM_Level                      = "Tag_Profile_IM_Level"
	ProfileTag_Tag_Profile_IM_Role                       = "Tag_Profile_IM_Role"
)

type ProfileItem struct {
	Tag   ProfileTag  `json:"Tag"`
	Value interface{} `json:"Value"`
}

type PortraitSetResponse struct {
	IMResponse
	ErrorDisplay string
}

type PortraitSetRequest struct {
	FromAccount string        `json:"From_Account"`
	ProfileItem []ProfileItem `json:"ProfileItem"`
}

func (c *Client) PortraitSet(ctx context.Context, accountId string, items []ProfileItem) *PortraitSetResponse {
	req := c.newRequest(ctx, Service_PROFILE, Command_PORTRAIT_SET)
	payload := &PortraitSetRequest{FromAccount: accountId, ProfileItem: items}
	result := &PortraitSetResponse{}
	_, err := req.SetResult(result).SetBody(&payload).Post(tencentCloudIMAPIEndpoint)
	if err != nil {
		result.internal = errors.Wrap(err, "err while portrait_set")
	}
	return result
}

func NewNickNameProfileItem(nickname string) ProfileItem {
	return ProfileItem{
		Tag: ProfileTag_Tag_Profile_IM_Nick, Value: nickname,
	}
}
