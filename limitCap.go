package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func limitCap(req Request) {

	for {
		mutex.Lock()
		fmt.Printf("Bucket size: %d\n", len(bucket))
		if len(bucket) >= req.size {
			fmt.Println("Request run ", req)
			for i := 0; i < req.size; i++ {
				<-bucket
			}
			fmt.Printf("Request %d finished\n", req.id)
			fmt.Printf("Bucket size: %d (after request run)\n", len(bucket))
			break
		} else {
			fmt.Printf("Request is too large!\n ID: %d  Size: %d\n", req.id, req.size)
			mutex.Unlock()
		}
		time.Sleep(time.Duration(1) * time.Second)
	}

	mutex.Unlock()
}

type Request struct {
	id   int
	size int
}

func run(req Request) {
	time.Sleep(time.Duration(1) * time.Second)
	limitCap(req)

}

func fill(bucket chan int, freq int) {
	for {
		time.Sleep(time.Duration(freq) * time.Second)

		//Se o mutex estiver fechado, pode ser que o tempo de adicionar token
		//seja maior do que o time.Sleep e aí não daria certo
		// mutex.Lock()
		if len(bucket) < cap(bucket) {
			bucket <- 1
			fmt.Printf("Token adicionado!\n")
		}
		// mutex.Unlock()

	}
}

var bucket chan int
var mutex sync.Mutex

func main() {

	B := 10

	R := 1

	bucket = make(chan int, B)
	//Filling the bucket at the beggining
	for i := 0; i < B; i++ {
		bucket <- 1
	}

	go fill(bucket, 1/R)

	// Defining the size of a Request: how many tokens it will consume from the bucket
	rand.Seed(time.Now().UnixNano())

	id := 0
	for i := 0; i < 5; i++ {
		go func() {
			for {
				id++
				size := rand.Intn(6)
				req := Request{id, size}
				run(req)
			}
		}()
	}

	joinCh := make(chan int)

	<-joinCh
}
