package emailx

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/hanyougame/glib/stores/emailx/config"
	"net"
	"net/smtp"
	"strings"
)

var Email EmailClient

// EmailClient 结构体
// 负责管理 SMTP 邮件发送

type EmailClient struct {
	smtpServer string
	smtpPort   int
	password   string
	sendEmail  string
	sendName   string
}

// Must 初始化 EmailClient
func Must(c config.Config) {
	Email = NewEmailClient(c)
}

// NewEmailClient 创建 EmailClient 实例
func NewEmailClient(c config.Config) EmailClient {
	return EmailClient{
		smtpServer: c.SmtpServer,
		smtpPort:   c.SmtpPort,
		password:   c.Password,
		sendEmail:  c.SendEmail,
		sendName:   c.SendName,
	}
}

// Send 发送邮件(注:测试110个邮件,发送大概50s)
func (e EmailClient) Send(subject, body string, toEmails ...string) error {
	if len(toEmails) == 0 || len(body) == 0 {
		return errors.New("recipient address or body is empty")
	}

	headers := map[string]string{
		"From":         fmt.Sprintf("%s <%s>", e.sendName, e.sendEmail),
		"Subject":      subject,
		"Content-Type": "text/html; charset=UTF-8",
	}

	var msgBuilder strings.Builder
	for k, v := range headers {
		msgBuilder.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msgBuilder.WriteString("\r\n" + body)

	auth := smtp.PlainAuth("", e.sendEmail, e.password, e.smtpServer)

	return e.sendMailWithTLS(auth, toEmails, msgBuilder.String())
}

// sendMailWithTLS 通过 TLS 发送邮件
func (e EmailClient) sendMailWithTLS(auth smtp.Auth, recipients []string, message string) error {
	addr := fmt.Sprintf("%s:%d", e.smtpServer, e.smtpPort)
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

	if err = client.Mail(e.sendEmail); err != nil {
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
