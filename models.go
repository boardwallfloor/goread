package main

import (
	"archive/zip"
	"encoding/xml"
	"html/template"
)

type Page struct {
	Page *zip.File
	Meta map[string]string
}

type ItemRef struct {
	XMLName xml.Name
	IdRef   string `xml:"idref,attr"`
}

type Item struct {
	XMLName    xml.Name
	Id         string `xml:"id,attr"`
	Href       string `xml:"href,attr"`
	Type       string `xml:"media-type,attr"`
	Properties string `xml:"properties,attr"`
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

type Reference struct {
	XMLName xml.Name `xml:"reference"`
	Type    string   `xml:"type,attr"`
	Title   string   `xml:"title,attr"`
	Href    string   `xml:"href,attr"`
}

type Guide struct {
	References []Reference `xml:"reference"`
}

type Package struct {
	XMLName  xml.Name `xml:"package"`
	Metadata Metadata `xml:"metadata"`
	Manifest Manifest `xml:"manifest"`
	Spine    Spine    `xml:"spine"`
	Guide    Guide    `xml:"guide"`
}

type Book struct {
	Title string
	Body  template.HTML
}
