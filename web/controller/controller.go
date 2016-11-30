package controller

import (
	"chatinsper/web/model"
	"flag"
	"net/http"
)

func Init() (*http.ServeMux, *http.ServeMux) {
	model.Init()
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
	// r.HandleFunc("/chat", model.WsHandler)
	r.Handle("/room", fileHandler)
	r.Handle("/", fileHandler)
	r.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	r.Handle("/controller/",
		http.StripPrefix("/controller/", controllerHandler))

	r2 := http.NewServeMux()
	r2.HandleFunc("/chat", model.WsHandler)

	return r, r2
}
