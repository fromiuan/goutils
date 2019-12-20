package sms

import (
	"errors"
)

const (
	SMS_ALIYUN  = "aliyun"
	SMS_TENCENT = "tencent"
	SMS_HUAWEI  = "huawei"
)

type Smser interface {
	SetDebug(bool)
	Send(string, string, interface{}) error
	StartAndGC(config interface{}) error
}

var adapters = make(map[string]Smser)

func Register(name string, adapter Smser) {
	if adapter == nil {
		panic("sms:register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("sms:repleace registry sms name")
	}
	adapters[name] = adapter
}

func NewClient(adapterName string, config interface{}) (Smser, error) {
	adapter, ok := adapters[adapterName]
	if !ok {
		return adapter, errors.New("sms:not find sms name")
	}
	err := adapter.StartAndGC(config)
	if err != nil {
		return adapter, err
	}
	return adapter, nil
}
