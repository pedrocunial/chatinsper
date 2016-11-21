package model

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
)

var connections map[*websocket.Conn]bool

func Init() {
	// init connections map variable (analog to a constructor to our pkg)
	connections = make(map[*websocket.Conn]bool)
}

func sendAll(msg [][]byte) {
	// for dealing with different connections at the same time
	for connection := range connections {
		err1 := connection.WriteMessage(websocket.TextMessage, msg[0])
		err2 := connection.WriteMessage(websocket.TextMessage, msg[1])
		if err1 != nil || err2 != nil {
			delete(connections, connection)
			return
		}
	}
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	// starting connection between the sockets
	connection, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	defer connection.Close() // close connection object
	connections[connection] = true
	for { // inf loop (the communication itself)
		_, msg, err := connection.ReadMessage()
		if err != nil {
			delete(connections, connection)
			return // break the inf loop
		}
		log.Println(string(msg))
		sendAll(parseMsg(msg))
	}
}

func parseMsg(msg []byte) [][]byte {
	function := func(c rune) bool {
		return c == ','
	}
	fields := strings.FieldsFunc(string(msg), function)
	return [][]byte{[]byte(fields[0]), []byte(fields[1])}
}
