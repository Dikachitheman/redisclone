package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	must(err)
	defer l.Close()
	fmt.Println("Redis Server is running on port 6379...")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConn(conn)
	}
}
func handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)

for {
	n, err := conn.Read(buf)
	if err != nil && err == io.EOF {
		fmt.Println("Client is done")
		break
	}
	if n != 0 {
		fmt.Println("Received: ", buf[:n])
		// fmt.Println("Received: ", n)
		fmt.Println("Received: \r\n", string(buf[:n]))
		sendPong(conn)
	} else {
		break
	}
}
}
func sendPong(conn net.Conn) {
resp := []byte("+PONG\r\n")

fmt.Println("Going to PONG")

_, err := conn.Write(resp)
must(err)
}

func must(err error) {
if err != nil {
	fmt.Println(err)
}
}