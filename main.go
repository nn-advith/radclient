package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/nn-advith/radclient/metrics"
	"github.com/nn-advith/radclient/requests"
)

// with go routines;
// multiple writer routine to send n requests
// one reader routine to read responses and possibly track
// track using sync.Map

var pendingReqs sync.Map // merge this into metrics tracker iguess; add a flag to save metrics if needed

func SendRequest(id uint8, conn net.Conn, avps map[string]interface{}, secret string, mt *metrics.Metrics) {
	newAR := requests.NewAccessRequest(id, avps, secret)
	encodedpacket := newAR.Encode()
	key := newAR.Identifier
	pendingReqs.Store(key, newAR.Authenticator)
	mt.Start(newAR.Identifier)
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
	n := 50
	printMetrics := true // get from config

	closeChannel := make(chan struct{})

	// metics starts
	var mtracker *metrics.Metrics
	if printMetrics {
		mtracker = metrics.NewMetrics()
		go mtracker.StartCollection()
	}

	go func() {
		p := make([]byte, 2048)
		for {
			n, err := bufio.NewReader(conn).Read(p)
			if err == nil {
				if printMetrics {
					mtracker.End(p[1])
				}
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

	// closer
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)

			empty := true
			pendingReqs.Range(func(key, value any) bool {
				// if it ranges then set empty to false
				empty = false
				return false
			})
			if empty && mtracker.IsEmpty() {
				closeChannel <- struct{}{}
			}
		}
	}()

	for i := range n {
		id := uint8(i)
		go SendRequest(id, conn, avps, secret, mtracker)
	}

	<-closeChannel
	os.Exit(0)

	// add decoding logic; dynamicaaly depending on which type of response is received.
	// first check code and then determine whcih struct to decode into

}
