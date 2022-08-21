package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

func two_phase_sleep(myTimeToSleep chan int, timesToSecondSleep chan int, join_ch chan int, wg *sync.WaitGroup, my_id int, second_id int) {

	v1 := rand.Intn(5)
	time.Sleep(time.Duration(v1) * time.Second)
	s := rand.Intn(10)
	timesToSecondSleep <- s

	// Espero todo mundo colocar nos seus respectivos canais
	// barrier
	wg.Done() // avisa que terminou de fazer a primeira fase
	fmt.Printf("Eu sou a GoRoutine %d e dormi por %d segundos\n", my_id, v1)
	fmt.Printf("Eu sou a GoRoutine %d e a próxima GoRoutine %d vai dormir por %d segundos\n\n", my_id, second_id, s)
	wg.Wait() // Espero todo mundo fazer a primeira fase

	v := <-myTimeToSleep
	time.Sleep(time.Duration(v) * time.Second)
	fmt.Printf("Eu sou a GoRoutine %d e eu dormi por %d segundos na segunda fase\n", my_id, v)
	join_ch <- 1

}

func main() {

	args := os.Args
	if len(args) == 1 {
		fmt.Println("É necessário passar a quantidade de GoRoutines a serem criadas")
		fmt.Println("Ex.:: go run two_phase_sleep 10")
		panic("Faltou o número de GoRoutines")
	}
	n, err := strconv.Atoi(args[1])
	if err != nil {
		panic(err)
	}

	channels := make([]chan int, n)

	var barrier sync.WaitGroup
	barrier.Add(n)

	join_ch := make(chan int)

	for i := 0; i < n; i++ {
		channels[i] = make(chan int, 1)
	}

	for i := 0; i < n; i++ {
		go two_phase_sleep(channels[i], channels[(i+1)%n], join_ch, &barrier, i, (i+1)%n)
	}

	for i := 0; i < n; i++ {
		<-join_ch
	}

	fmt.Printf("All goroutines finished\n")
}
