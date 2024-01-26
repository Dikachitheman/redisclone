package main

import (
	// "errors"
	"fmt"
	// "io"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func goid() int {

	buf := make([]byte, 32)
	n := runtime.Stack(buf[:], false)

	// fmt.Println(n)
	// fmt.Println(buf)
	// fmt.Println(string(buf[:n]))
	// fmt.Println(strings.TrimPrefix(string(buf[:n]), "goroutine "))

	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]

	// fmt.Println("w", idField)

	id, err := strconv.Atoi(idField)

	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}

	return id

}

func handleRequest(conn net.Conn) {

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)

	if err != nil {
		log.Fatal(err)
	}

	if n != 0 {
		fmt.Println("\nreceived", string(buffer[:n]))
	}

	conn.Write([]byte("xnxx"))
	fmt.Println("xoxo")
	fmt.Println("go: ", goid())
}

func main() {

	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	} else {
		fmt.Println("listening...")
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleRequest(conn)
	}

}