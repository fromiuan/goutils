package dingtalk

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	MessageTypeText     = "text"
	MessageTypeMarkdown = "markdown"

	DingTalkApi = "https://oapi.dingtalk.com/robot/send?access_token=?"
)

type Client struct {
	RobotURL string
}

func NewClient(token string) *Client {
	return &Client{
		RobotURL: fmt.Sprintf(DingTalkApi, token),
	}
}

func (d *Client) SendMessage(title string, msgtype string, content string) error {

	var message string

	switch msgtype {
	case MessageTypeText:
		message = fmt.Sprintf(`{"msgtype": "text", "text": {"content": "%s"}}`, content)
	case MessageTypeMarkdown:
		message = fmt.Sprintf(`{"msgtype": "markdown", "markdown": {"title": "%s", "text": "%s"}}`, title, content)
	default:
		message = fmt.Sprintf(`{"msgtype": "text", "text": {"content": "%s"}}`, content)
	}

	client := &http.Client{}
	request, _ := http.NewRequest("POST", d.RobotURL, bytes.NewBuffer([]byte(message)))
	request.Header.Set("Content-type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("dingtalk: Send Message | %v", err)
	}

	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("dingtalk: Send Message | %v", string(body))
	}

	_, err = ioutil.ReadAll(response.Body)
	return err
}
