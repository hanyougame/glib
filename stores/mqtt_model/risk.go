package mqtt_model

// RechargeChannelNewAlert 充值通道新告警
type RechargeChannelNewAlert struct {
	IDs []int64 `json:"ids"` // 告警id列表
}
