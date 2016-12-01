package model

import (
	"github.com/gorilla/websocket"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	// codes defined by RFC 6455
	CloseNormalClosure           = 1000
	CloseGoingAway               = 1001
	CloseProtocolError           = 1002
	CloseUnsupportedData         = 1003
	CloseNoStatusReceived        = 1005
	CloseAbnormalClosure         = 1006
	CloseInvalidFramePayloadData = 1007
	ClosePolicyViolation         = 1008
	CloseMessageTooBig           = 1009
	CloseMandatoryExtension      = 1010
	CloseInternalServerErr       = 1011
	CloseServiceRestart          = 1012
	CloseTryAgainLater           = 1013
	CloseTLSHandshake            = 1015
)

const (
	N             = 2 // number of data per message
	TextMessage   = 1
	BinaryMessage = 2
	CloseMessage  = 8
	PingMessage   = 9
	PongMessage   = 10
)

var connections map[*websocket.Conn]bool
var upgrader websocket.Upgrader // RFC 7692
var timer uint16

// var templates = template.Must(template.ParseGlob("web/view/*"))

type Page struct {
	Title string
	Body  []byte
}

func Init() {
	// init connections map variable (analog to a constructor to our pkg)
	connections = make(map[*websocket.Conn]bool)
	upgrader = websocket.Upgrader{ // Creating upgrader "object"
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// allows connection from anywhere
			return true
		},
	}
	timer = 0
}

func sendAll(msg [N][]byte) error {
	// for dealing with different connections at the same time
	var i uint8
	for connection := range connections {
		for i = 0; i < N; i++ {
			err := connection.WriteMessage(
				TextMessage, msg[i])
			if err != nil {
				delete(connections, connection)
				return err
			}
		}
	}
	return nil
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

func sendPing() error {
	// basic ping function for keeping the socket alive
	// RFC6544
	for connection := range connections {
		for i := 0; i < N; i++ {
			err := connection.WriteMessage(PingMessage, []byte{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	// starting connection between the sockets
	// connection, err := upgrader.Upgrade(w, r, nil, 1024, 1024)
	connection, err := upgrader.Upgrade(w, r, nil)
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
		if timer > 1000 {
			if err = sendPing(); err != nil {
				return
			}
			timer = 0
		} else {
			timer++
		}
		log.Println(string(msg))
		if err = sendAll(parseMsg(msg)); err != nil {
			return
		}
	}
}

func parseMsg(msg []byte) [N][]byte {
	function := func(c rune) bool {
		return c == ','
	}
	fields := strings.FieldsFunc(string(msg), function)
	var result [N][]byte
	for i := 0; i < N; i++ {
		result[i] = []byte(fields[i])
	}
	return result
}
