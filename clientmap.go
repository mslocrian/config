package main

var ClientInterfaces = map[string]ClientIf{"ribd": &RIBDClient{},
	"asicd":      &ASICDClient{},
	"arpd":       &ARPDClient{},
	"bgpd":       &BgpDClient{},
	"lacpd":      &LACPDClient{},
	"dhcprelayd": &DHCPRELAYDClient{},
	"local":      &LocalClient{},
	"ospfd":      &OSPFDClient{},
	"stpd":       &STPDClient{},
	"bfdd":       &BFDDClient{},
	"vrrpd":      &VRRPDClient{},
	"sysd":       &SYSDClient{},
}
