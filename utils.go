package euler

import (
	"net"
	"net/http"

	"github.com/miekg/dns"
)

//DOH dns over https //TODO
func DOH(req *dns.Msg, rt http.RoundTripper) (*dns.Msg, error) {
	// qtype := req.Question[0].Qtype
	// qname := req.Question[0].Name
	// var subnet net.IP
	// opt := req.IsEdns0()
	// if opt != nil {
	// 	for _, o := range opt.Option {
	// 		if _subnet, ok := o.(*dns.EDNS0_SUBNET); ok {
	// 			subnet = _subnet.Address
	// 		}
	// 	}
	// }
	return nil, nil
}

//EDNS0 apply edns0  this need more test
func EDNS0(m *dns.Msg, addr net.IP) {
	if addr == nil {
		return
	}
	option := m.IsEdns0()
	if option == nil {
		option = new(dns.OPT)
		option.Hdr.Name = "."
		option.Hdr.Rrtype = dns.TypeOPT
		m.Extra = append(m.Extra, option)
	}
	var subnet *dns.EDNS0_SUBNET
	for _, o := range option.Option {
		if _subnet, ok := o.(*dns.EDNS0_SUBNET); ok {
			subnet = _subnet
			break
		}
	}
	if subnet == nil {
		subnet = new(dns.EDNS0_SUBNET)
		option.Option = append(option.Option, subnet)
	}
	subnet.Code = dns.EDNS0SUBNET
	subnet.Address = addr
	if addr.To4() != nil {
		subnet.Family = 1
		subnet.SourceNetmask = 32
	} else {
		subnet.Family = 2
		subnet.SourceNetmask = 128
	}
	subnet.SourceScope = 0
}

//DNSMsgToAnswer just for AAAA?
func DNSMsgToAnswer(msg *dns.Msg) (dns.RR, net.IP) {
	if msg == nil {
		return nil, nil
	}
	for _, answer := range msg.Answer {
		switch v := answer.(type) {
		case *dns.A:
			if v != nil && len(v.A) != 0 {
				return v, v.A
			}
		case *dns.AAAA:
			if v != nil && len(v.AAAA) != 0 {
				return v, v.AAAA
			}
		}
	}
	return nil, nil
}
