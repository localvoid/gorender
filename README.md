# gorender

CLI util to render templates with [Go](http://golang.org)
[text.Template](http://golang.org/pkg/text/template/) template engine.

## Install

```sh
$ go install github.com/localvoid/gorender
```

## Usage

```sh
$ gorender -d data.json template.txt
```

## Flags

### -d [file]

Type: `String`  

Data file in JSON format.

### -i [pattern]

Type: `String`

Include templates.

### -b [name]

Type: `String`

Base template name.

### -html

Type: `Boolean`  
Default: `false`

Use [html.Template](http://golang.org/pkg/text/template/) package to
render templates.
