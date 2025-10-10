package requests

import (
	"crypto/md5"
	"fmt"
	"log"

	"github.com/nn-advith/radclient/avpcodec"
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

func (a *AVP) CalculateLength() uint8 {
	return uint8(len(a.Value) + 2)
}

func (a *AVP) Stream() []uint8 {
	return append([]uint8{a.Type, a.Length}, a.Value...)
}

// ACCESS-REQUEST

type AccessRequest struct {
	Code          uint8
	Identifier    uint8
	Length        uint16
	Authenticator [16]uint8
	Attributes    []AVP
}

func NewAccessRequest(id uint8, avps map[string]interface{}, secret string) AccessRequest {
	// compute length at end
	tempauth := [16]uint8(utils.GenerateAuthenticator())
	tempreq := AccessRequest{
		Code: uint8(dict.PacketType["AccessRequest"]),
		// Identifier:    utils.GenerateIdentifier(),
		Identifier:    id,
		Length:        uint16(0),
		Authenticator: tempauth,
		Attributes:    []AVP{},
	}

	ctx := avpcodec.AVPContext{
		Secret:        secret,
		Authenticator: tempauth,
	}
	for k, v := range avps {
		//construct temp avp and encode it, append to the attrubutes;
		// OPTIMISE for memory, this is probably not needed
		f := avpcodec.AVPEncodeMap[uint8(dict.AVP[k])]
		tempavp := AVP{
			Type:   uint8(dict.AVP[k]),
			Length: uint8(0),
			Value:  f(v, ctx),
		}
		tempavp.Length = tempavp.CalculateLength()
		tempreq.Attributes = append(tempreq.Attributes, tempavp)
	}

	return tempreq
}

func (a *AccessRequest) CalculateLength() uint16 {
	l := 20
	for i := range a.Attributes {
		l += int(a.Attributes[i].CalculateLength())
	}
	return uint16(l)
}

func (a *AccessRequest) Encode() []uint8 {
	// convert code to uint8
	// generate identifier
	// temprory length set to 00
	// randomly generate authenticator ( for now rad0 has msg-authenticator disabled so use normal authenticator)
	// encode each avp
	// combine and return

	// networek byte order i.e big endian ( virtually all protocols use this; refer investigations)
	// fmt.Printf("INFO: %x-%x-%x-%x\n", a.Code, a.Identifier, uint16(a.Length), a.Authenticator[:])
	//length update
	a.Length = a.CalculateLength()
	res := append([]uint8{a.Code, a.Identifier, uint8(a.Length >> 8), uint8(a.Length)}, a.Authenticator[0:len(a.Authenticator)]...)
	for i := range a.Attributes {
		res = append(res, a.Attributes[i].Stream()...)
	}
	return res
}

// ACCESS-ACCEPT

type AccessAccept struct {
	Code          uint8
	Identifier    uint8
	Length        uint16
	Authenticator [16]uint8
	Attributes    []AVP // howwwww to decode - sequential
}

func (a *AccessAccept) Decode(packet []uint8) {
	// decode packet and populate the AccessAccept struct
	a.Code = packet[0]
	a.Identifier = packet[1]
	a.Length = uint16(packet[2])>>8 + uint16(packet[3])
	a.Authenticator = [16]uint8(packet[4:20])

	fmt.Printf("attributes from avp: %x\n", packet[20:])
	// do opinionated parsing here; ignore AVPs you dont care about

}

func ValidateAccessAccept(packet []uint8, requestauth [16]uint8, secret string) bool {
	if uint8(packet[0]) != 2 {
		// not access accept
		log.Fatal("not access accept")
		return false
	}
	temp := make([]uint8, len(packet))
	copy(temp, packet)

	respauth := make([]uint8, 16)
	copy(respauth, temp[4:20])
	copy(temp[4:20], requestauth[:])
	temp = append(temp, []uint8(secret)...)
	md5op := md5.Sum(temp)
	// fmt.Printf("TEMP: %x\nRESPAUTH: %x\nMD5: %x\n", temp[:], respauth[:], md5op)
	return md5op == [16]byte(respauth)
}

// ACCESS-REJECT

type AccessReject struct {
	Code          uint8
	Identifier    uint8
	Length        uint16
	Authenticator [16]uint8
	Attributes    []AVP // if needed
}

func (a *AccessReject) Decode(packet []uint8) {
	a.Code = packet[0]
	a.Identifier = packet[1]
	a.Length = uint16(packet[2])>>8 + uint16(packet[3])
	a.Authenticator = [16]uint8(packet[4:20])

	fmt.Printf("attributes from avp: %x\n", packet[20:])
	// do opinionated parsing here; ignore AVPs you dont care about
}

func ValidateAccessReject(packet []uint8, requestauth [16]uint8, secret string) bool {
	if uint8(packet[0]) != 3 {
		// not access accept
		log.Fatal("not access reject")
		return false
	}
	temp := make([]uint8, len(packet))
	copy(temp, packet)

	respauth := make([]uint8, 16)
	copy(respauth, temp[4:20])
	copy(temp[4:20], requestauth[:])
	temp = append(temp, []uint8(secret)...)
	md5op := md5.Sum(temp)
	// fmt.Printf("TEMP: %x\nRESPAUTH: %x\nMD5: %x\n", temp[:], respauth[:], md5op)
	return md5op == [16]byte(respauth)
}
