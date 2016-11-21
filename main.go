package main

import (
	"chatinsper/web/model"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Server flags
	// port := flag.Int("port", 8888, "server port")
	port := os.Getenv("PORT")

	if port == "" {
		log.Println("$PORT not found, using 8888")
		port = "8888"
	}

	dir := flag.String("directory", "web/", "website views")
	flag.Parse()

	model.Init()

	// Handler for requests
	fileServer := http.Dir(*dir)
	fileHandler := http.FileServer(fileServer)
	http.Handle("/", fileHandler)
	http.HandleFunc("/chat", model.ChatHandler)

	fmt.Printf("On port %s\n", port)

	addr := ":" + port
	err := http.ListenAndServe(addr, nil)
	fmt.Println("Listening to: %s", addr)
	fmt.Println(err.Error())
}
