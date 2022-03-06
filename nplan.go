package main

import (
	"encoding/json"
	"os"

	"github.com/richartkeil/nplan/parser"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	scan := parser.Parse("./scans/scan.xml")

	json, err := json.MarshalIndent(scan, "", "  ")
	check(err)

	os.WriteFile("./dist/scan.json", json, 0644)
}
