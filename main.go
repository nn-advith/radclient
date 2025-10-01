package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/nn-advith/radclient/requests"
)

func main() {
	p := make([]byte, 2048)
	conn, err := net.Dial("udp", "192.168.56.10:1812")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	avps := map[string]interface{}{
		"User-Name":      "someuser",
		"User-Password":  "somepass",
		"NAS-IP-Address": "192.168.56.10",
		"NAS-Port":       "1812",
	}

	// IMPORTANT: for test only;
	secret := "radius"

	newAR := requests.NewAccessRequest(avps, secret)
	encodedpacket := newAR.Encode()
	fmt.Printf("%x\n", encodedpacket[:])
	pending := newAR.Identifier
	fmt.Printf("%x\n", pending)
	conn.Write(encodedpacket)

	// add feature to track pending requests; in a different routine i guess

	// add decoding logic
	AA := requests.AccessAccept{}

	n, err := bufio.NewReader(conn).Read(p)
	if err == nil {
		fmt.Printf("string: %x\n", p[:n]) // manually decode
		// decode
		AA.Decode(p[:n])
		fmt.Println(AA)
	} else {
		fmt.Printf("error %v\n", err)
	}
	conn.Close()

}
