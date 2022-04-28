package api

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type IPRange struct {
	cidrs []*net.IPNet
}

func IPRangeFromExternalSrc(url string) (IPRange, error) {

	resp, err := http.Get(url)

	if err != nil {
		return IPRange{}, err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	cidrs := strings.Split(string(content), "\r\n")
	subnets := []*net.IPNet{}

	for _, cidr := range cidrs {
		_, subnet, err := net.ParseCIDR(cidr)

		if err != nil {
			logrus.Warn("Failed to parse CIDR: ", err)
			continue
		}

		subnets = append(subnets, subnet)
	}

	return IPRange{
		cidrs: subnets,
	}, nil

}

func (r *IPRange) Contains(clientip string) bool {
	ip := net.ParseIP(clientip)

	for _, cidr := range r.cidrs {
		if cidr.Contains(ip) {
			return true
		}
	}

	return false
}

func (r *IPRange) Size() int {
	return len(r.cidrs)
}
