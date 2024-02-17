package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"

	"github.com/noncombatant/content2html"
)

func main() {
	templatePathname := flag.String("template", "template.html", "Pathname of template file")
	outputDirname := flag.String("out", "", "Pathname of output directory")
	flag.Parse()

	template, e := template.ParseFiles(*templatePathname)
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}

	for _, pathname := range flag.Args() {
		if e := content2html.GenerateHTMLFile(template, pathname, *outputDirname); e != nil {
			fmt.Fprintln(os.Stderr, e)
		}
	}
}
