package graph_event

const (
	FBEventApiUrl = "https://sdk-report.roibestopenapi.com/report/fb/event" // api url
)

type RBConfig struct {
	RBLinkID     string        `json:"rb_link_id"`     //由系统生成的链接 ID
	RBEventsName []RBEventName `json:"rb_events_name"` // 事件名称
	RBChannelID  string        `json:"rb_channel_id"`  // 埋点平台渠道
	RBRbPixelID  string        `json:"rb_pixel_id"`    // 像素ID
}
