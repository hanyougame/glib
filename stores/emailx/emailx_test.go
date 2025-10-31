package emailx

import (
	"context"
	"fmt"
	"testing"

	"github.com/hanyougame/glib/stores/emailx/config"
)

func TestEmailSend(t *testing.T) {
	cfg := config.Config{
		Debug:     false,
		Trace:     false,
		SendEmail: "OG@mail.og.game",
		SendName:  "hanyouweb system",
		ApiUser:   "JGvNzgSDkQAE_test_MhOXw1",
		ApiKey:    "09642b32c2e54d4fcd8ef9e0b9469bcf",
	}

	client := NewEmailClient(cfg)
	err := client.SendByEngageLab(context.Background(), "邮件标题", "你的验证码是000000", "leo@ogmail.org")
	if err != nil {
		fmt.Println("Send email error:", err)
	} else {
		fmt.Println("Send mail success!")
	}
	return
}
