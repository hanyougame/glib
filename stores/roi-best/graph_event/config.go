package graph_event

const (
	FBEventApiUrl = "https://sdk-report.roibestopenapi.com/report/fb/event" // api url
)

type Config struct {
	LinkID     string      `json:"link_id"`     //由系统生成的链接 ID
	EventsName []EventName `json:"events_name"` // 事件名称
	Extra      any         `json:"extra"`       // 事件的额外参数。如果对应的事件需要传递额外参数放在本字段中，ROIBest 将会透传至对应的广告平台。
	ChannelID  string      `json:"channel_id"`  // 埋点平台渠道
	RbPixelID  string      `json:"rb_pixel_id"` // 像素ID
}
