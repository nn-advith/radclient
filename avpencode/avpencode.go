package avpencode

import (
	"log"
	"net"
	"strconv"
)

type AVPValue interface{}

// acts as additional context in encoding that may or may not be used.
type AVPContext struct {
	Secret        string
	Authenticator [16]uint8
}

// function type
type AVPEncoder func(AVPValue, AVPContext) []uint8

var AVPEncodeMap = map[uint8]AVPEncoder{
	1: EncodeString,
	2: EncodePassword,
	4: EncodeIP,
	5: EncodeInteger32, //expects string
}

func EncodeString(value AVPValue, _ AVPContext) []uint8 {
	s := value.(string)
	return []byte(s)
}

func EncodeIP(value AVPValue, _ AVPContext) []uint8 {
	ip := value.(string)
	pip := net.ParseIP(ip)
	if pip == nil {
		log.Fatalf("unable to parse IP")
	}
	return pip.To4()
}

func EncodeInteger32(value AVPValue, _ AVPContext) []uint8 {
	i, err := strconv.Atoi(value.(string))
	if err != nil {
		log.Fatalf("string conversion failed")
	}
	r := make([]byte, 4)
	r[0] = uint8(i >> 24)
	r[1] = uint8(i >> 16)
	r[2] = uint8(i >> 8)
	r[3] = uint8(i)
	return r
}

// arg0: password, arg1: requestauth, arg2: secret
func EncodePassword(value AVPValue, ctx AVPContext) []uint8 {
	//calculate here
	return []uint8{}

}
