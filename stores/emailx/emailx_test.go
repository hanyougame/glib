package emailx

import (
	"fmt"
	"github.com/hanyougame/glib/stores/emailx/config"
	"testing"
)

func TestEmailSend(t *testing.T) {
	cfg := config.Config{
		Debug:      false,
		Trace:      false,
		SmtpServer: "sg-smtp.qcloudmail.com",
		SmtpPort:   465,
		Password:   "CMa1WE4nKvOiGbNb",
		SendEmail:  "system@hanyouweb.com",
		SendName:   "hanyouweb system",
	}

	Must(cfg)
	err := Email.Send("邮件标题", "你的验证码是000000", "xxxxxx@qq.com")
	if err != nil {
		fmt.Println("Send email error:", err)
	} else {
		fmt.Println("Send mail success!")
	}
	return
}
