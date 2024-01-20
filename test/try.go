package main

import (
	"fmt"
	"strconv"
)

func main() {

	buffer := []byte("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n")

	sBuffer := string(buffer)

	numStrings, _ := strconv.Atoi(sBuffer[1:2])

	fmt.Println(numStrings)

	arrayArray := make([]string, numStrings)

	sliceArray := sBuffer[numStrings:]

	fmt.Println(sliceArray)

	for i := 0; i < numStrings; i++ {

		a := 0
		b := 0

		for i := 0; i < len(sliceArray); i++ {
			if string(sliceArray[i]) == "$" {
				a = i + 1
				b = i

				fmt.Println("a ", a)
				fmt.Println("b ", b)
				break
			}
		}

		cutNumber, _ := strconv.Atoi(sliceArray[b+1 : a+1])
		fmt.Println("cutnumber ", cutNumber)

		arrayArray[i] = sliceArray[a+3 : cutNumber+a+3]

		newSliceArray := sliceArray[cutNumber+a+3+b:]

		sliceArray = newSliceArray

		fmt.Println("slice ", sliceArray)
		fmt.Println("arr ", arrayArray[i])

	}

}
