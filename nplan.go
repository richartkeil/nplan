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
	// Common flags
	nmapInputFlag := flag.String("nmap", "", "Set the path to the nmap input .xml file.")
	scan6InputFlag := flag.String("scan6", "", "Set the path to the scan6 input .txt file. For this to take effect the current model should already include MAC addresses.")
	exportFlag := flag.Bool("export", false, "Export the current model to a .drawio file.")
	// Config flags
	jsonFileFlag := flag.String("json", "./dist/model.json", "Set the path where to store the .json model file.")
	drawioOutputFlag := flag.String("drawio", "./dist/plan.drawio", "Set the path were to store the exported .drawio file.")
	resetModelFlag := flag.Bool("fresh", false, "Delete the previous .json model and build a new one. Use with caution.")

	flag.Parse()

	jsonData, err := os.ReadFile(*jsonFileFlag)
	if errors.Is(err, os.ErrNotExist) || *resetModelFlag {
		jsonData = []byte("{}")
		err = os.WriteFile(*jsonFileFlag, jsonData, 0644)
		check(err)
	}

	var scan core.Scan
	err = json.Unmarshal(jsonData, &scan)
	check(err)

	if *nmapInputFlag != "" {
		nmapScan := parser.ParseNmap(*nmapInputFlag)
		core.ComplementWithNmap(&scan, &nmapScan)
	}
	if *scan6InputFlag != "" {
		scan6Hosts := parser.ParseScan6(*scan6InputFlag)
		core.ComplementWithIPv6(&scan, &scan6Hosts)
	}

	json, err := json.MarshalIndent(scan, "", "  ")
	check(err)

	os.WriteFile(*jsonFileFlag, json, 0644)

	if *exportFlag {
		exporter.Export(*drawioOutputFlag, &scan)
	}
}
