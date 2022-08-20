package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func itemsStream(itemsCh chan string) {

	// i := 0
	for i := 0; i < 15; i++ {
		// for {
		itemsCh <- "item " + strconv.Itoa(i)
	}
	close(itemsCh)

}

type Bid struct {
	item      string
	bidValue  int
	bidFailed bool
}

func handle(nServers int, itemsCh <-chan string) chan Bid {

	bidCh := make(chan Bid)
	joinCh := make(chan int)

	for i := 0; i < nServers; i++ {
		go func() {
			for item := range itemsCh {
				bidCh <- bid(item)
			}
			joinCh <- 1

		}()
	}

	go func() {
		for i := 0; i < nServers; i++ {
			<-joinCh
		}
		close(joinCh)
		close(bidCh)
	}()

	return bidCh
}

func bid(item string) Bid {
	time.Sleep(time.Second * 5)
	rand.Seed(42)
	return Bid{item, rand.Intn(10), false}
}

func main() {

	itemsCh := make(chan string)

	go itemsStream(itemsCh)

	bidCh := handle(5, itemsCh)

	for bid := range bidCh {
		fmt.Println(bid)
	}

	fmt.Printf("Leilão encerrado\n")
}
