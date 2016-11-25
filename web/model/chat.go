package model

import (
	"flag"
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

func Init() *http.ServeMux {
	// init connections map variable (analog to a constructor to our pkg)
	connections = make(map[*websocket.Conn]bool)

	// define routers mappings
	views := flag.String("web directory", "web/view/", "website")
	css := flag.String("css directory", "web/css/", "css")
	controller := flag.String("controller directory", "web/controller",
		"controller")
	flag.Parse()

	r := http.NewServeMux()

	fileServer := http.Dir(*views)
	fileHandler := http.FileServer(fileServer)
	cssHandler := http.FileServer(http.Dir(*css))
	controllerHandler := http.FileServer(http.Dir(*controller))
	r.HandleFunc("/chat", WsHandler)
	r.Handle("/", fileHandler)
	r.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	r.Handle("/controller/",
		http.StripPrefix("/controller/", controllerHandler))

	return r
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

func loadPage(filepath string) (*Page, error) {
	body, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return &Page{Title: filepath, Body: body}, nil
}

func TemplateHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadPage("web/view/chat.html")
	if err != nil {
		p = &Page{Title: "ChatINSPER"}
	}
	t := template.New(string(p.Body)).Delims("<<", ">>")
	err = t.Execute(w, p)
	if err != nil {
		log.Println("Executing template:", err)
	}
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/web/view/chat.html", http.StatusFound)
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
