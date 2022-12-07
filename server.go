package main

import (
	"io"
	"log"
	"net/http"
)

func ServePage(w http.ResponseWriter, r *http.Request) {
	rc := HandlerEpub()
	for _, v := range rc {
		io.Copy(w, v)

	}
}

func ServeBook() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", ServePage)
	log.Println("Listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
