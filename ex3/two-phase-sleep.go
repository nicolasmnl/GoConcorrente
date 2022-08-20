package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func two_phase_sleep(myTimeToSleep chan int, timesToSecondSleep chan int, join_ch chan int, wg *sync.WaitGroup, my_id int, second_id int) {

	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	s := rand.Intn(10)
	timesToSecondSleep <- s

	// Espero todo mundo colocar nos seus respectivos canais
	// barrier
	wg.Done() // avisa que terminou de fazer a primeira fase
	fmt.Printf("Eu sou a GoRoutine %d e a próxima GoRoutine %d vai dormir por %d segundos\n", my_id, second_id, s)
	wg.Wait() // Espero todo mundo fazer a primeira fase

	v := <-myTimeToSleep
	time.Sleep(time.Duration(v) * time.Second)
	fmt.Printf("Eu sou a GoRoutine %d e eu dormi por %d segundos\n", my_id, v)
	join_ch <- 1

}

func main() {

	n := 5

	var channels [5]chan int

	var barrier sync.WaitGroup
	barrier.Add(n)

	join_ch := make(chan int)

	for i := 0; i < n; i++ {
		channels[i] = make(chan int, 1)
	}

	for i := 0; i < n; i++ {
		go two_phase_sleep(channels[i], channels[(i+1)%n], join_ch, &barrier, i, i+1)
	}

	for i := 0; i < n; i++ {
		<-join_ch
	}

	fmt.Printf("All goroutines finished\n")
}
