package mqtt_model

// UserPushOff 用户被挤下线通知
type UserPushOff struct {
	Ip        string `json:"ip"`         // ip
	Device    string `json:"device"`     // 设备型号
	LoginTime int64  `json:"login_time"` // 登录时间
}
