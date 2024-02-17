package exporter

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"os"

	"github.com/google/uuid"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/richartkeil/nplan/core"
)

// Hosts
var rows = 8
var hostWidth = 260
var hostHeight = 160
var additionalHeightPerPort = 20
var padding = 30

// Duplicate Fingerprint hosts display
var dupHostsFingerprintX = -440
var dupHostsFingerprintY = 0
var dupHostsFingerprintWidth = 350
var dupHostsFingerprintBaseHeight = 70
var dupHostsInsetX = 50
var dupHostsKeyOffsetX float32 = 0
var dupHostsKeyOffsetY float32 = 23.75

// Unidentified hosts
var unidentifiedHostsX = -700
var unidentifiedHostsY = 0
var unidentifiedHostsWidth = 260
var unidentifiedHostsHeight = 100

var hostGroups = make(map[core.HostKey][]core.Host)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Export(path string, scan *core.Scan) {
	cells := make([]MxCell, 0)
	cells = append(cells, MxCell{
		Id: "0",
	})
	cells = append(cells, MxCell{
		Id:     "1",
		Parent: "0",
	})
	cells = addHosts(cells, scan)
	cells = addHostsWithSameFingerprint(cells, scan)
	cells = addUnidentifiedHosts(cells, scan)

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

	os.WriteFile(path, output, 0644)
}

func addHosts(cells []MxCell, scan *core.Scan) []MxCell {
	currentX := 0
	currentY := 0
	for i, host := range scan.Hosts {
		cells = append(cells, MxCell{
			Id:     uuid.NewString(),
			Value:  getHostValue(host),
			Parent: "1",
			Style:  "rounded=1;whiteSpace=wrap;html=1;arcSize=2",
			Vertex: "1",
			MxGeometry: &MxGeometry{
				X:      fmt.Sprint(currentX),
				Y:      fmt.Sprint(currentY),
				Width:  fmt.Sprint(hostWidth),
				Height: fmt.Sprint(getHostHeight(&host)),
				As:     "geometry",
			},
		})
		currentY += getHostHeight(&host) + padding
		if (i+1)%rows == 0 {
			currentX += hostWidth + padding
			currentY = 0
		}
	}

	return cells
}

func addHostsWithSameFingerprint(cells []MxCell, scan *core.Scan) []MxCell {
	// Group hosts by Fingerprint address
	for _, host := range scan.Hosts {
		for _, port := range host.Ports {
			for _, hostKey := range port.HostKeys {
				if hostKey.Fingerprint != "" {
					hostGroups[hostKey] = append(hostGroups[hostKey], host)

				}
			}
		}
	}

	shuffledPallete := generateColorPallete()
	index := 0

	// For each group of hosts with the same Fingerprint create a box
	currentX := dupHostsFingerprintX
	currentY := dupHostsFingerprintY
	for duplicateKey, hosts := range hostGroups {
		duplicateKey.Color = shuffledPallete[index].Hex()

		// Do not show Fingerprints with only one host:
		if len(hosts) <= 1 {
			continue
		}


		value := fmt.Sprintf("<u>Identical SSH Key:</u><br>Type: <strong>%v</strong><br>Fingerprint: <strong>%v</strong>", duplicateKey.Type, duplicateKey.Fingerprint)
		id := uuid.NewString()
		cells = append(cells, MxCell{
			Id:     id,
			Value:  value,
			Parent: "1",
			Style:  fmt.Sprintf("rounded=1;whiteSpace=wrap;html=1;arcSize=2;align=left;spacingLeft=%v", dupHostsInsetX),
			Vertex: "1",
			MxGeometry: &MxGeometry{
				X:      fmt.Sprint(currentX),
				Y:      fmt.Sprint(currentY),
				Width:  fmt.Sprint(dupHostsFingerprintWidth),
				Height: fmt.Sprint(dupHostsFingerprintBaseHeight),
				As:     "geometry",
			},
		})
		cells = append(cells, makeKeyCell(id, duplicateKey.Color, dupHostsKeyOffsetX, dupHostsKeyOffsetY, "rotation=90;"))
		currentY += dupHostsFingerprintBaseHeight + padding
		index += 1
	}
	return cells
}

func addUnidentifiedHosts(cells []MxCell, scan *core.Scan) []MxCell {
	currentX := unidentifiedHostsX
	currentY := unidentifiedHostsY
	for _, host := range scan.UnidentifiedHosts {
		cells = append(cells, MxCell{
			Id:     uuid.NewString(),
			Value:  fmt.Sprintf("Unidentified host:<br><br>MAC: %v<br>IPv6: %v", host.MAC, host.IPv6),
			Parent: "1",
			Style:  "rounded=1;whiteSpace=wrap;html=1;arcSize=2",
			Vertex: "1",
			MxGeometry: &MxGeometry{
				X:      fmt.Sprint(currentX),
				Y:      fmt.Sprint(currentY),
				Width:  fmt.Sprint(unidentifiedHostsWidth),
				Height: fmt.Sprint(unidentifiedHostsHeight),
				As:     "geometry",
			},
		})
		currentY += unidentifiedHostsHeight + padding
	}
	return cells
}

func getHostHeight(host *core.Host) int {
	return hostHeight + len(host.Ports)*additionalHeightPerPort
}

func getHostValue(host core.Host) string {
	serviceColor := "#bbb"
	headerFontSize := 16
	value := ""

	// Addresses
	if host.Hostname != "" {
		value += fmt.Sprintf("<i>%v</i><br><br>", host.Hostname)
	}
	if host.IPv4 != "" {
		value += fmt.Sprintf(
			"<strong style=\"font-size: %vpx\">%v</strong><br>",
			headerFontSize,
			host.IPv4,
		)
	}
	if host.IPv6 != "" {
		value += fmt.Sprintf("IPv6: %v<br>", host.IPv6)
	}
	if host.MAC != "" {
		value += fmt.Sprintf("MAC: %v<br>", host.MAC)
	}

	// Ports
	if len(host.Ports) > 0 {
		value += "<br>"
	}
	for _, port := range host.Ports {
		value += fmt.Sprintf(":%v - %v<br>", port.Number, port.ServiceName)
		if port.ServiceVersion != "" {
			value += fmt.Sprintf(
				"<span style=\"color: %v\">(%v)</span><br>",
				serviceColor,
				port.ServiceVersion,
			)
		}
		for _, hostKey := range port.HostKeys {
			value += fmt.Sprintf(
				"<span style=\"color: %v\">(Key: %v)</span><br>",
				serviceColor,
				hostKey.Fingerprint,
			)
		}
	}

	// Misc
	if (host.OS != "") || (host.Hops != 0) {
		value += "<br>"
	}
	if host.OS != "" {
		value += fmt.Sprintf("OS: %v<br>", host.OS)
	}
	if host.Hops != 0 {
		value += fmt.Sprintf("Hops: %v<br>", host.Hops)
	}

	return value
}

func makeKeyCell(parentId string, color string, x float32, y float32, styleArgs string) MxCell {
	return MxCell{
		Id:     uuid.NewString(),
		Value:  "",
		Parent: parentId,
		Style:  fmt.Sprintf("shape=mxgraph.cisco19.key;fillColor=%v;strokeColor=none;%v", color, styleArgs),
		Vertex: "1",
		MxGeometry: &MxGeometry{
			X:      fmt.Sprint(x),
			Y:      fmt.Sprint(y),
			Width:  "50",
			Height: "22.5",
			As:     "geometry",
		},
	}
}

// Generate a color pallete for the keys and shuffle it
func generateColorPallete() []colorful.Color {
	pallete := colorful.FastHappyPalette(len(hostGroups))
	shuffledPallete := make([]colorful.Color, len(pallete))
	perm := rand.Perm(len(pallete))
	for i, v := range perm {
			shuffledPallete[v] = pallete[i]
	}
	return shuffledPallete
}
