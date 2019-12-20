package tools

import (
	"net"
	"os"
)

// SystemMac 获取本机的MAC地址
func SystemMac() interface{} {
	interfaces, _ := net.Interfaces()
	mac := make([]interface{}, 0)
	for _, inter := range interfaces {
		m := inter.HardwareAddr //获取本机MAC地址
		mac = append(mac, m)
	}
	return mac
}

func SystemHostname() string {
	if name, err := os.Hostname(); err == nil {
		return name
	}
	return ""
}
