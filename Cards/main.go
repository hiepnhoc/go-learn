package main

import "fmt"

func main() {
	// cards := newDeck()
	// // cards.print()
	// hand, remainingCards := deal(cards, 5)

	// hand.print()
	// remainingCards.print()

	// greeting := "Hi there!"

	// fmt.Println([]byte(greeting))

	// fmt.Println(cards.toString())
	// cards.saveToFile("test.txt")
	// fmt.Println(newDeckFromFile("test.txt"))
	// cards.shuffle()
	// cards.print()

	numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for _, number := range numbers {
		if number%2 == 0 {
			fmt.Println(number, " is a even")
		}

		if number%2 != 0 {
			fmt.Println(number, " is a odd")
		}
	}
}
