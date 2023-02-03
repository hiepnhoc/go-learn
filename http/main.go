package main

import (
	"fmt"
)

type logWriter struct {
}

type triangle struct {
	height float64
	base   float64
}

type square struct {
	sideLength float64
}

type shape interface {
	printArea() float64
}

func main() {
	// resp, err := http.Get("http://google.com")

	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	os.Exit(1)
	// }

	// lw := logWriter{}

	// io.Copy(lw, resp.Body)

	tri := triangle{base: 10, height: 10}
	squa := square{sideLength: 10}

	printArea(tri)
	printArea(squa)
}

func (logWriter) Write(bs []byte) (int, error) {
	fmt.Println(string(bs))
	fmt.Println("Len bytes:", len(bs))
	return len(bs), nil
}

func printArea(s shape) {
	fmt.Println(s.printArea())
}

func (t triangle) printArea() float64 {
	return 0.5 * t.base * t.height
}

func (s square) printArea() float64 {
	return s.sideLength * s.sideLength
}
