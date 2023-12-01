package exporter

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/richartkeil/nplan/core"
)

// Hosts
var rows = 8
var hostWidth = 260
var hostHeight = 160
var additionalHeightPerPort = 20
var padding = 30

// Duplicate Fingerprint hosts display
var dupHostsFingerprintX = -400
var dupHostsFingerprintY = 0
var dupHostsFingerprintWidth = 260
var dupHostsFingerprintHeightPerMac = 15
var dupHostsFingerprintBaseHeight = 70

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
	hostGroups := make(map[string][]core.Host)
	for _, host := range scan.Hosts {
		for _, port := range host.Ports {
			for _, hostKey := range port.HostKeys {
				if hostKey.Type == "ssh-rsa" && hostKey.Fingerprint != "" {
					hostGroups[hostKey.Fingerprint] = append(hostGroups[hostKey.Fingerprint], host)
				}
			}
		}
	}

	// For each group of hosts with the same Fingerprint create a box
	currentX := dupHostsFingerprintX
	currentY := dupHostsFingerprintY
	for mac, hosts := range hostGroups {
		// Do not show Fingerprints with only one host:
		if len(hosts) <= 1 {
			continue
		}

		value := fmt.Sprintf("Hosts with RSA Fingerprint<br><strong>%v</strong>:<br><br>", mac)
		for _, host := range hosts {
			value += fmt.Sprintf("%v<br>", host.IPv4)
		}
		cells = append(cells, MxCell{
			Id:     uuid.NewString(),
			Value:  value,
			Parent: "1",
			Style:  "rounded=1;whiteSpace=wrap;html=1;arcSize=2",
			Vertex: "1",
			MxGeometry: &MxGeometry{
				X:      fmt.Sprint(currentX),
				Y:      fmt.Sprint(currentY),
				Width:  fmt.Sprint(dupHostsFingerprintWidth),
				Height: fmt.Sprint(dupHostsFingerprintBaseHeight + len(hosts)*dupHostsFingerprintHeightPerMac),
				As:     "geometry",
			},
		})
		currentY += dupHostsFingerprintBaseHeight + len(hosts)*dupHostsFingerprintHeightPerMac + padding
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
			if hostKey.Type == "ssh-rsa" {
				value += fmt.Sprintf(
					"<span style=\"color: %v\">(RSA: %v)</span><br>",
					serviceColor,
					port.HostKeys[0].Fingerprint,
				)
			}
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
