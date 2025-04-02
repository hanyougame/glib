package adjustx

type AdEventType string

const (
	AdLoginEvent          AdEventType = "AdLoginEvent"          // 登录
	AdRegisterEvent       AdEventType = "AdRegisterEvent"       // 注册
	AdFirstRechargeEvent  AdEventType = "AdFirstRechargeEvent"  // 首充 携带金额 第一次充值成功
	AdSecondRechargeEvent AdEventType = "AdSecondRechargeEvent" // 二充 携带金额 第二次充值成功
	AdThirdRechargeEvent  AdEventType = "AdThirdRechargeEvent"  // 三充 携带金额 第三次充值成功
	AdFourthRechargeEvent AdEventType = "AdFourthRechargeEvent" // 四充 携带金额 第四次充值成功
	AdPurchaseEvent       AdEventType = "AdPurchaseEvent"       // 充值 携带金额 每次充值成功
	AdWithdrawEvent       AdEventType = "AdWithdrawEvent"       // 提现 携带金额 每次提现成功
)
