package main

import (
	"fmt"

	"./server"
	"./server/storage"
)

func main() {
	datastore := storage.New()
	server := server.New(&datastore)

	fmt.Println("Starting server...")
	server.Start()
	for {
	}
}
