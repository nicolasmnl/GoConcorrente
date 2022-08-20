package main

import (
	"fmt"
	"math/rand"
	"time"
)

func sleepRoutines(joinChan chan int) {
	time.Sleep(time.Second * time.Duration(rand.Intn(5)))
	fmt.Printf("GoRoutine terminou\n")
	joinChan <- 1
}

func main() {

	n := 5
	joinCh := make(chan int)

	for i := 0; i < n; i++ {
		go sleepRoutines(joinCh)
	}

	for i := 0; i < n; i++ {
		<-joinCh
	}

	close(joinCh)
	fmt.Printf("Todas as GoRoutines terminaram\n")
}
