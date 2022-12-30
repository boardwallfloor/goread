package main

import (
	"log"
)

func main() {

	zipReader := ReadEpub("test/Cooking with Wild Game 1 - Eda.epub")
	err := CheckMime(zipReader.File[0])
	if err != nil {
		log.Fatal(err)
	}

	// mappedZipFile, _, err := MapContent(zipReader)
	mappedZipFile, opf, err := MapContent(zipReader)
	if err != nil {
		log.Fatal(err)
	}

	structure, err := CreateStructure(opf)
	if err != nil {
		log.Fatal(err)
	}

	pageList := EnsurePageList(structure, mappedZipFile)

	bodyNode, err := GenerateNode(pageList)
	if err != nil {
		log.Fatal(err)
	}
	EncodeImage(mappedZipFile, bodyNode)
	book, err := RenderBody(bodyNode)
	if err != nil {
		log.Fatal(err)
	}

	Send(book)
}
