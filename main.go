package main

import (
	"bytes"
	"fmt"
	"github.com/ianr0bkny/go-sonos/ssdp"
	"log"
)

func main() {
	writer := new(bytes.Buffer)
	log.SetOutput(writer)
	mgr := ssdp.MakeManager()
	defer mgr.Close()
	mgr.Discover("en0", "11209", false)
	qry := ssdp.ServiceQueryTerms{ssdp.ServiceKey("schemas-upnp-org-MusicServices"): -1}
	result := mgr.QueryServices(qry)
	if devices, has := result[ssdp.ServiceKey("schemas-upnp-org-MusicServices")]; has {
		fmt.Printf("Found %d devices\n", len(devices))
		for _, dev := range devices {
			fmt.Printf("%s %s %s %s %s\n", dev.Product(), dev.ProductVersion(), dev.Name(), dev.Location(), dev.UUID())
		}
	}
}
