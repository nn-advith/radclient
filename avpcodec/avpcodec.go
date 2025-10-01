package avpcodec

import (
	"crypto/md5"
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

// for int32 i.e avps that are integers with 4 octet value fields. write some other for dynamic ints ig
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
	// length At least 18 and no larger than 130.
	// value The String field is between 16 and 128 octets long, inclusive.

	// procedure
	// check if length is multiple of 16 octet. if no, pad to next multiple and divide to blocks
	// sequentially calculate
	// c1 = p1 XOR MD5(secret + auth)
	// ci = pi XOR MD5(secret + ci-1)
	// return c1+c2+...+ci (concatenated)
	ptext := []uint8(value.(string))

	if len(ptext)%16 != 0 {
		// pad here
		ptext = append(ptext, make([]uint8, 16-len(ptext)%16)...)
	}
	opstr := []uint8{}
	for l := 0; l < len(ptext); l += 16 {
		// fmt.Printf("%x-", ptext[l:l+16])
		var temp []uint8
		if l == 0 {
			temp = append([]uint8(ctx.Secret), []uint8(ctx.Authenticator[:])...)
		} else {
			temp = append([]uint8(ctx.Secret), opstr[l-16:l]...)
		}
		md5op := md5.Sum(temp)
		xorop := make([]uint8, 16)
		for i := 0; i < 16; i++ {
			xorop[i] = ptext[l+i] ^ md5op[i]
		}
		opstr = append(opstr, xorop...)
	}
	return opstr

}

// add avp decode map
