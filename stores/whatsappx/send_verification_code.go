package whatsappx

import (
	"context"
	"errors"
	"github.com/piusalfred/whatsapp/message"
	"github.com/zeromicro/go-zero/core/logx"
)

func SendVerificationCode(ctx context.Context, recipient, code string) error {
	initTmpl := message.WithTemplateMessage(&message.Template{
		Name: GlobalSvc.conf.VerificationCodeTmplName,
		Language: &message.TemplateLanguage{
			Code: GlobalSvc.conf.VerificationCodeLanguage.ToWhatsAppLang(),
		},
	})

	initTmplMessage, err := message.New(recipient, initTmpl)
	if err != nil {
		logx.Errorf("error creating initial template message: %v\n", err)
		return err
	}

	response, err := GlobalSvc.baseClient.SendMessage(ctx, initTmplMessage)
	if err != nil {
		logx.Errorf("error sending initial template message: %v\n", err)
		return err
	}

	if !response.Success {
		logx.Errorf("send message failed. rsp:%#v\n", response)
		return errors.New("send message failed. rsp:%#v")
	}

	return nil
}
