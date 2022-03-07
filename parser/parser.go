package parser

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/richartkeil/nplan/core"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Parse(path string) core.Scan {
	data, err := os.ReadFile(path)
	check(err)

	var scan Scan
	xml.Unmarshal(data, &scan)

	return convertScan(scan)
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
		Protocol:    nmapPort.Protocol,
		Number:      nmapPort.Portid,
		ServiceName: nmapPort.Service.Name,
		// ServiceVersion: nmapPort.Service.Product,
		ServiceVersion: fmt.Sprintf("%v %v", nmapPort.Service.Product, nmapPort.Service.Version),
	}
}
