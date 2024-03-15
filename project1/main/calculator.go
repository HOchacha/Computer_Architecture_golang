package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Error : No input file specified")
		return
	}
	fmt.Println("The input file name : ", os.Args[1])

	//read input text file
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error : did you specified the input text file?")
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
