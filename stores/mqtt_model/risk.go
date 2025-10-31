package mqtt_model

// RiskRechargeChannelNewAlertNotify 充值通道新告警通知
type RiskRechargeChannelNewAlertNotify struct {
	AlertIDs []int64 `json:"alert_ids"` // 告警id列表
}
