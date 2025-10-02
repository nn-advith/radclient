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

	// add feature to track pending requests; in a different routine i guess
	pending := newAR.Identifier
	fmt.Printf("%x\n", pending)

	conn.Write(encodedpacket)

	// add decoding logic; dynamicaaly depending on which type of response is received.
	// first check code and then determine whcih struct to decode into
	AA := requests.AccessAccept{}

	n, err := bufio.NewReader(conn).Read(p)
	if err == nil {
		fmt.Printf("string: %x\n", p[:n]) // manually decode
		// decode
		AA.Decode(p[:n])
		if AA.Identifier == pending {
			fmt.Println("Received response for pending request")
		}
		valid := requests.ValidateAccessAccept(p[:n], newAR.Authenticator, secret)
		if valid {
			fmt.Println("valid response")
		} else {
			fmt.Println("not a valid response; response authenticator check failed")
		}

	} else {
		fmt.Printf("error %v\n", err)
	}
	conn.Close()

}
