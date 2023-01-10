package main

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func (app *App) ReadEpub(path string) *zip.ReadCloser {
	app.infoLog.Println("Reading epub")
	r, err := zip.OpenReader(path)
	if err != nil {
		log.Fatal(err)
	}
	size, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Reading :%s, with size :%d", path, size.Size())
	return r
}

// Search for content.opf
// &
// Map zip file with file name
func (app *App) MapContent(zr *zip.ReadCloser) (map[string]*zip.File, *zip.File, error) {
	app.infoLog.Println("Searching content.opf")
	epubMap := make(map[string]*zip.File, 0)
	var opf *zip.File

	opfRegex := regexp.MustCompile(`^.*\.opf$`)

	for _, v := range zr.File {
		epubMap[filepath.Base(v.Name)] = v
		if opfRegex.MatchString(filepath.Base(v.Name)) {
			opf = v
			app.infoLog.Println("content.opf found")
		}

	}

	if opf == nil {
		return nil, nil, errors.New("content.opf not found")
	}

	return epubMap, opf, nil
}

// Parse content.opf
func (app *App) CreateStructure(opf *zip.File) (Package, error) {
	app.infoLog.Println("Parsing content.opf")
	app.infoLog.Println("Opening xml zip")
	xmlFile, err := opf.Open()
	if err != nil {
		return Package{}, err
	}
	app.infoLog.Println("XML file open")

	app.infoLog.Println("Reading XML")
	xmlstream, err := io.ReadAll(xmlFile)
	if err != nil {
		return Package{}, err
	}

	app.infoLog.Println("Converting XML")
	var content Package
	err = xml.Unmarshal(xmlstream, &content)
	if err != nil {
		return Package{}, err
	}

	err = xmlFile.Close()
	if err != nil {
		return Package{}, err
	}
	app.infoLog.Println("Parsing content.opf succeed")
	return content, nil
}

// Ensuring a list that contain all valid key to mappedZipFile of the needed file since <spine> are not guaranteed to be directly corelated to reference
func (app *App) EnsurePageList(structure Package, mappedZipFile map[string]*zip.File) []Page {
	app.infoLog.Println("Validating book file list")
	pageList := []Page{}
	// func getnave is ranging of guides first(reference) then items(manifest) and return marker for nav file? nad map of items
	var navFile *zip.File
	for _, v := range structure.Manifest.Items {
		if v.Properties == "nav" {
			navFile = mappedZipFile[filepath.Base(v.Href)]
			// set nav here
		}
	}
	for _, v := range structure.Spine.ItemRefs {
		meta := make(map[string]string, 0)
		page := mappedZipFile[v.IdRef]
		meta["id"] = v.IdRef
		if page == nil {
			for _, item := range structure.Manifest.Items {
				if filepath.Base(item.Id) == v.IdRef {
					page = mappedZipFile[filepath.Base(item.Href)]
					meta["id"] = item.Href
				}
			}
		}
		if page.Name == navFile.Name {
			meta["type"] = "nav"
		}
		pageList = append(pageList, Page{Page: page, Meta: meta})

	}
	return pageList
}

// Generate html.Node within certain time
func (app *App) GenerateNode(pageList []Page, mappedZipFile map[string]*zip.File) ([]*html.Node, error) {
	app.infoLog.Println("Generating html.Node")
	timer := time.After(4 * time.Second)
	nodeList := []*html.Node{}
	for _, v := range pageList {
		select {
		case <-timer:
			app.infoLog.Println("Timer expired")
			return nodeList, nil
		default:
			body, err := app.GetBodyNode(v.Page)
			if err != nil {
				return nil, fmt.Errorf("%s, for file %s", err, v.Page.Name)
			}
			err = app.ProcessBody(mappedZipFile, body, v)
			if err != nil {
				return nil, fmt.Errorf("%s, for file %s", err, v.Page.Name)
			}
			nodeList = append(nodeList, body)
		}
	}
	return nodeList, nil
}

// Render the contents of the <body> tag as a string
func (app *App) RenderBody(nodes []*html.Node, meta Metadata) (Book, error) {
	var bodyTmpl bytes.Buffer
	for _, v := range nodes {
		err := html.Render(&bodyTmpl, v)
		if err != nil {
			return Book{}, err
		}
	}

	var strBody strings.Builder
	_, err := strBody.Write(bodyTmpl.Bytes())
	if err != nil {
		return Book{}, err
	}
	book := Book{
		Title: meta.Title,
		Body:  template.HTML(strBody.String()),
	}
	return book, nil
}
