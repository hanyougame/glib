package adjustx

type Config struct {
	Url        string            `json:"url"`
	ChannelID  int64             `json:"channel_id"`
	AppToken   string            `json:"app_token"`
	Auth       string            `json:"auth,optional"`
	EventCodes []EventCodeConfig `json:"event_codes"`
}

type EventCodeConfig struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type DefaultEvent struct {
	AppToken      string `form:"app_token"`
	S2S           int    `form:"s2s"`
	EventToken    string `form:"event_token"`     // 事件码
	AdId          string `form:"adid"`            // 外部设备ID
	CreatedAtUnix int64  `form:"created_at_unix"` // 事件时间戳
	Environment   string `form:"environment"`     // 环境 sandbox | production
	Authorization string `header:"authorization"` // Authorization
}

type RevenueEvent struct {
	AppToken      string `form:"app_token"`
	S2S           int    `form:"s2s"`
	EventToken    string `form:"event_token"`     // 事件码
	AdId          string `form:"adid"`            // 外部设备ID
	CreatedAtUnix int64  `form:"created_at_unix"` // 事件时间戳
	Environment   string `form:"environment"`     // 环境 sandbox | production
	Authorization string `header:"authorization"` // Authorization
	Currency string  `form:"currency"` // 币种
	Revenue  float64 `form:"revenue"`  // 金额
}
