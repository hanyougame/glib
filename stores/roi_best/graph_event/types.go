package graph_event

type RBEventName string

const (
	RBEventNameRegister      RBEventName = "CompleteRegistration" // 注册（标准事件）
	RBEventNameRecharge      RBEventName = "Purchase"             // 充值（标准事件）
	RBEventNameRechargeFirst RBEventName = "AddToCart"            // 首充（标准事件）
)

type ChannelID string

const (
	ChannelIDFacebook      ChannelID = "4"
	ChannelIDTiktok        ChannelID = "5"
	ChannelIDKwai          ChannelID = "9"
	ChannelIDGoogle        ChannelID = "10"
	ChannelIDXiaoBuWangLuo ChannelID = "20"
	ChannelIDOKSPIN        ChannelID = "21"
	ChannelIDSNAPTUBE      ChannelID = "22"
	ChannelIDBIGO          ChannelID = "23"
	ChannelIDAPPLUCK       ChannelID = "24"
	ChannelIDTrafficstars  ChannelID = "34"
)
