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

var defaultModelPath = "./dist/model.json"
var defaultDrawioPath = "./dist/drawio.json"

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
	jsonFileFlag := flag.String("json", defaultModelPath, "Set the path where to store the .json model file.")
	drawioOutputFlag := flag.String("drawio", defaultDrawioPath, "Set the path were to store the exported .drawio file.")
	resetModelFlag := flag.Bool("fresh", false, "Delete the previous .json model and build a new one. Use with caution.")

	flag.Parse()

	createDistFolder(jsonFileFlag, drawioOutputFlag)
	scan := loadModel(jsonFileFlag, resetModelFlag)

	if *nmapInputFlag != "" {
		nmapScan := parser.ParseNmap(*nmapInputFlag)
		core.ComplementWithNmap(scan, &nmapScan)
	}
	if *scan6InputFlag != "" {
		scan6Hosts := parser.ParseScan6(*scan6InputFlag)
		core.ComplementWithIPv6(scan, &scan6Hosts)
	}

	saveModel(jsonFileFlag, scan)

	if *exportFlag {
		exporter.Export(*drawioOutputFlag, scan)
	}
}

func createDistFolder(modelPath *string, drawioPath *string) {
	if *modelPath == defaultModelPath && *drawioPath == defaultDrawioPath {
		err := os.MkdirAll("./dist", 0755)
		check(err)
	}
}

func loadModel(modelPath *string, resetModel *bool) *core.Scan {
	jsonData, err := os.ReadFile(*modelPath)
	if errors.Is(err, os.ErrNotExist) || *resetModel {
		jsonData = []byte("{}")
		err = os.WriteFile(*modelPath, jsonData, 0644)
		check(err)
	}
	var scan core.Scan
	err = json.Unmarshal(jsonData, &scan)
	check(err)
	return &scan
}

func saveModel(modelPath *string, scan *core.Scan) {
	json, err := json.MarshalIndent(scan, "", "  ")
	check(err)
	os.WriteFile(*modelPath, json, 0644)
}
