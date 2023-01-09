package main

import (
	"archive/zip"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func contains(slice []string, element string) bool {
	for _, value := range slice {
		if value == element {
			return true
		}
	}
	return false
}

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

func (app *App) findNodeWithAttr(n *html.Node, attr, value string) *html.Node {
	// Check if the current node has the desired attribute and value
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == attr && a.Val == value {
				return n
			}
		}
	}

	// Recursively search the child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := app.findNodeWithAttr(c, attr, value); result != nil {
			return result
		}
	}

	// If the node and its children do not match, return nil
	return nil
}

func (app *App) EnsureNav(bodyNode *html.Node) {

	// node := app.findNodeWithAttr(bodyNode, "epub:type", "toc")
	// if node == nil {
	// 	notFoundError := fmt.Sprintf("attribute epub:type with value toc is not found on nav file `%s`", node.Data)
	// 	log.Fatal(errors.New(notFoundError))
	// }
	// log.Println("Modifying TOC file")
	// divNode := node.Parent
	// listNode := app.findNodeWithAttr(divNode, "epub:type", "list")
	var traverseNode func(*html.Node) error
	traverseNode = func(n *html.Node) error {
		if n.Type == html.ElementNode && n.Data == "a" {
			for i, v := range n.Attr {
				if v.Key == "href" {
					n.Attr[i].Val = fmt.Sprintf("#%s", filepath.Base(n.Attr[i].Val))
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverseNode(c)
		}
		return nil
	}
	traverseNode(bodyNode)

}

// Check for mimetype file
func (app *App) CheckMime(zr *zip.File) error {
	log.Println("Finding mimetype file")
	if zr.Name != "mimetype" {
		log.Fatal("mimetype file not found")
	}
	log.Println("mimetype found")
	mimeRc, err := zr.Open()
	if err != nil {
		return err
	}
	defer mimeRc.Close()

	st := new(strings.Builder)
	_, err = io.Copy(st, mimeRc)
	if err != nil {
		return err
	}
	mimetype := st.String()
	if mimetype != "application/epub+zip" {
		return errors.New("invalid mimetype")
	}
	log.Printf("mimetype : %s\n", mimetype)
	log.Println("mimetype confirmed")
	return nil
}

// Form html with asset encoded
func (app *App) GetBodyNode(page *zip.File) (*html.Node, error) {
	log.Println(page.Name)
	pageRc, err := page.Open()
	if err != nil {
		log.Fatal(err)
	}

	// Parse the HTML byte slice into a node tree
	log.Println("Create node")
	doc, err := html.Parse(pageRc)
	if err != nil {
		return nil, err
	}

	// Find the <body> tag
	var bodyNode *html.Node
	var findBody func(*html.Node)
	findBody = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			bodyNode = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findBody(c)
		}
	}
	log.Println("Finding body")
	findBody(doc)
	log.Println("Body found")
	return bodyNode, nil
}

// Change body element to div and encode any image
func (app *App) ProcessBody(mappedZipFile map[string]*zip.File, bodyNode *html.Node, page Page) error {
	var traverseBody func(*html.Node) error
	traverseBody = func(n *html.Node) error {
		if n.Type == html.ElementNode && n.Data == "body" {
			log.Println("Modifying body to div")
			n.Data = "div"
			idAttr := html.Attribute{Key: "id", Val: page.Meta["id"]}
			n.Attr = append(n.Attr, idAttr)
		}
		if n.Type == html.ElementNode {
			if n.Data == "img" || n.Data == "image" {
				log.Println("Encoding image")
				for i, v := range n.Attr {
					if v.Key == "src" {
						src := n.Attr[i]
						imageStream := mappedZipFile[filepath.Base(src.Val)]
						rc, err := imageStream.Open()
						if err != nil {
							return err
						}
						data, err := io.ReadAll(rc)
						if err != nil {
							return err
						}

						// Encode the image data as a base64 string
						encoded := base64.StdEncoding.EncodeToString(data)
						rc.Close()
						src.Val = fmt.Sprintf("data:image/jpeg;base64,%s", encoded)
						n.Attr[i].Val = src.Val
					}
					if v.Key == "href" && v.Namespace == "xlink" {
						src := n.Attr[i]

						imageStream := mappedZipFile[filepath.Base(src.Val)]
						rc, err := imageStream.Open()
						if err != nil {
							return err
						}
						data, err := io.ReadAll(rc)
						if err != nil {
							return err
						}
						ext := ".jpeg"
						allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
						if contains(allowedExtensions, filepath.Ext(src.Val)) {
							ext = filepath.Ext(src.Val)
							ext = strings.TrimLeft(ext, ".")
						} else {
							// TODO : add alt text if img, if svg add title and desc element as child
							ext = strings.TrimLeft(ext, ".")
						}
						// Encode the image data as a base64 string
						encoded := base64.StdEncoding.EncodeToString(data)
						rc.Close()
						src.Val = fmt.Sprintf("data:image/%s;base64,%s", ext, encoded)
						n.Attr[i].Val = src.Val
					}
				}

			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverseBody(c)
		}
		return nil
	}
	if page.Meta["type"] == "nav" {
		log.Println("Modifying TOC")
		app.EnsureNav(bodyNode)
	}
	log.Printf("Processing body  of %s\n", page.Page.Name)
	err := traverseBody(bodyNode)
	if err != nil {
		return err
	}
	log.Println("Finished processing")
	return nil
}
