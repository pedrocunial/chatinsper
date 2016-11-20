package main

import (
	"chatinsper/web/model"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Server flags
	port := flag.Int("port", 8888, "server port")
	dir := flag.String("directory", "web/", "website views")
	flag.Parse()

	model.Init()

	// Handler for requests
	fileServer := http.Dir(*dir)
	fileHandler := http.FileServer(fileServer)
	http.Handle("/", fileHandler)
	http.HandleFunc("/chat", model.ChatHandler)

	log.Printf("Running on port %d", *port)

	addr := fmt.Sprintf("localhost:%d", *port)
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}
