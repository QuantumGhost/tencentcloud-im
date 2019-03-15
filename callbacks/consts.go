package callbacks

const (
	groupCommandPrefix = "Group."
)

const (
	// 申请入群之前回调
	CallbackBeforeApplyJoinGroup = groupCommandPrefix + "CallbackBeforeApplyJoinGroup"
	// 拉人入群之前回调
	CallbackBeforeInviteJoinGroup = groupCommandPrefix + "CallbackBeforeInviteJoinGroup"
	// 新成员入群之后回调
	CallbackAfterNewMemberJoin = groupCommandPrefix + "CallbackAfterNewMemberJoin"
	// 群成员离开之后回调
	CallbackAfterMemberExit = groupCommandPrefix + "CallbackAfterMemberExit"
	// 创建群组之前回调
	CallbackBeforeCreateGroup = groupCommandPrefix + "CallbackBeforeCreateGroup"
	// 群组资料修改之后回调
	CallbackAfterGroupInfoChanged = groupCommandPrefix + "CallbackAfterGroupInfoChanged"
)
