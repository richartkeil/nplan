package exporter

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/richartkeil/nplan/parser"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Export(scan parser.Scan) {
	cols := 5
	width := 180
	height := 120
	padding := 30

	cells := make([]MxCell, 0)
	cells = append(cells, MxCell{
		Id: "0",
	})
	for i := 0; i < len(scan.Hosts); i++ {
		cells = append(cells, MxCell{
			Id:     uuid.NewString(),
			Value:  getHostValue(scan.Hosts[i]),
			Parent: "0",
			Style:  "rounded=1;whiteSpace=wrap;html=1;arcSize=2",
			Vertex: "1",
			MxGeometry: &MxGeometry{
				X:      fmt.Sprint((i % cols) * (width + padding)),
				Y:      fmt.Sprint((i / cols) * (height + padding)),
				Width:  fmt.Sprint(width),
				Height: fmt.Sprint(height),
				As:     "geometry",
			},
		})
	}

	mxFile := MxFile{
		Diagram: &Diagram{
			Id:   uuid.NewString(),
			Name: "Network Plan",
			MxGraphModel: &MxGraphModel{
				Root: &Root{
					MxCells: cells,
				},
				Dx:       "3000",
				Dy:       "2000",
				Grid:     "1",
				GridSize: "10",
				Guides:   "1",
				Tooltips: "1",
				Connect:  "1",
				Arrows:   "1",
			},
		},
	}

	output, err := xml.MarshalIndent(mxFile, "", "  ")
	check(err)

	os.WriteFile("./dist/drawio.xml", output, 0644)
}

func getHostValue(host parser.Host) string {
	serviceColor := "#bbb"

	value := ""
	if len(host.Hostnames) > 0 {
		value += fmt.Sprintf("<strong>%v</strong><br>", host.Hostnames[0].Name)
	}
	value += fmt.Sprintf("%v<br><br>", host.Address.Value)

	for i := 0; i < len(host.Ports); i++ {
		value += fmt.Sprintf(
			":%v - %v<br><span style=\"color: %v\">(%v %v)</span><br>",
			host.Ports[i].Portid,
			host.Ports[i].Service.Name,
			serviceColor,
			host.Ports[i].Service.Product,
			host.Ports[i].Service.Version,
		)
	}

	return value
}
