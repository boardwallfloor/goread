package main

import (
	"archive/zip"
	"fmt"
	"path/filepath"
	"regexp"
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
