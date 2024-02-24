package exporter

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/richartkeil/nplan/core"
)

// Keys
var keyHeight = 50
var keyWidth = 23

// Hosts
var rows = 8
var hostWidth = 260
var hostHeight = 160
var additionalHeightPerPort = 20
var hostKeyOffsetX = 235
var hostKeyOffsetY = 5
var hostKeyPadding = 10
var hostPadding = 30

// Duplicate Fingerprint hosts display
var dupHostsFingerprintX = -440
var dupHostsFingerprintY = 0
var dupHostsFingerprintWidth = 350
var dupHostsFingerprintBaseHeight = 70
var dupHostsInsetX = 50
var dupHostsKeyOffsetX = 0
var dupHostsKeyOffsetY = 24

// Unidentified hosts
var unidentifiedHostsX = -700
var unidentifiedHostsY = 0
var unidentifiedHostsWidth = 260
var unidentifiedHostsHeight = 100

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

	keyColorMap := findDuplicateHostKeyColors(scan)
	cells = addHosts(cells, scan, keyColorMap)
	cells = addDuplicateHostKeys(cells, scan, keyColorMap)
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

func findDuplicateHostKeyColors(scan *core.Scan) map[core.HostKey]string {
	// Collect number of hosts that share each key:
	hostsPerKey := make(map[core.HostKey]int)
	for _, host := range scan.Hosts {
		for _, port := range host.Ports {
			for _, hostKey := range port.HostKeys {
				if hostKey.Fingerprint != "" {
					hostsPerKey[hostKey] += 1
				}
			}
		}
	}

	// Generate a unique color palette for all keys that are shared by more than one host
	duplicateHostCount := 0
	for _, hosts := range hostsPerKey {
		if hosts > 1 {
			duplicateHostCount++
		}
	}

	palette := colorful.FastHappyPalette(duplicateHostCount)
	keyColorMap := make(map[core.HostKey]string)

	// Assign a color to each duplicate key
	keyIndex := 0
	for key, hosts := range hostsPerKey {
		if hosts <= 1 {
			continue
		}
		keyColorMap[key] = palette[keyIndex].Hex()
		keyIndex += 1
	}
	return keyColorMap
}

func addHosts(cells []MxCell, scan *core.Scan, keyColorMap map[core.HostKey]string) []MxCell {
	currentX := 0
	currentY := 0
	for i, host := range scan.Hosts {
		id := uuid.NewString()
		cells = append(cells, MxCell{
			Id:     id,
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

		// Add colored keys
		keyCount := 0
		for _, port := range host.Ports {
			for _, hostKey := range port.HostKeys {
				if keyColorMap[hostKey] != "" {
					color := keyColorMap[hostKey]
					cells = append(cells, makeKeyIconCell(id, color, hostKeyOffsetX, hostKeyOffsetY+(keyHeight+hostKeyPadding)*keyCount))
					keyCount += 1
				}
			}
		}

		currentY += getHostHeight(&host) + hostPadding
		if (i+1)%rows == 0 {
			currentX += hostWidth + hostPadding
			currentY = 0
		}
	}

	return cells
}

func addDuplicateHostKeys(cells []MxCell, scan *core.Scan, keyColorMap map[core.HostKey]string) []MxCell {
	// For each group of hosts with the same Fingerprint create a box
	currentX := dupHostsFingerprintX
	currentY := dupHostsFingerprintY
	for key, color := range keyColorMap {
		value := fmt.Sprintf("<u>Identical SSH Key:</u><br>Type: <strong>%v</strong><br>Fingerprint: <strong>%v</strong>", key.Type, key.Fingerprint)
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
		cells = append(cells, makeKeyIconCell(id, color, dupHostsKeyOffsetX, dupHostsKeyOffsetY))
		currentY += dupHostsFingerprintBaseHeight + hostPadding
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
		currentY += unidentifiedHostsHeight + hostPadding
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

func makeKeyIconCell(parentId string, color string, x int, y int) MxCell {
	return MxCell{
		Id:     uuid.NewString(),
		Value:  "",
		Parent: parentId,
		Style:  fmt.Sprintf("shape=mxgraph.cisco19.key;fillColor=%v;strokeColor=none;rotation=90", color),
		Vertex: "1",
		MxGeometry: &MxGeometry{
			X:      fmt.Sprint(x),
			Y:      fmt.Sprint(y),
			Width:  fmt.Sprint(keyHeight),
			Height: fmt.Sprint(keyWidth),
			As:     "geometry",
		},
	}
}
