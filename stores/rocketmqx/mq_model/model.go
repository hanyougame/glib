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
	SupplySign         bool   `json:"supply_sign"`          // 补单标志
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
	BetTime        int64  `json:"bet_time"`         // 投注时间
	BetAmount      int64  `json:"bet_amount"`       //投注金额(不包含撤单金额)
	ValidBetAmount int64  `json:"valid_bet_amount"` //有效投注金额
	GameId         int64  `json:"game_id"`          //游戏ID
	GameCategory   int64  `json:"game_category"`    //游戏类型
	PlatformID     int64  `json:"platform_id"`      // 游戏平台ID
	SettlementTime int64  `json:"settlement_time"`  // 领取时间
	BonusAmount    int64  `json:"bonus_amount"`     // 派奖金额
	UserWinAmount  int64  `json:"user_win_amount"`  // 用户输赢金额 有负数
}

// UserPromotionBonusNotify 用户优惠奖励领取通知
type UserPromotionBonusNotify struct {
	UserId          int64  `json:"user_id"` //用户ID
	UserAccount     string `json:"user_account"`
	CurrencyCode    string `json:"currency_code"`
	BonusAmount     int64  `json:"bonus_amount"`     //彩金金额
	ReceiveTime     int64  `json:"receive_time"`     // 领取时间
	PromotionSource int64  `json:"promotion_source"` // 优惠来源，同后台 constants.PromotionSource
}

type UserRegisterNotify struct {
	ParentID     int64  `json:"parent_id"`     // 上级ID
	UserId       int64  `json:"user_id"`       // 用户ID
	UserAccount  string `json:"user_account"`  // 用户账号
	RegisterTime int64  `json:"register_time"` // 注册时间
	CurrencyCode string `json:"currency_code"` // 币种
}

// MqDepWdlTripartiteMsg 游戏转入转出第三方余额通知
type MqDepWdlTripartiteMsg struct {
	GamePlatformId   int64  `json:"game_platform_id"`  //游戏平台key  唯一标识
	GameId           int64  `json:"game_id"`           //游戏标志
	UserId           int64  `json:"user_id"`           //用户ID
	CurrencyCode     string `json:"currency_code"`     //货币Code
	PlatformAmount   int64  `json:"platform_amount"`   // 整数大于0 平台金额格式
	TripartiteAmount string `json:"tripartite_amount"` // 保留2位小数 三方金额格式
	GameOrderNo      string `json:"game_order_no"`     // 订单号
	TransferType     int64  `json:"transfer_type"`     //交易类型 1转入,2转出
	GameName         string `json:"game_name"`         //游戏名称
}

// MqttPublish 发布mqtt消息
type MqttPublish struct {
	ClientIDList []string `json:"client_id_list"` // 客户端id列表
	ServiceID    string   `json:"service_id"`     // 服务id
	Type         int64    `json:"type"`           // 1 全部 2 部分
	Message      string   `json:"message"`        // 发送的内容
	Code         int64    `json:"code"`           // 事件的code
}

type OrderTimeout struct {
	OrderNo string `json:"order_no"`      // 订单号
	UserID  int64  `json:"user_id"`       // 用户ID
}

// 充值日志
type RechargeLog struct {
	ID                  int64  `json:"id"`                          // Id            
	RechargeType        int8   `json:"recharge_type"`               // 充值类型 1-在线充值 2-转账充值 3-客服代充
	PaymentId           int64  `json:"payment_id"`                  // paymentId
	RechargeCategoryId  int64  `json:"recharge_category_id"`  	// 充值大类Id
	RechargeMerchantId  int64  `json:"recharge_merchant_id"`   	// 充值商户Id
	RechargeChannelId   int64  `json:"recharge_channel_id"`     	// 充值渠道Id
	RechargeOrderNumber string `json:"recharge_order_number"` 	// 充值订单号
	RechargeAmount      int64  `json:"recharge_amount"`             // 充值金额
	GiftAmount          int64  `json:"gift_amount"`             	// 赠送金额
	ReceiveAmount       int64  `json:"receive_amount"`             	// 实际到账金额
	UserId              int64  `json:"user_id"`                  	// 充值用户
	CurrencyCode        string `json:"currency_code"`               // 充值币种
	RechargeTime        int64  `json:"recharge_time"`               // 充值时间 通过该时间进行统计
	SuccessTime         int64  `json:"success_time"`                // 充值成功时间
	Status              int    `json:"status"`                      // 充值状态 0-无状态（订单被创建） 1-充值成功 2-充值失败
}
