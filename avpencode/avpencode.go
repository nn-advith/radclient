package avpencode

import (
	"log"
	"net"
)

func EncodeIP(ip string) []byte {
	pip := net.ParseIP(ip)
	if pip == nil {
		log.Fatalf("unable to parse IP")
	}
	return pip.To4()
}
