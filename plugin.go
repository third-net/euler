package euler

import "github.com/miekg/dns"

//Plugin plugin
type Plugin interface {
	Handler(req *dns.Msg) []dns.RR
}
