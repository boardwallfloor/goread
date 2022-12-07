package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
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

func Parse(zf []*zip.File, filename string) {

	log.Println("Parsing XML")
	xmlf, err := os.Open("test.xml")
	if err != nil {
		fmt.Println(err)
	}
	log.Println("XML file open")
	defer xmlf.Close()

	xmlstream, err := io.ReadAll(xmlf)
	if err != nil {
		log.Println(err)
	}

	var content Package
	err = xml.Unmarshal(xmlstream, &content)
	if err != nil {
		log.Println(err)
	}

	// Log Content
	// log.Println("Metadata")
	// log.Printf("Title : %s\n", content.Metadata.Title)
	// log.Printf("Author : %s\n", content.Metadata.Creator)
	// log.Println("Item")
	// for _, v := range content.Manifest.Items {
	// 	log.Printf("id : %s, href : %s, media-type: %s\n", v.Id, v.Href, v.Type)
	// }
	// log.Println("Spine")
	// for _, v := range content.Spine.ItemRefs {
	// 	log.Printf("idref : %s\n", v.IdRef)
	// }
}
