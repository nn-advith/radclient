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

	newAR := requests.NewAccessRequest()
	// fmt.Printf("%x\n", newAR.Encode())
	conn.Write(newAR.Encode())
	// fmt.Fprintf(conn, "\x01\x42\x00\x36\x3f\x4e\x6a\x1d\x9f\x8b\x0a\x43\x12\x56\x78\x9a\xbc\xde\xf0\x12\x01\x0a\x73\x6f\x6d\x65\x75\x73\x65\x72\x02\x12\x7b\x2b\x7f\x20\x23\x7a\xf8\xac\x02\x87\xa5\x73\x34\xd5\x79\x7c\x04\x06\xc0\xa8\x38\x00")
	n, err := bufio.NewReader(conn).Read(p)
	if err == nil {
		fmt.Printf("raw: %v\n", p[:n])
		fmt.Printf("string: %x\n", p[:n]) // manually decode
	} else {
		fmt.Printf("Some error %v\n", err)
	}
	conn.Close()

	// op := avpencode.EncodeIP("192.168.56.10")
	// fmt.Printf("%x\n", op)

}
