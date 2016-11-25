package model

import (
	"github.com/gorilla/websocket"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const N uint8 = 2 // number of data per message

var connections map[*websocket.Conn]bool

// var templates = template.Must(template.ParseGlob("web/view/*"))

type Page struct {
	Title string
	Body  []byte
}

func Init() {
	// init connections map variable (analog to a constructor to our pkg)
	connections = make(map[*websocket.Conn]bool)
}

func sendAll(msg [N][]byte) {
	// for dealing with different connections at the same time
	var i uint8
	for connection := range connections {
		for i = 0; i < N; i++ {
			err := connection.WriteMessage(
				websocket.TextMessage, msg[i])
			if err != nil {
				delete(connections, connection)
				return
			}
		}
	}
}

func loadPage(filepath string, title string) (*Page, error) {
	body, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func TemplateHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/room/"):]
	p, err := loadPage("/room.html", title)
	if err != nil {
		p = &Page{Title: "ChatINSPER"}
	}
	t, err := template.New("room.html").
		Delims("<<", ">>").ParseFiles("room.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	t.Execute(w, p)
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
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

func parseMsg(msg []byte) [N][]byte {
	var i uint8
	function := func(c rune) bool {
		return c == ','
	}
	fields := strings.FieldsFunc(string(msg), function)
	var result [N][]byte
	for i = 0; i < N; i++ {
		result[i] = []byte(fields[i])
	}
	return result
}
