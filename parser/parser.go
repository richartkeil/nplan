package parser

import (
	"encoding/xml"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Parse(path string) Scan {
	data, err := os.ReadFile(path)
	check(err)

	var scan Scan
	xml.Unmarshal(data, &scan)

	return scan
}
