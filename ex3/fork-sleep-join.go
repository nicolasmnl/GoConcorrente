package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func sleepRoutines(joinChan chan int, my_id int) {
	v := rand.Intn(10)
	fmt.Printf("ID: [%d] Vou dormir por %d segundos\n", my_id, v)
	time.Sleep(time.Second * time.Duration(v))
	fmt.Printf("ID: [%d] Acordei :)\n", my_id)
	joinChan <- 1
}

func main() {

	// n := 5

	args := os.Args
	if len(args) == 1 {
		fmt.Println("É necessário passar a quantidade de GoRoutines a serem criadas")
		fmt.Println("Ex.:: go run fork-sleep-join 10")
		panic("Faltou o número de GoRoutines")
	}
	n, err := strconv.Atoi(args[1])
	if err != nil {
		panic(err)
	}

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
