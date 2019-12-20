package mail

import (
	"gopkg.in/gomail.v2"
)

type Client struct {
	Host     string // 地址
	Port     int    // 端口
	User     string // 邮箱账号
	Password string // 邮箱密码
	To       string // 收件人
	Subject  string // 主题
}

func NewClient(host, user, password string, port int) *Client {
	return &Client{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
	}
}

func (c *Client) Send(to, subject, content string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", c.User)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)
	d := gomail.NewDialer(c.Host, c.Port, c.User, c.Password)
	err := d.DialAndSend(m)
	return err
}
