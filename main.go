package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"

	client "github.com/khubo/go-torrent/client"
	decoder "github.com/khubo/go-torrent/decoder"
)

func main() {

	// get the torrent file from user passed args
	file := os.Args[1]
	
	// generate a random client id for our peer
	var clientId [20]byte
	port := 19119
	_, err := rand.Read(clientId[:])
	if err != nil {
		log.Fatalln(err)
	}

	//decode the torrent file
	torrentFile, err := decoder.Decode(file)
	if err != nil {
		log.Fatalln(err)
	}

	peers, err := client.RequestPeers(torrentFile, clientId, port)
	
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%T", peers)
}