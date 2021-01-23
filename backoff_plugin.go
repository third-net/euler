package euler

import (
	"fmt"

	"github.com/miekg/dns"
)

//BackoffPlugin backoff plugin
type BackoffPlugin struct {
	transport *DNSTransport
}

func (p BackoffPlugin) Handler(req *dns.Msg) []dns.RR {
	resp, err := p.transport.Exchange(req)
	if err == nil && resp != nil {
		return resp.Answer
	}
	fmt.Printf("backoff plugin exchange err: %v \r\n", err)
	return nil
}
