package adjustx

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest/httpc"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"net/url"
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

	m := structs.Map(e)
	formValues := url.Values{}
	for k := range m {
		formValues.Set(k, cast.ToString(m[k]))
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, sender.url, nil)
	if err != nil {
		logger.Errorf("[Adjust] S2SEvent request Fail error: %v", err)
		return
	}
	req.URL.RawQuery = formValues.Encode()
	req.Header.Set(httpx.ContentType, "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", sender.GetAuth())
	response, err := httpc.DoRequest(req)
	if err != nil {
		logger.Errorf("[Adjust] S2SEvent request Fail error: %v", err)
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
