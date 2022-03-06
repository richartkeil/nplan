package main

import (
	"github.com/richartkeil/nplan/exporter"
	"github.com/richartkeil/nplan/parser"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	scan := parser.Parse("./scans/scan.xml")
	exporter.Export(scan)
}
