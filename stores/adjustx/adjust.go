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

func (sender *Sender) Send(e interface{}) error {

	logger := logx.WithContext(sender.ctx)

	m := structs.Map(e)
	formValues := url.Values{}
	for k := range m {
		formValues.Set(k, cast.ToString(m[k]))
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, sender.url, nil)
	if err != nil {
		return fmt.Errorf("S2SEvent new request err: %v", err)
	}
	req.URL.RawQuery = formValues.Encode()
	req.Header.Set(httpx.ContentType, "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", sender.GetAuth())
	response, err := httpc.DoRequest(req)
	if err != nil {
		return fmt.Errorf("S2SEvent do request err: %v", err)
	}
	defer response.Body.Close()

	// 解析 JSON 响应
	var data map[string]string
	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return fmt.Errorf("S2SEvent response err: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("S2SEvent response status err, status: %v, data: %v", response.StatusCode, data)
	}

	logger.Infof("[Adjust] S2SEvent response Success: %v, Response: %v", response.StatusCode, data)

	return nil
}
