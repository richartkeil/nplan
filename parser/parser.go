package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/richartkeil/nplan/core"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ParseNmap(path string) core.Scan {
	data, err := os.ReadFile(path)
	check(err)

	var scan Scan
	xml.Unmarshal(data, &scan)

	return convertScan(scan)
}

func ParseScan6(path string) []core.Host {
	data, err := os.ReadFile(path)
	check(err)

	output := strings.TrimSpace(string(data))

	// If scan has global and local addresses, we only want global addresses
	if strings.Contains(output, "Global addresses") {
		output = strings.Split(output, "Global addresses:\n")[1]
	}
	lines := strings.Split(output, "\n")

	var hosts []core.Host
	for _, line := range lines {
		cols := strings.Split(line, " @ ")
		hosts = append(hosts, core.Host{
			IPv6: cols[0],
			MAC:  strings.ToUpper(cols[1]),
		})
	}

	return hosts
}

func convertScan(scan Scan) core.Scan {
	var hosts []core.Host
	for _, host := range scan.Hosts {
		hosts = append(hosts, convertHost(host))
	}

	return core.Scan{
		Hosts: hosts,
	}
}

func convertHost(nmapHost Host) core.Host {
	var host core.Host

	for _, address := range nmapHost.Address {
		if address.Type == "ipv4" {
			host.IPv4 = address.Value
		} else if address.Type == "ipv6" {
			host.IPv6 = address.Value
		} else if address.Type == "mac" {
			host.MAC = address.Value
		}
	}

	for _, hostname := range nmapHost.Hostnames {
		host.Hostname = hostname.Name
	}

	for _, port := range nmapHost.Ports {
		host.Ports = append(host.Ports, convertPort(port))
	}

	host.Hops = nmapHost.Distance.Value
	host.OS = getHostOS(nmapHost.OSMatches)

	return host
}

func convertPort(nmapPort Port) core.Port {
	version := ""
	if nmapPort.Service.Version != "" || nmapPort.Service.Product != "" {
		version = fmt.Sprintf("%v %v", nmapPort.Service.Product, nmapPort.Service.Version)
	}

	port := core.Port{
		Protocol:       nmapPort.Protocol,
		Number:         nmapPort.Portid,
		ServiceName:    nmapPort.Service.Name,
		ServiceVersion: version,
	}
	for _, table := range nmapPort.Tables {
		port.HostKeys = append(port.HostKeys, convertKey(table))
	}
	return port
}

func getHostOS(nmapOSMatches []OS) string {
	if len(nmapOSMatches) < 1 {
		return ""
	}
	// Nmap sorts OS matches by accuracy, so we take the first one:
	match := nmapOSMatches[0]
	return fmt.Sprintf("%v (%v%%)", match.Name, match.Accuracy)
}

func convertKey(nmapTable Table) core.HostKey {
	var key core.HostKey
	for _, element := range nmapTable.Elements {
		if element.Key == "type" {
			key.Type = element.Value
		}
		if element.Key == "key" {
			key.Key = element.Value
		}
		if element.Key == "fingerprint" {
			key.Fingerprint = element.Value
		}
	}
	return key
}
