package main

import (
	"bytes"
	"fmt"
	"github.com/ianr0bkny/go-sonos"
	"github.com/ianr0bkny/go-sonos/ssdp"
	"github.com/ianr0bkny/go-sonos/upnp"
	"log"
)

func main() {
	writer := new(bytes.Buffer)
	log.SetOutput(writer)
	mgr := ssdp.MakeManager()
	defer mgr.Close()
	if err := mgr.Discover("en0", "11209", false); err != nil {
		log.Fatalf("Error discovering devices: %s", err)
	}
	qry := ssdp.ServiceQueryTerms{ssdp.ServiceKey("schemas-upnp-org-MusicServices"): -1}
	result := mgr.QueryServices(qry)
	if devices, has := result[ssdp.ServiceKey("schemas-upnp-org-MusicServices")]; has {
		fmt.Printf("Found %d devices\n", len(devices))
		for _, d := range devices {
			fmt.Printf("%s %s %s %s %s\n", d.Product(), d.ProductVersion(), d.Name(), d.Location(), d.UUID())
			svc_map, err := upnp.Describe(d.Location())
			if nil != err {
				log.Fatalf("Error describing device: %s", err)
			}
			fmt.Printf("Found %d services\n", len(svc_map))
			for key, svcs := range svc_map {
				fmt.Printf("key: %s, services: %#v\n", key, svcs)
			}
			s := sonos.Connect(d, nil, sonos.SVC_MUSIC_SERVICES)
			currentVolume, err := s.GetAutoplayVolume()
			if err != nil {
				log.Fatalf("Error: %s\n", err)
			}
			fmt.Printf("Current volume: %d\n", currentVolume)
		}
	}
}
