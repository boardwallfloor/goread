package main

import (
	"archive/zip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"

	"golang.org/x/net/html"
)

func PrintKey(epubMap map[string]*zip.File) {
	for i := range epubMap {
		fmt.Println(i)
	}
}

func PrintFileName(epubMap map[string]*zip.File) {
	for _, v := range epubMap {
		fmt.Println(filepath.Base(v.Name))
	}
}

func GetFirstTag(epubMap map[string]*zip.File, tag string) *zip.File {
	for _, v := range epubMap {
		rgx := regexp.MustCompile(tag)
		if rgx.MatchString(v.Name) {
			return v
		}

	}
	return nil
}

func EncodeImage(mappedZipFile map[string]*zip.File, bodyNode []*html.Node) {
	var traverseBody func(*html.Node)
	traverseBody = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			for i, v := range n.Attr {
				if v.Key == "src" {
					src := n.Attr[i]
					imageStream := mappedZipFile[filepath.Base(src.Val)]
					rc, err := imageStream.Open()
					if err != nil {
						log.Fatal(err)
					}
					data, err := ioutil.ReadAll(rc)
					if err != nil {
						log.Fatal(err)
					}

					// Encode the image data as a base64 string
					encoded := base64.StdEncoding.EncodeToString(data)
					rc.Close()
					src.Val = fmt.Sprintf("data:image/jpeg;base64,%s", encoded)
					n.Attr[i].Val = src.Val
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverseBody(c)
		}
	}
	for _, v := range bodyNode {
		traverseBody(v)
	}
}
