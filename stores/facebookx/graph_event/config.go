package graph_event

const (
	GraphEventApiUrl     = "https://graph.facebook.com/%s/%s/events?access_token=%s" // api url
	GraphEventApiVersion = "v22.0"                                                   // api 版本号
)

type Config struct {
	PixelID     string      `gorm:"pixel_id" json:"pixel_id"`                // 像素ID
	AccessToken string      `gorm:"column:access_token" json:"access_token"` // token
	EventsName  []EventName `gorm:"column:events_name" json:"events_name"`   // 事件名称
}
