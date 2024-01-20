package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

const (
	HOST = "localhost"
	PORT = "6379"
	TYPE = "tcp"
)

func main() {
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)

	if err != nil {
		fmt.Println("resolvetcpaddr failed:", err.Error())
		os.Exit(1)
	}


	for i := 0; i < 3; i++ {

		conn, err := net.DialTCP(TYPE, nil, tcpServer)
		if err != nil {
			fmt.Println("Dial failed:", err.Error())
			os.Exit(1)
		}
		_, err = conn.Write([]byte("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"))
		
		if err != nil {
			fmt.Println("write data failed: ", err.Error())
			os.Exit(1)
		} else {
			fmt.Println("sending: " )
		}

		received := make([]byte, 1024)
		_, err = conn.Read(received)
		if err != nil {
			fmt.Println("read data failed: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("recived message:", string(received))

		conn.Close()

		time.Sleep(2 * time.Second)
	}
}