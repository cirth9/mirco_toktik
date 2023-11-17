package main

import (
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.OpenFile("./dir/test.txt", os.O_RDONLY, 777)
	if err != nil {
		log.Println(err)
	}
	create, err := os.Create("./dir/test3.txt")
	if err != nil {
		log.Println(err)
	}
	written, err := io.Copy(create, file)
	if err != nil {
		log.Println(err)
	}
	log.Println(written)
}
