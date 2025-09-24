package graph_event

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hanyougame/glib/utils/httpc"
	"net/http"
)

type Service struct {
	ctx    context.Context
	config RBConfig
}

func New(ctx context.Context, conf RBConfig) *Service {
	return &Service{
		ctx:    ctx,
		config: conf,
	}
}

// SendEvents 发送事件
func (s *Service) SendEvents(rbLinkID string, eventName RBEventName, extra any) error {
	var err error
	var response *resty.Response
	response, err = httpc.Do(s.ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"event_name": eventName,
			"link_id":    rbLinkID,
			"extra":      extra,
		}).
		Post(s.genUrl())
	if err != nil {
		return fmt.Errorf("send event err: %+v", err)
	}
	if response == nil {
		return fmt.Errorf("response empty")
	}
	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("send event response err, code: %d, body: %s", response.StatusCode(), string(response.Body()))
	}

	return nil
}

// genUrl 构建url
func (s *Service) genUrl() string {
	return fmt.Sprintf(FBEventApiUrl)
}
