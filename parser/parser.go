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

	lowerPart := strings.Split(string(data), "Global addresses:\n")[1]
	lines := strings.Split(lowerPart, "\n")

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

	return host
}

func convertPort(nmapPort Port) core.Port {
	return core.Port{
		Protocol:       nmapPort.Protocol,
		Number:         nmapPort.Portid,
		ServiceName:    nmapPort.Service.Name,
		ServiceVersion: fmt.Sprintf("%v %v", nmapPort.Service.Product, nmapPort.Service.Version),
	}
}
