package graph_event

const (
	FBEventApiUrl = "https://sdk-report.roibestopenapi.com/report/fb/event" // api url
)

type RBConfig struct {
	RBLinkID     string        `json:"link_id"`     //由系统生成的链接 ID
	RBEventsName []RBEventName `json:"events_name"` // 事件名称
	RBChannelID  string        `json:"channel_id"`  // 埋点平台渠道
	RBRbPixelID  string        `json:"pixel_id"`    // 像素ID
}
