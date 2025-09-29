package utils

import (
	crand "crypto/rand"
	"log"
	mrand "math/rand"
	"net"
)

func EncodeIP(ip string) []byte {
	pip := net.ParseIP(ip)
	if pip == nil {
		log.Fatalf("unable to parse IP")
	}
	return pip.To4()
}

func GenerateIdentifier() byte {
	return byte(mrand.Intn(256))
}

func GenerateAuthenticator() []uint8 {
	auth := make([]uint8, 16)
	crand.Read(auth)
	return auth
}
