package main

import "fmt"

func main() {
	// var colors map[string]string

	colors := map[string]string{
		"red":   "#ff0000",
		"green": "#4bf745",
		"white": "#ffffff",
	}

	// colors := make(map[int]string)

	// fmt.Println(colors)
	colors["yellow"] = "123123"
	delete(colors, "yellow")
	printMap(colors)
}

func printMap(c map[string]string) {
	for key, value := range c {
		fmt.Println("hex code for ", key, "is ", value)
	}
}
