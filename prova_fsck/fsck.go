package main

import (
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var fsckFiles int
var fsckDirs int
var dmgdFiles int
var dmgdDirs int

var mutexDmgFiles sync.Mutex
var mutexDmgDirs sync.Mutex

var mutexFsckFiles sync.Mutex
var mutexFsckDirs sync.Mutex
var mutexParents sync.Mutex
var mutexCorrupted sync.Mutex
var mutexVisitedDirs sync.Mutex

func readDir(path string, filesOutChannel chan<- string, parents map[string]string, visitedDirs map[string]bool) {
	err := filepath.WalkDir(path, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
			return err
		}
		if !dirEntry.IsDir() {
			filesOutChannel <- path
		} else {
			mutexVisitedDirs.Lock()
			visitedDirs[path] = false
			mutexVisitedDirs.Unlock()
		}
		mutexParents.Lock()
		parents[path] = filepath.Dir(path)

		mutexParents.Unlock()

		return nil
	})

	if err != nil {
		panic(err)
	}

	close(filesOutChannel)

}

func readFirstByte(filesInChannel <-chan string, pathCorrupted map[string]bool, joinCh chan int) {

	for filePath := range filesInChannel {
		time.Sleep(2 * time.Second)
		file, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("ERROR: ", err)
		}
		if len(file) > 0 {
			go fsckFile(filePath, file, pathCorrupted)

		}
	}
	joinCh <- 1

}

func fsckDir(path string) bool {
	rn := rand.Intn(2)
	dirCorrupted := false
	if rn%2 == 0 {
		mutexDmgDirs.Lock()
		dmgdDirs++
		dirCorrupted = true
		mutexDmgDirs.Unlock()
	}
	mutexFsckDirs.Lock()
	fsckDirs++
	mutexFsckDirs.Unlock()
	return dirCorrupted
}

func fsckFile(filePath string, file []byte, pathCorrupted map[string]bool) {
	mutexFsckFiles.Lock()
	fsckFiles++
	mutexFsckFiles.Unlock()
	if file[0]%2 == 0 {
		mutexCorrupted.Lock()
		pathCorrupted[filePath] = true
		mutexCorrupted.Unlock()
		mutexDmgFiles.Lock()
		dmgdFiles++
		fmt.Printf("File: %s is corrupted! | First Byte: %b\n", filePath, file[0])
		mutexDmgFiles.Unlock()

	} else {
		mutexCorrupted.Lock()
		pathCorrupted[filePath] = false
		mutexCorrupted.Unlock()
	}

}

func checkedFsckDirs(parents map[string]string, pathCorrupted map[string]bool, visitedDirs map[string]bool, joinCh chan int) {

	for fileCorrupted := range pathCorrupted {
		p := parents[fileCorrupted]
		if !visitedDirs[p] {
			checkCorrupted(p, parents, visitedDirs)

		}
	}

	joinCh <- 1

}

func checkCorrupted(path string, parents map[string]string, visitedDirs map[string]bool) {
	if !visitedDirs[path] {
		mutexVisitedDirs.Lock()
		visitedDirs[path] = true
		mutexVisitedDirs.Unlock()
		if fsckDir(path) {
			checkCorrupted(parents[path], parents, visitedDirs)
		}
	} else {
		return
	}

}

func report() {
	for {
		time.Sleep(1 * time.Second)
		fmt.Printf("\nfscked_files %d damaged_files %d fscked_dirs %d damaged_dirs %d\n", fsckFiles, dmgdFiles, fsckDirs, dmgdDirs)
	}
}

func main() {

	args := os.Args
	if len(args) == 1 {
		fmt.Println("É necessário passar o caminho do diretório root a ser passado")
		fmt.Println("Ex.:: go run fsck teste_dir")
		panic("Faltou o caminho do diretório root")
	}
	filepath := args[1]
	// Sequencial

	filesChannel := make(chan string, 200)
	joinCh := make(chan int)

	dmgdFiles = 0
	dmgdDirs = 0

	fsckFiles = 0
	fsckDirs = 0

	// Vai guardar os pais de cada arquivo ou diretório
	parents := make(map[string]string)

	visitedDirs := make(map[string]bool)

	// Vai dizer de um dado caminho(diretorio ou arquivo) tá corrompido
	pathCorrupted := make(map[string]bool)

	go readDir(filepath, filesChannel, parents, visitedDirs)
	go readFirstByte(filesChannel, pathCorrupted, joinCh)

	go report()

	// go func() {
	// 	for {
	// 		time.Sleep(1 * time.Second)
	// 		fmt.Print("Parents map: \n", parents)
	// 		fmt.Print("\nPaths corrupted: \n", pathCorrupted)
	// 	}

	// }()

	<-joinCh
	checkedFsckDirs(parents, pathCorrupted, visitedDirs, joinCh)
	<-joinCh
	// time.Sleep(4 * time.Second)

}
