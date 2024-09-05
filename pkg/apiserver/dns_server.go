package apiserver

import (
	"log"

	"github.com/miekg/dns"
)

type DNSServer struct {
	DNSHandler *DNSHandler
}

func NewDNSServer(dnsHandler *DNSHandler) *DNSServer {
	return &DNSServer{DNSHandler: dnsHandler}
}

func (s *DNSServer) StartDNSServer() {
	dns.HandleFunc("cluster.local.", s.DNSHandler.DNSQueryHandler)
	server := &dns.Server{Addr: ":53", Net: "udp"}
	log.Printf("Starting DNS server on %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start DNS server: %v", err)
	}
}
