package mqtt_model

// EnterQuitGameDelWdlReply 进入游戏转账信息/退出游戏转账信息
type EnterQuitGameDelWdlReply struct {
	GameId         int64 `json:"game_id"`
	Balance        int64 `json:"balance"`
	GamePlatformId int64 `json:"game_platform_id"`
}
