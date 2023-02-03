package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	argsWithProg := os.Args[1]

	f, error := os.Open(argsWithProg)

	if error != nil {
		fmt.Println("err", error)
		os.Exit(1)
	}

	io.Copy(os.Stdout, f)
}
