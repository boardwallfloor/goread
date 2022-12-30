package main

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

type ItemRef struct {
	XMLName xml.Name
	IdRef   string `xml:"idref,attr"`
}

type Item struct {
	XMLName xml.Name
	Id      string `xml:"id,attr"`
	Href    string `xml:"href,attr"`
	Type    string `xml:"media-type,attr"`
}

type Metadata struct {
	XMLName xml.Name
	Title   string `xml:"title"`
	Creator string `xml:"creator"`
}

type Manifest struct {
	XMLName xml.Name
	Items   []Item `xml:"item"`
}

type Spine struct {
	XMLName  xml.Name
	ItemRefs []ItemRef `xml:"itemref"`
}

type Package struct {
	XMLName  xml.Name `xml:"package"`
	Metadata Metadata `xml:"metadata"`
	Manifest Manifest `xml:"manifest"`
	Spine    Spine    `xml:"spine"`
}

type Page struct {
	Title string
	Body  template.HTML
}

func ReadEpub(path string) *zip.ReadCloser {
	log.Println("Reading epub")
	r, err := zip.OpenReader(path)
	if err != nil {
		log.Fatal(err)
	}
	return r
}

// Check for mimetype file
func CheckMime(zr *zip.File) error {
	log.Println("Finding mimetype file")
	if zr.Name != "mimetype" {
		log.Fatal("mimetype file not found")
	}
	log.Println("mimetype found")
	mimeRc, err := zr.Open()
	if err != nil {
		return err
	}
	defer mimeRc.Close()

	st := new(strings.Builder)
	_, err = io.Copy(st, mimeRc)
	if err != nil {
		return err
	}
	mimetype := st.String()
	if mimetype != "application/epub+zip" {
		return errors.New("invalid mimetype")
	}
	log.Printf("mimetype : %s\n", mimetype)
	log.Println("mimetype confirmed")
	return nil
}

// Search for content.opf
// &
// Map zip file with file name
func MapContent(zr *zip.ReadCloser) (map[string]*zip.File, *zip.File, error) {
	log.Println("Searching content.opf")
	epubMap := make(map[string]*zip.File, 0)
	var opf *zip.File
	for _, v := range zr.File {
		epubMap[filepath.Base(v.Name)] = v
		if filepath.Base(v.Name) == "content.opf" {
			opf = v
			log.Println("content.opf found")
		}
	}

	if filepath.Base(opf.Name) != "content.opf" {
		return nil, nil, errors.New("content.opf not found")
	}
	return epubMap, opf, nil
}

// Parse content.opf
func CreateStructure(opf *zip.File) (Package, error) {
	log.Println("Parsing content.opf")
	log.Println("Opening xml zip")
	xmlFile, err := opf.Open()
	if err != nil {
		return Package{}, err
	}
	log.Println("XML file open")

	log.Println("Reading XML")
	xmlstream, err := io.ReadAll(xmlFile)
	if err != nil {
		return Package{}, err
	}

	log.Println("Converting XML")
	var content Package
	err = xml.Unmarshal(xmlstream, &content)
	if err != nil {
		return Package{}, err
	}

	err = xmlFile.Close()
	if err != nil {
		return Package{}, err
	}
	log.Println("Parsing content.opf succeed")
	return content, nil
}

func EnsurePageList(structure Package, mappedZipFile map[string]*zip.File) []*zip.File {
	pageList := []*zip.File{}
	for _, v := range structure.Spine.ItemRefs {
		page := mappedZipFile[v.IdRef]

		if page == nil {
			for _, item := range structure.Manifest.Items {
				if v.IdRef == item.Id {
					page = mappedZipFile[filepath.Base(item.Href)]
					pageList = append(pageList, page)
				}
			}
		} else {
			pageList = append(pageList, page)
		}
	}
	return pageList
}

// Form html with asset encoded
func ProcessBody(page *zip.File) (string, error) {
	//* xhtml give are based on the structure on spine
	//* Find body
	//* If img tag exist replace src with encoded image

	pageRc, err := page.Open()
	if err != nil {
		log.Fatal(err)
	}

	// Parse the HTML byte slice into a node tree
	doc, err := html.Parse(pageRc)
	if err != nil {
		return "", err
	}

	// Find the <body> tag
	var body *html.Node
	var findBody func(*html.Node)
	findBody = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			body = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findBody(c)
		}
	}
	findBody(doc)

	// Render the contents of the <body> tag as a string
	var bodyTmpl bytes.Buffer
	err = html.Render(&bodyTmpl, body)
	if err != nil {
		return "", err
	}
	var strBody strings.Builder
	_, err = strBody.Write(bodyTmpl.Bytes())
	if err != nil {
		return "nil", err
	}
	return strBody.String(), nil
}

func ProcessPage(epubMap map[string]*zip.File) {
	// xhtml give are based on the structure on spine
	// Find body
	// If img tag exist replace src with encoded image
	firstTag := getFirstTag(epubMap, ".xhtml")
	testRead, err := firstTag.Open()
	if err != nil {
		log.Fatal(err)
	}
	// var page bytes.Buffer
	// _, err = io.Copy(&page, testRead)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Parse the HTML byte slice into a node tree
	doc, err := html.Parse(testRead)
	if err != nil {
		log.Fatal(err)
	}

	// Find the <body> tag
	var body *html.Node
	var findBody func(*html.Node)
	findBody = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			body = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findBody(c)
		}
	}
	findBody(doc)

	// Render the contents of the <body> tag as a string
	var bodyTmpl bytes.Buffer
	err = html.Render(&bodyTmpl, body)
	if err != nil {
		log.Fatal(err)
	}

	//* Send it
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		para := Page{
			Title: "My Page",
			Body:  template.HTML(bodyTmpl.String()),
		}
		tmpl, err := template.ParseFiles("index.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, para)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	log.Println("Listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
