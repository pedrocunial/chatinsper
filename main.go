package main

import (
	"chatinsper/web/model"
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
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
	views := flag.String("web directory", "web/view/", "website")
	css := flag.String("css directory", "web/css/", "css")
	controller := flag.String("controller directory", "web/controller",
		"controller")
	flag.Parse()

	model.Init()

	r := http.NewServeMux()

	// Handler for requests
	// r.HandleFunc("/", model.TemplateHandler)
	fileServer := http.Dir(*views)
	fileHandler := http.FileServer(fileServer)
	cssHandler := http.FileServer(http.Dir(*css))
	controllerHandler := http.FileServer(http.Dir(*controller))
	r.HandleFunc("/chat", model.WsHandler)
	r.Handle("/", fileHandler)
	r.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	r.Handle("/controller/",
		http.StripPrefix("/controller/", controllerHandler))

	fmt.Printf("On port %s\n", port)

	addr := ":" + port
	err := http.ListenAndServe(addr, handlers.CompressHandler(r))
	fmt.Println("Listening to: %s", addr)
	fmt.Println(err.Error())
}
