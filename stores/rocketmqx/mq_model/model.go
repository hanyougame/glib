package mq_model

// UserLoginNotify 用户登陆通知
type UserLoginNotify struct {
	UserId                 int64 //用户ID
	FirstToday             bool  //今天首次(有最好)
	LoginTime              int64 //登陆时间
	FirstRechargeTime      int64 // 首充时间
	FirstBetSettlementTime int64 // 首次投注时间
}

// UserRechargeNotify 用户成功充值推送
type UserRechargeNotify struct {
	UserId         int64 //用户ID
	RechargeAmount int64 //充值金额
	RechargeTime   int64 //充值时间
	FirstSign      bool  //首充标志
}

// UserWithdrawNotify 用户提现推送
type UserWithdrawNotify struct {
	UserId          int64 //用户ID
	WithdrawAmount  int64 //提现金额
	WithdrawEndTime int64 // 提现结束时间
}

// UserBetSettlementNotify 用户投注结算通知
type UserBetSettlementNotify struct {
	UserId         int64 //用户ID
	BetAmount      int64 //投注金额(不包含撤单金额)
	WinAmount      int64 //中奖金额
	ValidBetAmount int64 //有效投注金额
	GameId         int64 //游戏ID
	GameCategory   int64 //游戏类型
	SettlementTime int64 // 领取时间
}

// UserPromotionBonusNotify 用户优惠奖励领取通知
type UserPromotionBonusNotify struct {
	UserId          int64 //用户ID
	BonusAmount     int64 //彩金金额
	ReceiveTime     int64 // 领取时间
	PromotionSource int64 // 优惠来源
}
