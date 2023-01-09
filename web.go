package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Send HTML
func (app *App) Send(book Book) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		tmpl, err := template.ParseFiles("index.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	port := fmt.Sprintf(":%s", app.port)
	log.Printf("Listening on port %s\n", app.port)
	log.Fatal(http.ListenAndServe(port, nil))
}
