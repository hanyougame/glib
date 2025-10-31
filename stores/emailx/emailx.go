package emailx

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/hanyougame/glib/stores/emailx/config"
	"github.com/hanyougame/glib/utils/httpc"
)

var Email EmailClient

// EmailClient 结构体
// 负责管理 SMTP 邮件发送

type EmailClient struct {
	Config config.Config
}

// Must 初始化 EmailClient
func Must(c config.Config) {
	Email = NewEmailClient(c)
}

// NewEmailClient 创建 EmailClient 实例
func NewEmailClient(c config.Config) EmailClient {
	return EmailClient{Config: c}
}

// Send 发送邮件(注:测试110个邮件,发送大概50s)
func (e EmailClient) Send(subject, body string, toEmails ...string) error {
	if len(toEmails) == 0 || len(body) == 0 {
		return errors.New("recipient address or body is empty")
	}

	headers := map[string]string{
		"From":         fmt.Sprintf("%s <%s>", e.Config.SendName, e.Config.SendEmail),
		"Subject":      subject,
		"Content-Type": "text/html; charset=UTF-8",
	}

	var msgBuilder strings.Builder
	for k, v := range headers {
		msgBuilder.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msgBuilder.WriteString("\r\n" + body)

	auth := smtp.PlainAuth("", e.Config.SendEmail, e.Config.Password, e.Config.SmtpServer)

	return e.sendMailWithTLS(auth, toEmails, msgBuilder.String())
}

// sendMailWithTLS 通过 TLS 发送邮件
func (e EmailClient) sendMailWithTLS(auth smtp.Auth, recipients []string, message string) error {
	addr := fmt.Sprintf("%s:%d", e.Config.SmtpServer, e.Config.SmtpPort)
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return fmt.Errorf("failed to establish TLS connection: %w", err)
	}
	defer conn.Close()

	host, _, _ := net.SplitHostPort(addr)
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth error: %w", err)
	}

	if err = client.Mail(e.Config.SendEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	for _, recipient := range recipients {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
		writer, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to start data transfer: %w", err)
		}
		if _, err = writer.Write([]byte("To: " + recipient + "\r\n" + message)); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}
		if err = writer.Close(); err != nil {
			return fmt.Errorf("failed to finalize message: %w", err)
		}
	}
	return nil
}

func (e EmailClient) SendByEngageLab(ctx context.Context, subject, body string, toEmails ...string) error {
	originalBytes := []byte(e.Config.ApiUser + ":" + e.Config.ApiKey)
	auth := "Basic " + base64.StdEncoding.EncodeToString(originalBytes)
	url := "https://email.api.engagelab.cc/v1/mail/send"

	resp, err := httpc.Do(ctx).SetBody(map[string]any{
		"from":    e.Config.SendEmail,
		"to":      toEmails,
		"subject": subject,
		"body": map[string]any{
			"subject": subject,
			"content": map[string]string{
				"html": body,
			},
		},
	}).
		SetHeaders(map[string]string{
			"Authorization": auth,
		}).
		Post(url)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to send email: %s", resp.String())
	}
	return nil
}
