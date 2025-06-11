package adjustx

type Config struct {
	Url        string            `json:"url"`
	ChannelID  int64             `json:"channel_id"`
	AppToken   string            `json:"app_token"`
	Auth       string            `json:"auth,optional"`
	EventCodes []EventCodeConfig `json:"event_codes"`
	PixelId    string            `json:"pixel_id"`
}

type EventCodeConfig struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type DefaultEvent struct {
	AppToken      string `form:"app_token" structs:"app_token"`
	S2S           int    `form:"s2s" structs:"s2s"`
	EventToken    string `form:"event_token" structs:"event_token"`         // 事件码
	AdId          string `form:"adid" structs:"adid"`                       // 外部设备ID
	CreatedAtUnix int64  `form:"created_at_unix" structs:"created_at_unix"` // 事件时间戳
	Environment   string `form:"environment" structs:"environment"`         // 环境 sandbox | production
	Authorization string `header:"authorization" structs:"-"`               // Authorization
	IpAddress     string `form:"ip_address" structs:"ip_address"`           // IP地址
	UserAgent     string `form:"user_agent" structs:"user_agent"`           // useragent
	GpsAdId       string `form:"gps_adid" structs:"gps_adid,omitempty"`     // GpsAdId
}

type RevenueEvent struct {
	AppToken      string  `form:"app_token" structs:"app_token"`
	S2S           int     `form:"s2s" structs:"s2s"`
	EventToken    string  `form:"event_token" structs:"event_token"`         // 事件码
	AdId          string  `form:"adid" structs:"adid"`                       // 外部设备ID
	CreatedAtUnix int64   `form:"created_at_unix" structs:"created_at_unix"` // 事件时间戳
	Environment   string  `form:"environment" structs:"environment"`         // 环境 sandbox | production
	Authorization string  `header:"authorization" structs:"-"`               // Authorization
	Currency      string  `form:"currency" structs:"currency"`               // 币种
	Revenue       float64 `form:"revenue" structs:"revenue"`                 // 金额
	IpAddress     string  `form:"ip_address" structs:"ip_address"`           // IP地址
	UserAgent     string  `form:"user_agent" structs:"user_agent"`           // useragent
	GpsAdId       string  `form:"gps_adid" structs:"gps_adid,omitempty"`     // GpsAdId
}
