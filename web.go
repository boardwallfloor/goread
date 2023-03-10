package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

func HandleHtml(book Book) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// var zipBuff bytes.Buffer
		// gz := gzip.NewWriter(&zipBuff)
		gz := gzip.NewWriter(w)
		defer gz.Close()

		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var buf bytes.Buffer

		err = tmpl.Execute(&buf, book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		m := minify.New()
		m.AddFunc("text/html", html.Minify)
		minifiedHTML, err := m.String("text/html", buf.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		gz.Write([]byte(minifiedHTML))
		// log.Println(len(buf.String()), len(minifiedHTML), len(zipBuff.Bytes()))
	}

}

// Send HTML
func (app *App) Send(book Book) {
	log.Println("Sending html")
	http.HandleFunc("/", HandleHtml(book))
	port := fmt.Sprintf(":%s", app.port)
	log.Printf("Listening on port %s\n", app.port)
	log.Fatal(http.ListenAndServe(port, nil))
}
