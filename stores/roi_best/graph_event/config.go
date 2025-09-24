package graph_event

import (
	"database/sql/driver"
	"fmt"
	"github.com/zeromicro/go-zero/core/jsonx"
)

const (
	FBEventApiUrl = "https://sdk-report.roibestopenapi.com/report/fb/event" // api url
)

type RBConfig struct {
	RBLinkID     string        `json:"rb_link_id"`     //由系统生成的链接 ID
	RBEventsName []RBEventName `json:"rb_events_name"` // 事件名称
	RBChannelID  string        `json:"rb_channel_id"`  // 埋点平台渠道
	RBRbPixelID  string        `json:"rb_pixel_id"`    // 像素ID
}

func (s *RBConfig) Scan(value interface{}) error {
	var (
		err    error
		result RBConfig
	)

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid config val: %+v", value)
	}
	if err = jsonx.Unmarshal(bytes, &result); err != nil {
		return fmt.Errorf("decode config err: %+v", err)
	}
	*s = result

	return nil
}

func (s *RBConfig) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}

	return jsonx.Marshal(s)
}
