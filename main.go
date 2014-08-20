package main

import (
	"encoding/json"
	"flag"
	"fmt"
	htmlTemplate "html/template"
	"io/ioutil"
	"os"
	textTemplate "text/template"
)

var dataPath = flag.String("d", "data.json", "JSON data file")
var useHtml = flag.Bool("html", false, "use `html.Template`")

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gorender -d [datafile] [templates...] \n\n")
	flag.PrintDefaults()
	os.Exit(64)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	templates := flag.Args()
	if len(templates) == 0 {
		usage()
	}

	if len(*dataPath) == 0 {
		usage()
	}

	for _, t := range templates {
		if _, err := os.Stat(t); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Template file \"%s\" is missing\n", t)
			os.Exit(66)
		}
	}

	if _, err := os.Stat(*dataPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Data file \"%s\" is missing\n", *dataPath)
		os.Exit(66)
	}

	var tpl interface{}
	var err error

	if *useHtml {
		tpl, err = htmlTemplate.ParseFiles(templates...)
	} else {
		tpl, err = textTemplate.ParseFiles(templates...)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(65)
	}

	raw, err := ioutil.ReadFile(*dataPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var data map[string]interface{}

	err = json.Unmarshal(raw, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid data file: %s\n", err)
		os.Exit(65)
	}

	if *useHtml {
		t := tpl.(*htmlTemplate.Template)
		err = t.Execute(os.Stdout, data)
	} else {
		t := tpl.(*textTemplate.Template)
		err = t.Execute(os.Stdout, data)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(65)
	}

	os.Exit(0)
}
