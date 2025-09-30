package requests

import (
	"github.com/nn-advith/radclient/dict"
	"github.com/nn-advith/radclient/utils"
)

type Request interface {
	Encode() []uint8
	Decode() []uint8
	Validate() []uint8 // temporary for now; basically to enfore the rules on type of packet as per rfc
}

type Attribute interface {
	Validate() // for now
}

type AVP struct {
	Type   uint8
	Length uint8
	Value  []byte
}

type AccessRequest struct {
	Code          uint8
	Identifier    uint8
	Length        uint16
	Authenticator [16]uint8
	Attributes    []AVP
}

func NewAccessRequest(avps map[string]string, secret string) AccessRequest {
	// compute length

	// for k, v := range avps {
	// 	//construct temp avp and encode it, append to the attrubutes;
	// 	// OPTIMISE for memory, this is probably not needed
	// 	op:= avpencode.AVPEncodMap[uint8(dict.AVP[k])].
	// 	tempavp := AVP{
	// 		 Type: uint8(dict.AVP[k]),
	// 		 Length: uint8(0),
	// 		 Value: op,
	// 	}
	// }

	tempreq := AccessRequest{
		Code:          uint8(dict.PacketType["AccessRequest"]),
		Identifier:    utils.GenerateIdentifier(),
		Length:        uint16(0),
		Authenticator: [16]uint8(utils.GenerateAuthenticator()),
		Attributes:    []AVP{},
	}
	return tempreq
}

func (a *AccessRequest) Encode() []uint8 {
	// convert code to uint8
	// generate identifier
	// temprory length set to 00
	// randomly generate authenticator ( for now rad0 has msg-authenticator disabled so use normal authenticator)
	// encode each avp
	// combine and return

	// networek byte order i.e big endian ( virtually all protocols use this; refer investigations)
	res := append([]uint8{a.Code, a.Identifier, uint8(a.Length >> 8), uint8(a.Length)}, a.Authenticator[0:len(a.Authenticator)]...)
	return res
}
