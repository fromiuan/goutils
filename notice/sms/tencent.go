package sms

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const TencentSMSURL = "https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=?&random=?"

type SMSTencent struct {
	Debug bool
	SMSTencentOpetion
}

type SMSTencentOpetion struct {
	SecretId  string
	SecretKey string
	Sign      string
}

type tencentRequest struct {
	Ext    string   `json:"ext"`
	Extend string   `json:"extend"`
	Params []string `json:"params"`
	Sig    string   `json:"sig"`
	Sign   string   `json:"sign"`
	Tel    struct {
		Mobile     string `json:"mobile"`
		Nationcode string `json:"nationcode"`
	} `json:"tel"`
	Time  int64  `json:"time"`
	TplID string `json:"tpl_id"`
}

type tencentRespone struct {
	Result int    `json:"result"`
	Errmsg string `json:"errmsg"`
	Ext    string `json:"ext"`
	Fee    int    `json:"fee"`
	Sid    string `json:"sid"`
}

func (t *SMSTencent) StartAndGC(config interface{}) error {
	option, ok := config.(SMSTencentOpetion)
	if !ok {
		return errors.New("sms:cfg must SMSTencentOpetion struct")
	}
	t.SMSTencentOpetion = option
	return nil
}

func (t *SMSTencent) SetDebug(b bool) {
	t.Debug = true
}

func (t *SMSTencent) Send(mobile, templeteId string, params interface{}) error {
	if t.Debug {
		return nil
	}

	paramsData, ok := params.([]string)
	if !ok {
		return errors.New("sms:params must type []string")
	}

	random := t.randomStr(12)
	jsonByte, err := t.newTencentRequest(mobile, templeteId, random, paramsData)
	if err != nil {
		return errors.New("生成json字符串错误")
	}

	err = t.post(random, jsonByte)
	if err != nil {
		return err
	}
	return nil
}

func (t *SMSTencent) post(random string, jsonByte []byte) error {
	req, err := http.NewRequest("POST", fmt.Sprintf(TencentSMSURL, t.SecretId, random), bytes.NewBuffer(jsonByte))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var rsp tencentRespone
	xml.Unmarshal(respBytes, &rsp)
	//处理return code.
	if rsp.Result != 0 {
		return errors.New("短信发送失败，原因:" + rsp.Errmsg)
	}
	return nil
}

func (t *SMSTencent) newTencentRequest(mobile, templeteId, random string, params []string) ([]byte, error) {
	now := time.Now().Unix()
	var m tencentRequest
	m.Ext = ""
	m.Extend = ""
	m.Params = params
	m.Sig = t.smsCalcSign("appkey=" + t.SecretId + "&random=" + random + "&time=" + strconv.FormatInt(now, 10) + "&mobile=" + mobile)
	m.Sign = t.Sign
	m.Tel.Mobile = mobile
	m.Tel.Nationcode = "86"
	m.Time = now
	m.TplID = templeteId
	return json.Marshal(m)
}

func (t *SMSTencent) smsCalcSign(m string) string {
	h := sha256.New()
	h.Write([]byte(m))

	return fmt.Sprintf("%x", h.Sum(nil))
}

//RandomStr 随机生成字符串
func (t *SMSTencent) randomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func init() {
	Register(SMS_TENCENT, &SMSTencent{})
}
