// Essa solução não tá completa pois em algum momento seria necessário fechar o canal
// Mas, por incrível que pareça, ela funciona

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

func search_files(path string, files_ch chan string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			search_files(path+"/"+file.Name(), files_ch)
		} else {
			fmt.Printf("Parent directory: %s\n", filepath.Dir(path+"/"+file.Name()))
			fmt.Printf("File %s is damaged\n", path+"/"+file.Name())
			files_ch <- path + "/" + file.Name()
		}
	}
}

func read_files(files_ch chan string) {
	for filename := range files_ch {

		content, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		if content[0]%2 == 0 {
			fmt.Printf("First byte: %d", content[0])
			fmt.Printf("\nFiilename: %s \n", filename)
		}
	}

}

func main() {

	filepath := "./ex3/test_dir"
	files_ch := make(chan string)
	go read_files(files_ch)
	search_files(filepath, files_ch)

	fmt.Printf("Terminou")
}
