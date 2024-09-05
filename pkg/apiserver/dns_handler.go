package apiserver

import (
	"maden/pkg/etcd"
	"maden/pkg/shared"

	"net"

	"github.com/miekg/dns"
)

type DNSHandler struct {
	Repo etcd.DNSRepository
}

func NewDNSHandler(repo etcd.DNSRepository) *DNSHandler {
	return &DNSHandler{Repo: repo}
}

func (h *DNSHandler) DNSQueryHandler(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := r.Question[0].Name
		ip, err := h.Repo.ResolveService(domain)
		if err != nil {
			shared.Log.Errorf("Failed to resolve service: %v", err)
			w.WriteMsg(&msg)
			return
		}

		msg.Answer = append(msg.Answer, &dns.A{
			Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
			A:   net.ParseIP(ip),
		})

	}
	w.WriteMsg(&msg)
}
