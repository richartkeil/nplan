package core

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
