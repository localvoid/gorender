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

var (
	dataPath = flag.String("d", "", "JSON data file")
	include  = flag.String("i", "", "Include templates")
	baseName = flag.String("b", "", "Base template name")
	useHtml  = flag.Bool("html", false, "use html.Template package for rendering")
)

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

	if len(*include) > 0 {
		if *useHtml {
			t := tpl.(*htmlTemplate.Template)
			tpl, err = t.ParseGlob(*include)
		} else {
			t := tpl.(*textTemplate.Template)
			tpl, err = t.ParseGlob(*include)
		}
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(65)
	}

	var data map[string]interface{}

	if len(*dataPath) > 0 {
		raw, err := ioutil.ReadFile(*dataPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		err = json.Unmarshal(raw, &data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid data file: %s\n", err)
			os.Exit(65)
		}
	}

	if *useHtml {
		t := tpl.(*htmlTemplate.Template)
		if len(*baseName) > 0 {
			err = t.ExecuteTemplate(os.Stdout, *baseName, data)
		} else {
			err = t.Execute(os.Stdout, data)
		}
	} else {
		t := tpl.(*textTemplate.Template)
		if len(*baseName) > 0 {
			err = t.ExecuteTemplate(os.Stdout, *baseName, data)
		} else {
			err = t.Execute(os.Stdout, data)
		}
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(65)
	}

	os.Exit(0)
}
