package mq_model

// UserLoginNotify 用户登陆通知
type UserLoginNotify struct {
	UserId      int64  `json:"user_id"` //用户ID
	UserAccount string `json:"user_account"`
	FirstToday  bool   `json:"first_today"` //今天首次(有最好)
	LoginTime   int64  `json:"login_time"`  //登陆时间
}

// UserRechargeNotify 用户成功充值推送
type UserRechargeNotify struct {
	UserId             int64  `json:"user_id"`              // 用户ID
	UserAccount        string `json:"user_account"`         // 用户账号
	CurrencyCode       string `json:"currency_code"`        // 币种
	RechargeAmount     int64  `json:"recharge_amount"`      // 充值金额
	RechargeTime       int64  `json:"recharge_time"`        // 充值时间
	FirstSign          bool   `json:"first_sign"`           // 首充标志
	RechargeType       int64  `json:"recharge_type"`        // 充值类型 1-在线充值 2-转账充值 3-客服代充
	RechargeMerchantId int64  `json:"recharge_merchant_id"` // 充值商户Id
	RechargeCategoryId int64  `json:"recharge_category_id"` // 充值类别Id
	RechargeChannelId  int64  `json:"recharge_channel_id"`  // 充值通道Id
	GiftAmount         int64  `json:"gift_amount"`          // 赠送金额
}

// UserWithdrawNotify 用户提现推送
type UserWithdrawNotify struct {
	UserId          int64  `json:"user_id"`           //用户ID
	UserAccount     string `json:"user_account"`      //用户名
	CurrencyCode    string `json:"currency_code"`     // 币种
	WithdrawAmount  int64  `json:"withdraw_amount"`   //提现金额
	WithdrawTime    int64  `json:"withdraw_time"`     // 提现时间
	WithdrawEndTime int64  `json:"withdraw_end_time"` // 提现结束时间
}

// UserBetSettlementNotify 用户投注结算通知
type UserBetSettlementNotify struct {
	UserId         int64  `json:"user_id"` //用户ID
	UserAccount    string `json:"user_account"`
	CurrencyCode   string `json:"currency_code"`
	BetAmount      int64  `json:"bet_amount"`       //投注金额(不包含撤单金额)
	WinAmount      int64  `json:"win_amount"`       //中奖金额
	ValidBetAmount int64  `json:"valid_bet_amount"` //有效投注金额
	GameId         int64  `json:"game_id"`          //游戏ID
	GameCategory   int64  `json:"game_category"`    //游戏类型
	PlatformID     int64  `json:"platform_id"`      // 游戏平台ID
	SettlementTime int64  `json:"settlement_time"`  // 领取时间
}

// UserPromotionBonusNotify 用户优惠奖励领取通知
type UserPromotionBonusNotify struct {
	UserId          int64  `json:"user_id"` //用户ID
	UserAccount     string `json:"user_account"`
	CurrencyCode    string `json:"currency_code"`
	BonusAmount     int64  `json:"bonus_amount"`     //彩金金额
	ReceiveTime     int64  `json:"receive_time"`     // 领取时间
	PromotionSource int64  `json:"promotion_source"` // 优惠来源
}

// AgentCommissionReceiveNotify 代理佣金领取通知
type AgentCommissionReceiveNotify struct {
	UserId        int64  `json:"user_id"`        // 用户Id
	CurrencyCode  string `json:"currency_code"`  // 用户币种
	ReceiveTime   int64  `json:"receive_time"`   // 领取时间
	ReceiveAmount int64  `json:"receive_amount"` // 领取金额
}

type UserRegisterNotify struct {
	ParentID     int64  `json:"parent_id"`     // 上级ID
	UserId       int64  `json:"user_id"`       // 用户ID
	UserAccount  string `json:"user_account"`  // 用户账号
	RegisterTime int64  `json:"register_time"` // 注册时间
	CurrencyCode string `json:"currency_code"` // 币种
}
