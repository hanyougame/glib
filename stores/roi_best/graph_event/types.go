package graph_event

type RBEventName string

const (
	RBEventNameRegister      RBEventName = "CompleteRegistration" // 注册（标准事件）
	RBEventNameRecharge      RBEventName = "Purchase"             // 充值（标准事件）
	RBEventNameRechargeFirst RBEventName = "AddToCart"            // 首充（标准事件）
)
