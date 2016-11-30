package main

import (
	"chatinsper/web/controller"
	"fmt"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"
)

func openconnection(port string, r *http.ServeMux) {
	fmt.Println(port)
	addr := ":" + port
	err := http.ListenAndServe(addr, handlers.CompressHandler(r))
	fmt.Println("Listening to: %s", addr)
	fmt.Println(err.Error())
}

func main() {
	// Server flags
	r1, r2 := controller.Init()

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT not found, using 8888")
		port = "8888"
	}
	fmt.Printf("On port %s\n", port)
	// threads
	go openconnection(port, r1)
	go openconnection("5050", r2)
}
