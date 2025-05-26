package graph_event

import (
	"database/sql/driver"
	"fmt"
	"github.com/zeromicro/go-zero/core/jsonx"
)

const (
	GraphEventApiUrl     = "https://graph.facebook.com/%s/%s/events?access_token=%s" // api url
	GraphEventApiVersion = "v22.0"                                                   // api 版本号
)

type Config struct {
	PixelID     string      `gorm:"pixel_id" json:"pixel_id"`                // 像素ID
	AccessToken string      `gorm:"column:access_token" json:"access_token"` // token
	EventsName  []EventName `gorm:"column:events_name" json:"events_name"`   // 事件名称
}

func (s *Config) Scan(value interface{}) error {
	var (
		err    error
		result Config
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

func (s *Config) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}

	return jsonx.Marshal(s)
}
