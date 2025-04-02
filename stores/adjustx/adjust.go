package adjustx

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest/httpc"
	"net/http"
)

type Sender struct {
	ctx          context.Context
	url          string
	ChannelId    int64
	AppToken     string
	Auth         string
	EventCodeMap map[AdEventType]string
}

func NewAdjustSender(ctx context.Context, eventConfig Config) *Sender {
	eventMap := lo.SliceToMap(eventConfig.EventCodes, func(c EventCodeConfig) (AdEventType, string) {
		return AdEventType(c.Name), c.Code
	})
	return &Sender{
		ctx:          ctx,
		url:          eventConfig.Url,
		ChannelId:    eventConfig.ChannelID,
		AppToken:     eventConfig.AppToken,
		Auth:         eventConfig.Auth,
		EventCodeMap: eventMap,
	}
}

func (sender *Sender) GetCode(event AdEventType) string {
	val, ok := sender.EventCodeMap[event]
	if ok {
		return val
	}
	logx.Errorf("channel %d, event %s Not Exist", sender.ChannelId, event)
	return ""
}

func (sender *Sender) GetAppToken() string {
	return sender.AppToken
}

func (sender *Sender) GetEnv(env string) string {
	if env == service.DevMode || env == service.TestMode {
		return "sandbox"
	}
	return "production"
}

func (sender *Sender) GetAuth() string {
	return fmt.Sprintf("Bearer %s", sender.Auth)
}

func (sender *Sender) Send(e interface{}) {

	logger := logx.WithContext(sender.ctx)

	response, err := httpc.Do(sender.ctx, http.MethodPost, sender.url, e)
	if err != nil {
		logger.Errorf("[Adjust] S2SEvent response Fail error: %v", err)
		return
	}

	defer response.Body.Close()

	// 解析 JSON 响应
	var data map[string]string
	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		logger.Errorf("[Adjust] S2SEvent response Fail error: %v", err)
		return
	}

	if response.StatusCode != http.StatusOK {
		logger.Errorf("[Adjust] S2SEvent response Fail: %v, %v", response.StatusCode, data)
		return
	}
	logger.Infof("[Adjust] S2SEvent response Success: %v, Response: %v", response.StatusCode, data)
}
