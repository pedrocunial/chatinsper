package model

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
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
