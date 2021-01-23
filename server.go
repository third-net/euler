package euler

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

var pluginChains []Plugin
var cache = NewCache(defaultCacheSize)

type replyer struct {
	resp dns.ResponseWriter
}

func (r replyer) Reply(req *dns.Msg, answer ...dns.RR) error {
	msg := new(dns.Msg)

	msg.Id = req.Id
	msg.Response = true
	msg.Opcode = req.Opcode
	msg.Rcode = dns.RcodeSuccess
	if len(req.Question) > 0 {
		msg.Question = make([]dns.Question, 1)
		msg.Question[0] = req.Question[0]
	}

	msg.RecursionAvailable = true
	msg.Answer = answer
	return r.resp.WriteMsg(msg)
}

//ListenAndServe listen & up server
func ListenAndServe(addr string, plugins ...Plugin) error {
	mux := dns.NewServeMux()
	mux.HandleFunc(".", handler)
	pluginChains = plugins
	//append default plugin
	pluginChains = append(pluginChains, BackoffPlugin{transport: NewDNSTransport("8.8.8.8:53", "tcp")})

	errChan := make(chan error)
	srv := &dns.Server{Addr: addr, Net: "udp", Handler: mux}
	go func(srv *dns.Server) {
		errChan <- srv.ListenAndServe()
	}(srv)

	return <-errChan
}

func handler(resp dns.ResponseWriter, req *dns.Msg) {
	r := replyer{
		resp: resp,
	}
	//Question holds a DNS question. Usually there is just one.
	ques := req.Question[0]
	qname := ques.Name
	//only cache A record
	cacheable := ques.Qtype == dns.TypeA
	fmt.Printf("QUERY name:%s type:%d cacheable:%t \r\n", qname, ques.Qtype, cacheable)
	//dhcp query?
	if strings.HasSuffix(qname, `.DHCP\ HOST.`) {
		//error handler
		_ = r.Reply(req)
		return
	}

	domain := qname[:len(qname)-1]
	//cache hit
	if cacheable {
		if answer, ok := cache.GET(domain); ok {
			_ = r.Reply(req, answer...)
			return
		}
	}

	//apply plugin
	for _, plugin := range pluginChains {
		answer := plugin.Handler(req)
		if answer != nil {
			_ = r.Reply(req, answer...)
			if cacheable {
				cache.ADD(domain, answer...)
			}
			return
		}
	}
	_ = r.Reply(req)
}
