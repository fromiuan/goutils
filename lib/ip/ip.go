package ip

import (
	"github.com/17mon/go-ipip/ipip"
)

var districtDB *ipip.DistrictDB

func init() {
	districtDB = ipip.NewDistrictDB()
	if err := districtDB.Load("ip.data"); err != nil {
		panic("load ip data error")
	}
}

func GetIpLocation(ip string) (ipip.District, error) {
	return districtDB.Find(ip)
}
