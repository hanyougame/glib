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
	config Config
}

func New(ctx context.Context, conf Config) *Service {
	return &Service{
		ctx:    ctx,
		config: conf,
	}
}

// SendEvents 发送事件
func (s *Service) SendEvents(request *EventsSendRequest) error {
	var err error
	if err = request.verify(); err != nil {
		return fmt.Errorf("verify request err: %+v", err)
	}
	var response *resty.Response
	response, err = httpc.Do(s.ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
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
	return fmt.Sprintf(GraphEventApiUrl, GraphEventApiVersion, s.config.PixelID, s.config.AccessToken)
}
