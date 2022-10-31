package main

// Docs: https://svrooij.io/sonos-api-docs/sonos-communication.html

import (
	"fmt"
	"log"
	"net"
)

const ListenPort = 1900

func main() {
	addr := net.UDPAddr{
		Port: ListenPort,
		IP:   net.ParseIP("0.0.0.0"),
	}
	listenConn, err := net.ListenUDP("udp", &addr) // code does not block here
	if err != nil {
		log.Fatalf("Error listening on port %d: %s", ListenPort, err)
	}
	defer listenConn.Close()
	go func() {
		var buf [1024]byte
		for {
			rlen, remote, err := listenConn.ReadFromUDP(buf[:])
			if err != nil {
				log.Fatalf("error during read: %s", err)
			}
			log.Printf("read %d bytes from %v: %s\n", rlen, remote, buf[:rlen])
		}
	}()

	sendConn, err := net.Dial("tcp", "239.255.255.250:1900")
	if err != nil {
		log.Fatalf("error dialing send: %s", err)
	}
	var msg = "M-SEARCH * HTTP/1.1\n"
	var headers = []struct{ key, value string }{
		{"HOST", fmt.Sprintf("239.255.255.250:%d", ListenPort)},
		{"MAN", "ssdp:discover"},
		{"MX", "1"},
		{"ST", "urn:schemas-upnp-org:device:ZonePlayer:1"},
	}
	for _, header := range headers {
		msg += fmt.Sprintf("%s: %s\n", header.key, header.value)
	}
	fmt.Fprintf(sendConn, msg)
	log.Printf("sent %d bytes ssdp\n", len(msg))
	//log.Printf("%s\n", msg)
	<-make(chan struct{})
}
