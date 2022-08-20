package main

import (
	"fmt"
	"math/rand"
	"time"
)

func sleepRoutines(joinChan chan int, my_id int) {
	v := rand.Intn(5)
	fmt.Printf("ID: [%d] Vou dormir por %d segundos\n", my_id, v)
	time.Sleep(time.Second * time.Duration(v))
	fmt.Printf("ID: [%d] Acordei :)\n", my_id)
	joinChan <- 1
}

func main() {

	n := 5
	joinCh := make(chan int)

	for i := 0; i < n; i++ {
		go sleepRoutines(joinCh, i)
	}

	for i := 0; i < n; i++ {
		<-joinCh
	}

	close(joinCh)
	fmt.Printf("Todas as GoRoutines terminaram\n")
}
