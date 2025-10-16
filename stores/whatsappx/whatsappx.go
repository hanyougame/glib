package whatsappx

import (
	"context"
	"github.com/piusalfred/whatsapp/config"
	"github.com/piusalfred/whatsapp/message"
	whttp "github.com/piusalfred/whatsapp/pkg/http"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

var GlobalSvc *WhatsAppX

type WhatsAppX struct {
	conf       *Config
	baseClient *message.BaseClient
}

func NewWhatsAppX(c *Config) (*WhatsAppX, error) {
	coreClient := whttp.NewSender[message.Message]()
	coreClient.SetHTTPClient(http.DefaultClient)
	reader := config.ReaderFunc(func(ctx context.Context) (*config.Config, error) {
		conf := &config.Config{
			BaseURL:           c.BaseURL,
			APIVersion:        c.APIVersion,
			AccessToken:       c.AccessToken,
			PhoneNumberID:     c.PhoneNumberID,
			BusinessAccountID: c.BusinessAccountID,
			AppSecret:         c.AppSecret,
			AppID:             c.AppID,
			SecureRequests:    c.SecureRequests,
		}
		return conf, nil
	})
	baseClient, err := message.NewBaseClient(coreClient, reader)
	if err != nil {
		logx.Errorf("error creating base client. err:%v", err)
		return nil, err
	}

	svc := &WhatsAppX{
		conf:       c,
		baseClient: baseClient,
	}
	GlobalSvc = svc
	return svc, nil
}

func Must(c *Config) {
	_, err := NewWhatsAppX(c)
	if err != nil {
		logx.Errorf("create whatsapp client failed. err:%v", err)
	}
}
