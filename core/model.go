package core

type Scan struct {
	Hosts []Host `json:"hosts"`
}

type Host struct {
	IPv4     string `json:"ipv4"`
	IPv6     string `json:"ipv6"`
	MAC      string `json:"mac"`
	Hostname string `json:"hostname"`
	Ports    []Port `json:"ports"`
	OS       string `json:"os,omitempty"`
	Hops 	 int 	`json:"distance,omitempty"`
}

type Port struct {
	Protocol       string    `json:"protocol"`
	Number         int       `json:"number"`
	ServiceName    string    `json:"service_name"`
	ServiceVersion string    `json:"service_version"`
	HostKeys       []HostKey `json:"host_keys,omitempty"`
}

type HostKey struct {
	Type        string `json:"type"`
	Key         string `json:"key"`
	Fingerprint string `json:"fingerprint"`
}
