package core

type Scan struct {
	Hosts             []Host             `json:"hosts"`
	UnidentifiedHosts []UnidentifiedHost `json:"unidentified_hosts"`
}

type UnidentifiedHost struct {
	IPv6 string `json:"ipv6"`
	MAC  string `json:"mac"`
}

type Host struct {
	IPv4     string `json:"ipv4"`
	IPv6     string `json:"ipv6"`
	MAC      string `json:"mac"`
	Hostname string `json:"hostname"`
	Ports    []Port `json:"ports"`
	Hops     int    `json:"hop_distance,omitempty"`
	OS       string `json:"os,omitempty"`
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
	Color 		 	string `json:"color"`
}
