package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
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

func mapRespCommand(buffer []byte) map[int]string {

	sBuffer := string(buffer)
	respMap := make(map[int]string)

	numStrings, _ := strconv.Atoi(sBuffer[1:2])

	// fmt.Println(numStrings)

	arrayArray := make([]string, numStrings)

	sliceArray := sBuffer[numStrings:]

	// fmt.Println(sliceArray)

	for i := 0; i < numStrings; i++ {

		a := 0
		b := 0

		for i := 0; i < len(sliceArray); i++ {
			if string(sliceArray[i]) == "$" {
				a = i + 1
				b = i

				// fmt.Println("a ", a)
				// fmt.Println("b ", b)
				break
			}
		}

		cutNumber, _ := strconv.Atoi(sliceArray[b+1 : a+1])
		// fmt.Println("cutnumber ", cutNumber)

		arrayArray[i] = sliceArray[a+3 : cutNumber+a+3]

		newSliceArray := sliceArray[cutNumber+a+3+b:]

		sliceArray = newSliceArray

		// fmt.Println("slice ", sliceArray)
		// fmt.Println("arr ", arrayArray[i])

		respMap[i] = arrayArray[i]

	}

	return respMap

}

func echo(conn net.Conn, message string) {

	response := fmt.Sprintf(message)
	conn.Write([]byte(response))

}

func errorFunc(conn net.Conn) {

	response := fmt.Sprintf("error")
	conn.Write([]byte(response))
}

func handleRequest(conn net.Conn) {

	buffer := make([]byte, 1024)
	respMap := make(map[int]string)

	commandMap := map[string]interface{}{

		"ECHO":  echo,
		"error": errorFunc,
	}

	for {

		n, err := conn.Read(buffer)

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("\nreceived", string(buffer[:n]))

		respMap = mapRespCommand(buffer[:n])

		for index, value := range commandMap {

			if respMap[0] == index {
				value.(func(net.Conn, string))(conn, respMap[1])
			}

			if index == "error" {
				value.(func(net.Conn))(conn)
			}

		}

	}

	fmt.Println("go: ", goid())
}

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
