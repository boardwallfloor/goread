package main

import (
	"bytes"
	"html/template"
	"log"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func main() {

	zipReader := ReadEpub("test/Cooking with Wild Game 1 - Eda.epub")
	err := CheckMime(zipReader.File[0])
	if err != nil {
		log.Fatal(err)
	}

	mappedZipFile, opf, err := MapContent(zipReader)
	if err != nil {
		log.Fatal(err)
	}

	// _, err = CreateStructure(opf)
	structure, err := CreateStructure(opf)
	if err != nil {
		log.Fatal(err)
	}

	pageList := EnsurePageList(structure, mappedZipFile)

	timer := time.After(4 * time.Second)
	temp := []*html.Node{}
	for _, v := range pageList {
		select {
		case <-timer:
			return
		default:
			body, err := ProcessBody(v)
			temp = append(temp, body)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// Render the contents of the <body> tag as a string
	var bodyTmpl bytes.Buffer
	for _, v := range temp {
		err = html.Render(&bodyTmpl, v)
		if err != nil {
			log.Fatal(err)
		}
	}

	var strBody strings.Builder
	_, err = strBody.Write(bodyTmpl.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	book := Book{
		Title: "My Page",
		Body:  template.HTML(strBody.String()),
	}
	Send(book)
}
