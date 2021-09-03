package main

import (
	"flag"
	"go/parser"
	"go/token"
	"path"
	"regexp"

	"github.com/jreamy/go-opt/json/parse"
)

var trigger = regexp.MustCompile("// go-opt: (.* )?json( .*)?")

var dirname string

func main() {

	// Parse the file name from the command line
	// Unless it has already been set by the tests
	if dirname == "" {
		flag.StringVar(&dirname, "f", ".", "go directory to parse")
		flag.Parse()
	}

	// Parse the source files
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dirname, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	for _, pkg := range pkgs {
		p := &parse.Package{Package: pkg}
		p.Parse(trigger)

		optFilename := path.Join(dirname, pkg.Name+"_go_opt.go")
		optTestname := path.Join(dirname, pkg.Name+"_go_opt_test.go")

		JSONMarshalers(p, "buf")

		if err := UseTemplate(optFilename, optTemplate, p); err != nil {
			panic(err)
		}

		if err := UseTemplate(optTestname, testTemplate, p); err != nil {
			panic(err)
		}
	}
}
