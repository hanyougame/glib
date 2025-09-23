package graph_event

type EventName string

const (
	EventNameRegister      EventName = "CompleteRegistration" // 注册（标准事件）
	EventNameRecharge      EventName = "Purchase"             // 充值（标准事件）
	EventNameRechargeFirst EventName = "AddToCart"            // 首充（标准事件）
)
