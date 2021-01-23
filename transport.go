package euler

import (
	"net"
	"time"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

type DNSTransport struct {
	nameserver string
	net        string // ["TCP"|"UDP"]
}

func NewDNSTransport(nameserver, net string) *DNSTransport {
	return &DNSTransport{
		nameserver: nameserver,
		net:        net,
	}
}

func (t *DNSTransport) Query(domain string, qtype uint16, edns0ip ...net.IP) (*dns.Msg, error) {
	req := &dns.Msg{}
	req.SetQuestion(dns.Fqdn(domain), qtype)
	if edns0ip != nil {
		EDNS0(req, edns0ip[0])
	}
	return t.Exchange(req)
}

func (t *DNSTransport) Exchange(req *dns.Msg) (*dns.Msg, error) {
	conn, err := net.DialTimeout(t.net, t.nameserver, queryTimeout)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer conn.Close()
	c := new(dns.Conn)
	c.Conn = conn
	opt := req.IsEdns0()
	if opt != nil && opt.UDPSize() >= dns.MinMsgSize {
		c.UDPSize = opt.UDPSize()
	}
	c.SetWriteDeadline(time.Now().Add(queryTimeout))
	if err = c.WriteMsg(req); err != nil {
		return nil, errors.WithStack(err)
	}
	c.SetReadDeadline(time.Now().Add(queryTimeout))
	msg, err := c.ReadMsg()
	if err == nil && msg.Id != req.Id {
		err = dns.ErrId
	}
	return msg, errors.WithStack(err)
}
