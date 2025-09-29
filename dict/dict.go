package dict

var PacketType = map[string]int{
	"AccessRequest":   1,
	"AccessAccept":    2,
	"AccessReject":    3,
	"AccessChallenge": 11,
}

var AVP = map[string]int{
	"User-Name":      1,
	"User-Password":  2,
	"NAS-IP-Address": 4,
	"NAS-Port":       5,
}
