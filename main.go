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
}
