package main

import (
	"fmt"
	"strings"
	"sync"
	// "bytes"
	// "errors"
	"runtime"
	"strconv"
)

func main() {

	wg := sync.WaitGroup{}
	done := make(chan struct{})

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		fmt.Println("this is goroutine: ", goid())
		done <- struct{}{}
		wg.Done()
	}(&wg)
	wg.Wait()

	fmt.Println("this is main, id: ", goid())

	<-done
}

func goid() (int) {

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