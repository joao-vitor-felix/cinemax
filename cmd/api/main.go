package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	var port int
	flag.IntVar(&port, "port", 8080, "server port")
	flag.Parse()

	//TODO: setup db connection
	//TODO: setup routes

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		// Handler: router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}

	log.Fatal(srv.ListenAndServe())
}
