package main

import (
	"archive/zip"
	"errors"
	"io"
	"log"
	"regexp"
	"strings"
)

func HandlerEpub() []io.ReadCloser {
	r, err := zip.OpenReader("test.epub")
	if err != nil {
		log.Fatal(err)
	}
	// defer r.Close()
	// checkMime(r.File)

	// Parse(r.File, "container.xml")
	firstPage := regexp.MustCompile(`.xhtml`)
	var protopages []zip.File

	for _, v := range r.File {
		if firstPage.MatchString(v.Name) {
			protopages = append(protopages, *v)
		}
		if len(protopages) >= 18 {
			break
		}
	}

	var pages []io.ReadCloser
	for _, v := range protopages {
		page, err := v.Open()
		if err != nil {
			log.Fatal(err)
		}
		pages = append(pages, page)

	}
	if err != nil {
		log.Fatal(err)
	}
	// _,err := io.Copy()

	return pages

}

func checkMime(zf []*zip.File) (string, error) {
	log.Println("Finding mimetype")
	fileHeader := zf[0]
	switch fileHeader.Name {
	case "mimetype":
		log.Println("mimetype file found")
		fileRc, err := fileHeader.Open()
		if err != nil {
			return "", err
		}
		st := new(strings.Builder)
		_, err = io.Copy(st, fileRc)
		if err != nil {
			return "", err
		}
		mimetype := st.String()
		if mimetype != "application/epub+zip" {
			return "", errors.New("invalid mimetype")
		}
		log.Printf("mimetype : %s\n", mimetype)
		log.Println("mimetype confirmed")
		fileRc.Close()
	default:
		return "", errors.New("unexpected error occurs")
	}
	return "", errors.New("unexpected exception")
}

func main() {
	// handlerEpub()
	ServeBook()
}
