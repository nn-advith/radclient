package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/nn-advith/radclient/requests"
)

// with go routines;
// multiple writer routine to send n requests
// one reader routine to read responses and possibly track
// track using sync.Map

var pendingReqs sync.Map

func SendRequest(id uint8, conn net.Conn, avps map[string]interface{}, secret string) {
	newAR := requests.NewAccessRequest(id, avps, secret)
	encodedpacket := newAR.Encode()
	key := newAR.Identifier
	pendingReqs.Store(key, newAR.Authenticator)
	conn.Write(encodedpacket)
}

func main() {

	conn, err := net.Dial("udp", "192.168.56.10:1812")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	defer conn.Close()

	avps := map[string]interface{}{
		"User-Name":      "someuser",
		"User-Password":  "somepass",
		"NAS-IP-Address": "192.168.56.10",
		"NAS-Port":       "1812",
	}

	// IMPORTANT: for test only;
	secret := "radius"

	// newAR := requests.NewAccessRequest(avps, secret)
	// encodedpacket := newAR.Encode()
	// fmt.Printf("%x\n", encodedpacket[:])

	// add feature to track pending requests; in a different routine i guess
	// pending := newAR.Identifier
	// fmt.Printf("%x\n", pending)

	// conn.Write(encodedpacket)

	n := 10

	go func() {
		p := make([]byte, 2048)
		for {
			n, err := bufio.NewReader(conn).Read(p)
			if err == nil {
				AA := requests.AccessAccept{}
				AA.Decode(p[:n])
				key := AA.Identifier
				if v, ok := pendingReqs.Load(key); ok {
					valid := requests.ValidateAccessAccept(p[:n], v.([16]byte), secret)
					if valid {
						fmt.Printf("READ: Response Authenticator check \033[032mPASSED\033[0m\n")
					} else {
						fmt.Printf("READ: Response Authenticator check \033[031mFAILED\033[0m\n")
					}
				} else {
					// porbably not some request which we sent; idk how we got; ignore
					fmt.Println("dud")
				}
			} else {
				fmt.Printf("error: %v\n", err)
			}
		}
	}()

	for i := range n {
		id := uint8(i)
		go SendRequest(id, conn, avps, secret)
	}

	time.Sleep(5 * time.Second)
	// add decoding logic; dynamicaaly depending on which type of response is received.
	// first check code and then determine whcih struct to decode into
	// AA := requests.AccessAccept{}

	// m, err := bufio.NewReader(conn).Read(p)
	// if err == nil {
	// 	fmt.Printf("string: %x\n", p[:n]) // manually decode
	// 	// decode
	// 	AA.Decode(p[:m])
	// 	// if AA.Identifier == pending {
	// 	// 	fmt.Println("Received response for pending request")
	// 	// }
	// 	// valid := requests.ValidateAccessAccept(p[:n], newAR.Authenticator, secret)
	// 	// if valid {
	// 	// 	fmt.Println("valid response")
	// 	// } else {
	// 	// 	fmt.Println("not a valid response; response authenticator check failed")
	// 	// }

	// } else {
	// 	fmt.Printf("error %v\n", err)
	// }

}
