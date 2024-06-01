package networking

import (
	"log"
	"net"

	"github.com/miekg/dns"
)


func StartDNSServer() {
	dns.HandleFunc("cluster.local.", handleDNSQuery)
	server := &dns.Server{Addr: ":53", Net: "udp"}
	log.Printf("Starting DNS server on %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start DNS server: %v", err)
	}
}

func handleDNSQuery(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := r.Question[0].Name
		ip, exists := serviceIPMap[domain]
		if exists {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(ip),
			})
		}
	}
	w.WriteMsg(&msg)
}

var serviceIPMap = map[string]string{
	"redis-service.cluster.local.": "10.0.0.10",
}

func registerService(name, ip string) {
    serviceIPMap[name] = ip
}

func deregisterService(name string) {
    delete(serviceIPMap, name)
}