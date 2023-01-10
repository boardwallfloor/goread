package main

import (
	"compress/gzip"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Send HTML
func (app *App) Send(book Book) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")
		tmpl, err := template.ParseFiles("index.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		gz := gzip.NewWriter(w)
		defer gz.Close()

		err = tmpl.Execute(gz, book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})
	port := fmt.Sprintf(":%s", app.port)
	log.Printf("Listening on port %s\n", app.port)
	log.Fatal(http.ListenAndServe(port, nil))
}
