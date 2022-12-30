package main

import (
	"archive/zip"
	"fmt"
	"path/filepath"
	"regexp"
)

func printKey(epubMap map[string]*zip.File) {
	for i, _ := range epubMap {
		fmt.Println(i)
	}
}

func printFileName(epubMap map[string]*zip.File) {
	for _, v := range epubMap {
		fmt.Println(filepath.Base(v.Name))
	}
}

func getFirstTag(epubMap map[string]*zip.File, tag string) *zip.File {
	for _, v := range epubMap {
		rgx := regexp.MustCompile(tag)
		if rgx.MatchString(v.Name) {
			return v
		}

	}
	return nil
}
