package exporter

import "encoding/xml"

type MxFile struct {
	XMLName xml.Name `xml:"mxfile"`
	Diagram *Diagram `xml:"diagram"`
}

type Diagram struct {
	XMLName      xml.Name      `xml:"diagram"`
	Id           string        `xml:"id,attr"`
	Name         string        `xml:"name,attr"`
	MxGraphModel *MxGraphModel `xml:"mxGraphModel"`
}

type MxGraphModel struct {
	XMLName  xml.Name `xml:"mxGraphModel"`
	Root     *Root    `xml:"root"`
	Dx       string   `xml:"dx,attr"`
	Dy       string   `xml:"dy,attr"`
	Grid     string   `xml:"grid,attr"`
	GridSize string   `xml:"gridSize,attr"`
	Guides   string   `xml:"guides,attr"`
	Tooltips string   `xml:"tooltips,attr"`
	Connect  string   `xml:"connect,attr"`
	Arrows   string   `xml:"arrows,attr"`
}

type Root struct {
	XMLName xml.Name `xml:"root"`
	MxCells []MxCell `xml:"mxCell"`
}

type MxCell struct {
	XMLName    xml.Name    `xml:"mxCell"`
	Id         string      `xml:"id,attr"`
	Value      string      `xml:"value,attr,omitempty"`
	Style      string      `xml:"style,attr,omitempty"`
	MxGeometry *MxGeometry `xml:"mxGeometry,omitempty"`
	Parent     string      `xml:"parent,attr,omitempty"`
	Vertex     string      `xml:"vertex,attr,omitempty"`
}

type MxGeometry struct {
	XMLName xml.Name `xml:"mxGeometry"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	As      string   `xml:"as,attr"`
}
