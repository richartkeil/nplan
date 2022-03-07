package main

import (
	"encoding/json"
	"errors"
	"flag"
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
	nmapInputFlag := flag.String("nmap", "", "nmap input file path")
	scan6InputFlag := flag.String("scan6", "", "scan6 input file path")
	jsonFileFlag := flag.String("json", "./dist/scan.json", "intermediate json file path")
	drawioOutputFlag := flag.String("drawio", "./dist/plan.drawio", "drawio output file path")

	flag.Parse()

	jsonData, err := os.ReadFile(*jsonFileFlag)
	if errors.Is(err, os.ErrNotExist) {
		jsonData = []byte("{}")
		err = os.WriteFile(*jsonFileFlag, jsonData, 0644)
		check(err)
	}

	// Read existing JSON
	var scan core.Scan
	err = json.Unmarshal(jsonData, &scan)
	check(err)

	if *nmapInputFlag != "" {
		nmapScan := parser.ParseNmap(*nmapInputFlag)
		scan = nmapScan
		// core.ComplementWithNmap(&scan, &nmapScan)
	}
	if *scan6InputFlag != "" {
		scan6Hosts := parser.ParseScan6(*scan6InputFlag)
		core.ComplementWithIPv6(&scan, &scan6Hosts)
	}

	json, err := json.MarshalIndent(scan, "", "  ")
	check(err)

	os.WriteFile(*jsonFileFlag, json, 0644)

	if *drawioOutputFlag != "" {
		exporter.Export(*drawioOutputFlag, &scan)
	}
}
