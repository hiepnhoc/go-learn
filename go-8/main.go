package main

import (
	"fmt"
	"math/rand"
	"time"
)

func startSender(name string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 1; i <= 5; i++ {
			c <- (name + " hello ")
			time.Sleep(time.Second)
		}
	}()
	return c
}

func fetchAPI(model string) string {
	time.Sleep(time.Duration(time.Duration(rand.Intn(1e3)) * time.Millisecond))
	return model
}

func query(model string) string {
	time.Sleep(time.Duration(time.Duration(rand.Intn(2e3)) * time.Millisecond))
	return model
}

func queryFirst(servers ...string) <-chan string {
	c := make(chan string)

	for _, serv := range servers {
		go func(s string) {
			c <- query(s)
		}(serv)
	}

	return c
}

func main() {
	// sender := startSender("Ti")
	// teo := startSender("Teo")
	// // for i := 1; i <= 5; i++ {
	// // 	fmt.Println(<-sender)
	// // }

	// for {
	// 	select {
	// 	case msgTi := <-sender:
	// 		fmt.Println(msgTi)
	// 	case msgTeo := <-teo:
	// 		fmt.Println(msgTeo)
	// 	}
	// }

	// response := make(chan string)
	// var result []string

	// go func() {
	// 	response <- fetchAPI("User")
	// }()

	// go func() {
	// 	response <- fetchAPI("categori")
	// }()

	// go func() {
	// 	response <- fetchAPI("product")
	// }()

	// for i := 1; i <= 3; i++ {
	// 	result = append(result, <-response)
	// }

	// fmt.Println(result)

	result := queryFirst("server 1", " server 2", "server 3")

	fmt.Println((<-result))

}
