package main

import "fmt"

type contactInfo struct {
	email   string
	zipCode int
}

type person struct {
	firstName string
	lastName  string
	contact   contactInfo
}

func main() {
	jim := person{
		firstName: "Jim",
		lastName:  "Party",
		contact: contactInfo{
			email:   "jim@gmail.com",
			zipCode: 94000,
		},
	}

	jimPointer := &jim
	jimPointer.updateName("hiepln")
	jimPointer.print()
}

func (p person) print() {
	fmt.Printf("%+v", p)
}

func (pointToPerson *person) updateName(newFirstName string) {
	(*pointToPerson).firstName = newFirstName
}
