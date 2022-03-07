package main

import (
	"encoding/json"
	"os"

	"github.com/richartkeil/nplan/core"
	"github.com/richartkeil/nplan/exporter"
	"github.com/richartkeil/nplan/parser"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	scan := parser.ParseNmap("./scans/scan_with_mac.xml")

	hosts := parser.ParseScan6("./scans/scan6.txt")
	core.ComplementWithIPv6(&scan, &hosts)

	json, err := json.MarshalIndent(scan, "", "  ")
	check(err)

	os.WriteFile("./dist/scan.json", json, 0644)
	exporter.Export(scan)
}
