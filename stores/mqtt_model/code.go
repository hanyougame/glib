package mqtt_model

const (
	UserRemoteLoginCode  = 1001 + iota // 用户异地登录通知
	UserPushOffCode                    // 用户被挤下线通知
	UserStatusChangeCode               // 用户状态变更
)

const (
	EnterQuitGameCode = 2001 + iota // 进出游戏通知
)

const (
	ActivityRechargePopWindow   = 3001 + iota // 活动充值弹窗
	ActivityRechargeRedEnvelope               // 活动抢红包
)

const (
	WithdrawApplyCode             = 4001 + iota // 发起提现申请
	WithdrawFailCode                            // 提现三方回调失败
	WithdrawRiskControlReviewCode               // 提现风控审核
)

const (
	GameAwardMonitoringCode = 5001 + iota // 游戏派奖监控后台警告
)

const (
	ExportSuccessCode = 6001 + iota
	ExportFailCode
)
