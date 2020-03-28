package main

import (
	"fmt"
	"log"
	"os"
	decoder "github.com/khubo/go-torrent/decoder" 
)

func main() {

	file := os.Args[1]
	torrentFile, err := decoder.Decode(file)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(torrentFile)
}