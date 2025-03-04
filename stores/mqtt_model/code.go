package mqtt_model

const (
	UserRemoteLoginCode  = 1001 + iota // 用户异地登录通知
	UserPushOffCode                    // 用户被挤下线通知
	UserStatusChangeCode               // 用户状态变更
)

const (
	EnterQuitGameCode = 2001 + iota // 进出游戏通知
)
