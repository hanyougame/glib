package emailx

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/hanyougame/glib/stores/emailx/config"
	"log"
	"net"
	"net/smtp"
)

var Email EmailClient

//发送邮件

type EmailClient struct {
	smtpServer string // smtp服务器地址
	smtpPort   int    // smtp服务器地址
	password   string // 发送密码
	sendEmail  string // 发送地址
	sendName   string // 发送人名称
}

func Must(c config.Config) {
	Email = NewEmailClient(c)
}

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
func (e EmailClient) Send(subject, body string, toEmail ...string) error {
	if len(toEmail) == 0 || len(body) == 0 {
		return errors.New("address or body is empty")
	}
	header := make(map[string]string)
	header["From"] = e.sendName + " <" + e.sendEmail + ">"
	//header["To"] = toEmail[0]
	header["Subject"] = subject
	//html格式邮件
	header["Content-Type"] = "text/html; charset=UTF-8"
	//纯文本格式邮件
	//header["Content-Type"] = "text/plain; charset=UTF-8"
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body
	auth := smtp.PlainAuth(
		"",
		e.sendEmail,
		e.password,
		e.smtpServer,
	)
	return sendMailWithTLS(
		fmt.Sprintf("%s:%d", e.smtpServer, e.smtpPort),
		auth,
		e.sendEmail,
		toEmail,
		message,
	)
}

// Dial return a smtp client
func dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("tls.Dial Error:", err)
		return nil, err
	}

	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

// sendMailWithTLS send email with tls
func sendMailWithTLS(smtpAddr string, auth smtp.Auth, from string,
	to []string, msg string) (err error) {
	//create smtp client
	c, err := dial(smtpAddr)
	if err != nil {
		return errors.New("Create smtp client error:" + err.Error())
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return errors.New("SMTP Error during AUTH:" + err.Error())
			}
		}
	}
	for _, addr := range to {
		newMsg := fmt.Sprintf("%s: %s\r\n", "To", addr) + msg
		if err = c.Mail(from); err != nil {
			return err
		}
		if err = c.Rcpt(addr); err != nil {
			return err
		}
		w, err := c.Data()
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(newMsg))
		if err != nil {
			return err
		}
		err = w.Close()
		if err != nil {
			return err
		}
	}
	return c.Quit()
}
