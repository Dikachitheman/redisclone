package main
import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)
//*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n
const (
	Array = '*'
	Bulk  = '$'
)
func decode(r *bufio.Reader) ([]string, error) {
	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	switch b {
	case Bulk:
		vbarr, err := validBytes(r)
		if err != nil {
			return nil, err
		}
		count, err := strconv.Atoi(string(vbarr))
		if err != nil {
			return nil, err
		}
		buf := make([]byte, count+2)
		_, err = io.ReadFull(r, buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		return []string{string(bytes.TrimSpace(buf))}, nil
	default:
		return nil, fmt.Errorf("could not decode the stream")
	}
}
func validBytes(r *bufio.Reader) ([]byte, error) {
	barr, err := r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return barr[:len(barr)-2], nil
}
func decodeArray(r *bufio.Reader) ([]string, error) {
	vbarr, err := validBytes(r)
	if err != nil {
		return nil, err
	}
	l, err := strconv.Atoi(string(vbarr))
	if err != nil {
		return nil, err
	}
	cmds := []string{}
	for i := 0; i < l; i++ {
		vals, err := decode(r)
		if err == io.EOF {
			return vals, nil
		}
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, vals...)
	}
	return cmds, nil
}
type input struct {
	raw  []byte
	cmds []string
}
func (i *input) parse() {
	r := bufio.NewReader(bytes.NewReader(i.raw))
	b, err := r.ReadByte()
	if err != nil {
		log.Fatal("error reading 1st byte: ", err)
	}
	switch b {
	case Array:
		cmds, err := decodeArray(r)
		if err != nil {
			log.Fatal("error parsing input stream: ", err)
		}
		i.cmds = cmds
	}

}
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		barr := make([]byte, 1024)
		_, err := conn.Read(barr)

		if err == io.EOF {
			log.Println("client is done")
			return
		}
		if err != nil {
			
			log.Fatal(err)
		}
		in := input{raw: barr}
		in.parse()
		switch strings.ToUpper(in.cmds[0]) {
		case "PING":
			mustCopy(conn, strings.NewReader("+PONG\r\n"))
		case "ECHO":

			mustCopy(conn, strings.NewReader(fmt.Sprintf("%s%s%s", "+", in.cmds[1], "\r\n")))
		default:
			mustCopy(conn, strings.NewReader("+PONG\r\n"))
		}
	}
}
func main() {
	l, err := net.Listen("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}