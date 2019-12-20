package sms

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const AliyunSMSURL = "http://dysmsapi.aliyuncs.com/"

type SMSAliyun struct {
	Debug bool
	SMSAliyunOption
}

type sendSmsResponse struct {
	Message   string
	RequestId string
	BizId     string
	Code      string
}

type SMSAliyunOption struct {
	AccessKeyId  string
	AccessSecret string
	SignName     string
}

func (a *SMSAliyun) StartAndGC(cfg interface{}) error {
	option, ok := cfg.(SMSAliyunOption)
	if !ok {
		return errors.New("sms:cfg must SMSAliyunOption struct")
	}
	a.SMSAliyunOption = option
	return nil
}

func (a *SMSAliyun) SetDebug(b bool) {
	a.Debug = b
}

func (a *SMSAliyun) Send(mobile string, templateCode string, params interface{}) error {
	if a.Debug {
		return nil
	}
	_, ok := params.(map[string]string)
	if !ok {
		return errors.New("sms:params must type map[string]string")
	}
	jsn, _ := json.Marshal(params)
	url := a.getParamUrl(mobile, string(jsn), templateCode)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return err
	}
	resp.Body.Close()

	var result sendSmsResponse
	err = xml.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	if result.Code != "OK" {
		return errors.New(result.Message)
	}
	return nil
}

func (a *SMSAliyun) getParamUrl(mobile string, message string, templateCode string) string {

	var params = map[string]string{
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureNonce":   fmt.Sprintf("%d", rand.Int63()),
		"AccessKeyId":      a.AccessKeyId,
		"SignatureVersion": "1.0",
		"Timestamp":        a.getGMTTime(),
		"Format":           "",
		"Action":           "SendSms",
		"Version":          "2017-05-25",
		"RegionId":         "cn-hangzhou",
		"PhoneNumbers":     mobile,
		"SignName":         a.SignName,
		"TemplateParam":    message,
		"TemplateCode":     templateCode,
		"OutId":            "123",
	}

	var sortKeys []string
	for k, _ := range params {
		sortKeys = append(sortKeys, k)
	}

	sort.Strings(sortKeys)

	var str string
	for _, key := range sortKeys {
		str += "&"
		str += a.urlEncode(key)
		str += "="
		str += a.urlEncode(params[key])
	}

	encryptKey := fmt.Sprintf("%s%s%s%s%s", "GET", "&", a.urlEncode("/"), "&", a.urlEncode(str[1:]))
	encryptText := a.AccessSecret + "&"

	return AliyunSMSURL + a.urlEncode(a.hmacSHA1(encryptText, encryptKey)) + str
}

func (a *SMSAliyun) getGMTTime() string {
	now := time.Now()
	year, mon, day := now.UTC().Date()
	hour, min, sec := now.UTC().Clock()
	utcTime := fmt.Sprintf(`%d-%02d-%02dT%02d:%02d:%02dZ`, year, mon, day, hour, min, sec)
	return utcTime
}

func (a *SMSAliyun) urlEncode(str string) string {
	result := url.QueryEscape(str)
	result = strings.Replace(result, "+", "%20", -1)
	result = strings.Replace(result, "*", "%2A", -1)
	result = strings.Replace(result, "%7E", "~", -1)
	return result
}

func (a *SMSAliyun) hmacSHA1(encryptText, encryptKey string) string {
	key := []byte(encryptText)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(encryptKey))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func init() {
	Register(SMS_ALIYUN, &SMSAliyun{})
}
