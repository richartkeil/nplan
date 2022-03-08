package parser

import "encoding/xml"

type Scan struct {
	XMLName xml.Name `xml:"nmaprun" json:"-"`
	Hosts   []Host   `xml:"host" json:"hosts"`
}

type Host struct {
	XMLName   xml.Name   `xml:"host" json:"-"`
	Address   []Address  `xml:"address" json:"address"`
	Hostnames []Hostname `xml:"hostnames>hostname" json:"hostnames"`
	Ports     []Port     `xml:"ports>port" json:"ports"`
}

type Address struct {
	XMLName xml.Name `xml:"address" json:"-"`
	Value   string   `xml:"addr,attr" json:"value"`
	Type    string   `xml:"addrtype,attr" json:"type"`
}

type Hostname struct {
	XMLName xml.Name `xml:"hostname" json:"-"`
	Name    string   `xml:"name,attr" json:"name"`
}

type Port struct {
	XMLName  xml.Name `xml:"port" json:"-"`
	Protocol string   `xml:"protocol,attr" json:"protocol"`
	Portid   int      `xml:"portid,attr" json:"portid"`
	Service  Service  `xml:"service" json:"service"`
	Tables   []Table  `xml:"script>table" json:"keys"`
}

type Service struct {
	XMLName xml.Name `xml:"service" json:"-"`
	Name    string   `xml:"name,attr" json:"name"`
	Product string   `xml:"product,attr" json:"product"`
	Version string   `xml:"version,attr" json:"version"`
}

type Table struct {
	XMLName  xml.Name  `xml:"table" json:"-"`
	Elements []Element `xml:"elem" json:"elements"`
}

type Element struct {
	XMLName xml.Name `xml:"elem" json:"-"`
	Key     string   `xml:"key,attr" json:"key"`
	Value   string   `xml:",innerxml" json:"value"`
}
