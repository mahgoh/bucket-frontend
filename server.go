package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	port string = ":8080"
)

func serve() {
	// serve dist dir
	fs := http.FileServer(http.Dir("./dist"))
	http.Handle("/", fs)

	go func() {
		if err := http.ListenAndServe(port, nil); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Printf("[INFO] Listening on %s ...\n", port)
}
