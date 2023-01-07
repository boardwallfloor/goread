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
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func ReadEpub(path string) *zip.ReadCloser {
	log.Println("Reading epub")
	r, err := zip.OpenReader(path)
	if err != nil {
		log.Fatal(err)
	}
	return r
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

// Ensuring a list that contain all valid key to mappedZipFile of the needed file since <spine> are not guaranteed to be directly corelated to reference
func EnsurePageList(structure Package, mappedZipFile map[string]*zip.File) []Page {
	log.Println("Validating book file list")
	pageList := []Page{}
	// func getnave is ranging of guides first(reference) then items(manifest) and return marker for nav file? nad map of items

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
			//* change to map call to items from getNav()
		}

		pageList = append(pageList, Page{Page: page, Meta: meta})

	}
	return pageList
}

// Generate html.Node within certain time
func GenerateNode(pageList []Page, mappedZipFile map[string]*zip.File) ([]*html.Node, error) {
	log.Println("Generating html.Node")
	timer := time.After(4 * time.Second)
	nodeList := []*html.Node{}
	for _, v := range pageList {
		select {
		case <-timer:
			log.Println("Timer expired")
			return nodeList, nil
		default:
			body, err := GetBodyNode(v.Page)
			if err != nil {
				return nil, fmt.Errorf("%s, for file %s", err, v.Page.Name)
			}
			err = ProcessBody(mappedZipFile, body, v)
			if err != nil {
				return nil, fmt.Errorf("%s, for file %s", err, v.Page.Name)
			}
			nodeList = append(nodeList, body)
		}
	}
	return nodeList, nil
}

// Render the contents of the <body> tag as a string
func RenderBody(nodes []*html.Node) (Book, error) {
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
		Title: "My Page",
		Body:  template.HTML(strBody.String()),
	}
	return book, nil
}
