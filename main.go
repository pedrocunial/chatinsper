package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func chatHandler(w http.ResponseWriter, r *http.Request) {
	// starting connection between the sockets
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
	} else if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close() // close connection object
	for {              // inf loop (the communication itself)
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return // break the inf loop
		}
		log.Println(string(msg))
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return // break the inf loop
		}
	}
}

func main() {
	// Server flags
	port := flag.Int("port", 8888, "server port")
	dir := flag.String("directory", "web/", "web files directory")
	flag.Parse()

	// Handler for requests
	fileServer := http.Dir(*dir)
	fileHandler := http.FileServer(fileServer)
	http.Handle("/", fileHandler)
	http.HandleFunc("/chat", chatHandler)

	log.Printf("Running on port %d", *port)

	addr := fmt.Sprintf("localhost:%d", *port)
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}
