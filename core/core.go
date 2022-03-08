package core

// Merges an existing scan with hosts obtained from a scan6 file.
func ComplementWithIPv6(scan *Scan, ipv6Hosts *[]Host) *Scan {
	for i, existingHost := range scan.Hosts {
		for _, ipv6Host := range *ipv6Hosts {
			if existingHost.MAC == ipv6Host.MAC {
				scan.Hosts[i].IPv6 = ipv6Host.IPv6
			}
		}
	}
	return scan
}

// Merges an existing scan with a new scan. For that, we check for all new
// hosts if there is already a host with the same IPv4. If there is, we merge
// the two hosts, prioritizing information from the new host. If there is not,
// we just add the new host.
func ComplementWithNmap(scan *Scan, newScan *Scan) *Scan {
	for _, newHost := range newScan.Hosts {
		existingHost := findHostByIPv4(scan, newHost.IPv4)
		if existingHost != nil {
			mergeHost(existingHost, newHost)
		} else {
			scan.Hosts = append(scan.Hosts, newHost)
		}
	}
	return scan
}

func findHostByIPv4(scan *Scan, ipv4 string) *Host {
	for i, existingHost := range scan.Hosts {
		if existingHost.IPv4 == ipv4 {
			return &scan.Hosts[i]
		}
	}
	return nil
}

func mergeHost(hostPtr *Host, newHost Host) *Host {
	if newHost.IPv6 != "" {
		hostPtr.IPv6 = newHost.IPv6
	}
	if newHost.MAC != "" {
		hostPtr.MAC = newHost.MAC
	}
	if newHost.Hostname != "" {
		hostPtr.Hostname = newHost.Hostname
	}

	for _, newPort := range newHost.Ports {
		existingPort := findPortByNumber(hostPtr, newPort.Number)
		if existingPort != nil {
			mergePort(existingPort, newPort)
		} else {
			hostPtr.Ports = append(hostPtr.Ports, newPort)
		}
	}
	return hostPtr
}

func findPortByNumber(host *Host, number int) *Port {
	for i, port := range host.Ports {
		if port.Number == number {
			return &host.Ports[i]
		}
	}
	return nil
}

func mergePort(port *Port, newPort Port) *Port {
	if newPort.Protocol != "" {
		port.Protocol = newPort.Protocol
	}
	if newPort.ServiceName != "" {
		port.ServiceName = newPort.ServiceName
	}
	if newPort.ServiceVersion != "" {
		port.ServiceVersion = newPort.ServiceVersion
	}
	mergeKeys(port, newPort)
	return port
}

func mergeKeys(port *Port, newPort Port) {
	if port.HostKeys == nil {
		port.HostKeys = newPort.HostKeys
	} else {
		for _, hostKey := range newPort.HostKeys {
			existingKey := findKeyByType(port, hostKey)
			if existingKey != nil {
				*existingKey = hostKey
			} else {
				port.HostKeys = append(port.HostKeys, hostKey)
			}
		}
	}
}

func findKeyByType(port *Port, newKey HostKey) *HostKey {
	for i, key := range port.HostKeys {
		if key.Type == newKey.Type {
			return &port.HostKeys[i]
		}
	}
	return nil
}
