package main

import (
	"bufio"
	"fmt"
	"log"
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
		"User-Password":  "w",
		"NAS-IP-Address": "192.168.56.10",
		"NAS-Port":       "1812",
	}

	// IMPORTANT: for test only;
	secret := "radius"

	n := 1

	go func() {
		p := make([]byte, 2048)
		for {
			n, err := bufio.NewReader(conn).Read(p)
			if err == nil {
				fmt.Printf("%x\n", p[0])
				switch p[0] {
				case 2:
					fmt.Printf("\033[032mAccess-Accept\033[0m\n")
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
						pendingReqs.Delete(key)
					} else {
						// porbably not some request which we sent; idk how we got; ignore
						fmt.Println("dud")
					}
				case 3:
					fmt.Printf("\033[031mAccess-Reject\033[0m\n")
					AA := requests.AccessReject{}
					AA.Decode(p[:n])
					key := AA.Identifier
					if v, ok := pendingReqs.Load(key); ok {
						valid := requests.ValidateAccessReject(p[:n], v.([16]byte), secret)
						if valid {
							fmt.Printf("READ: Response Authenticator check \033[032mPASSED\033[0m\n")
						} else {
							fmt.Printf("READ: Response Authenticator check \033[031mFAILED\033[0m\n")
						}
						pendingReqs.Delete(key)
					} else {
						// porbably not some request which we sent; idk how we got; ignore
						fmt.Println("dud")
					}
				default:
					log.Println("Unknown message code")
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

	time.Sleep(3 * time.Second)
	// add decoding logic; dynamicaaly depending on which type of response is received.
	// first check code and then determine whcih struct to decode into

}
