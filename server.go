package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	port string = ":8080"
)

func serve(target string) {
	// serve dist dir
	fs := http.FileServer(http.Dir(target))
	http.Handle("/", fs)

	fmt.Printf("[INFO] Listening on http://localhost%s ...\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}

}

func serveInBackground(target string) {
	// serve dist dir
	fs := http.FileServer(http.Dir(target))
	http.Handle("/", fs)

	go func() {
		if err := http.ListenAndServe(port, nil); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Printf("[INFO] Listening on http://localhost%s ...\n", port)
}
