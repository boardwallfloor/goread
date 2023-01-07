package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
)

func main() {
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	zipReader := ReadEpub("test/Cooking with Wild Game 1 - Eda.epub")
	err := CheckMime(zipReader.File[0])
	if err != nil {
		log.Fatal(err)
	}

	// mappedZipFile, _, err := MapContent(zipReader)
	// _, opf, err := MapContent(zipReader)
	mappedZipFile, opf, err := MapContent(zipReader)
	if err != nil {
		log.Fatal(err)
	}

	structure, err := CreateStructure(opf)
	if err != nil {
		log.Fatal(err)
	}

	// EnsureNav(structure.Guide.References, structure.Manifest.Items, mappedZipFile)
	pageList := EnsurePageList(structure, mappedZipFile)

	bodyNode, err := GenerateNode(pageList, mappedZipFile)
	if err != nil {
		log.Fatal(err)
	}

	book, err := RenderBody(bodyNode)
	if err != nil {
		log.Fatal(err)
	}

	Send(book)
	zipReader.Close()
}
