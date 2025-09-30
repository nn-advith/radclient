package avpencode

import (
	"log"
	"net"
)

var AVPEncodMap = map[uint8]interface{}{
	1: EncodeString,
	// 2: EncodePassword,
	4: EncodeIP,
	5: EncodeString,
}

func EncodeString(s string) []uint8 {
	return []byte(s)
}

func EncodeIP(ip string) []uint8 {
	pip := net.ParseIP(ip)
	if pip == nil {
		log.Fatalf("unable to parse IP")
	}
	return pip.To4()
}

// arg0: password, arg1: requestauth, arg2: secret
// func EncodePassword(password string, auth [16]uint8, secret string) []uint8{
// 	//calculate here

// }
