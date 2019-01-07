package consts

import (
	"database/sql/driver"
	"github.com/QuantumGhost/tencentcloud-im/errors"
)

type ApplyJoinOption string

const (
	DisableApply   ApplyJoinOption = "DisableApply"
	NeedPermission                 = "NeedPermission"
	FreeAccess                     = "FreeAccess"
)

func (a *ApplyJoinOption) Scan(src interface{}) error {
	switch v := src.(type) {
	default:
		return errors.ErrInvalidType
	case string:
		*a = ApplyJoinOption(v)
		return nil
	}
}

func (a ApplyJoinOption) Value() (driver.Value, error) {
	return string(a), nil
}

type GroupType string

const (
	// 公开群
	Public GroupType = "Public"
	// 私有群
	Private = "Private"
	// 聊天室
	ChatRoom = "ChatRoom"
	// 互动直播聊天室
	AVChatRoom = "AVChatRoom"
	// 在线成员广播大群
	BChatRoom = "BChatRoom"
)

func (t GroupType) Value() (driver.Value, error) {
	return string(t), nil
}

func (t *GroupType) Scan(src interface{}) error {
	switch v := src.(type) {
	default:
		return errors.ErrInvalidType
	case string:
		*t = GroupType(v)
		return nil
	}
}

type GroupRole string

const (
	GroupRoleAdmin GroupRole = "Admin"
)
