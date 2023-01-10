package main

import (
	"flag"
	"io"
	"log"
	"os"
	"runtime/pprof"
)

type App struct {
	title    string
	port     string
	mode     string
	infoLog  *log.Logger
	errorLog *log.Logger
}

func (app *App) read() {
	zipReader := app.ReadEpub(app.title)
	err := app.CheckMime(zipReader.File[0])
	if err != nil {
		log.Fatal(err)
	}

	mappedZipFile, opf, err := app.MapContent(zipReader)
	if err != nil {
		log.Fatal(err)
	}

	structure, err := app.CreateStructure(opf)
	if err != nil {
		log.Fatal(err)
	}

	pageList := app.EnsurePageList(structure, mappedZipFile)

	bodyNode, err := app.GenerateNode(pageList, mappedZipFile)
	if err != nil {
		log.Fatal(err)
	}

	book, err := app.RenderBody(bodyNode, structure.Metadata)
	if err != nil {
		log.Fatal(err)
	}
	if app.mode == "reader" {
		app.Send(book)
	}
	zipReader.Close()
}

func main() {
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	title := flag.String("read", "", "read book")
	port := flag.String("port", "8080", "set port")
	mode := flag.String("mode", "", "set mode")
	verbose := flag.String("verbose", "discrete", "set log output noise")
	flag.Parse()

	if *title == "" {
		log.Panic("Unable to read without setting book name")
	}
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	if *verbose == "quiet" {
		infoLog.SetOutput(io.Discard)
		errorLog.SetOutput(io.Discard)
	}
	app := App{title: *title, port: *port, mode: *mode, infoLog: infoLog, errorLog: errorLog}
	app.read()
}
