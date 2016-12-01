package main

import (
	"chatinsper/web/controller"
	"fmt"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"
)

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://chatinsper.herokuapp.com"+
		r.RequestURI, http.StatusMovedPermanently)
}

func main() {
	// Server flags
	// port := flag.Int("port", 8888, "server port")
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT not found, using 8888")
		port = "8888"
	}
	port = "443"

	r := controller.Init()

	fmt.Printf("On port %s\n", port)

	addr := ":" + port
	err := http.ListenAndServe(addr, handlers.CompressHandler(r))
	fmt.Println("Listening to: %s", addr)
	fmt.Println(err.Error())
}
