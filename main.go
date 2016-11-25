package main

import (
	"chatinsper/web/model"
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

	r := model.Init()

	fmt.Printf("On port %s\n", port)

	addr := ":" + port
	err := http.ListenAndServe(addr, handlers.CompressHandler(r))
	fmt.Println("Listening to: %s", addr)
	fmt.Println(err.Error())
}
