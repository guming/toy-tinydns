package main

import (
	"fmt"
	"github.com/miekg/dns"
	"net"
	"os"
	"strings"
)

func resolve(name string) net.IP {
	nameserver := net.ParseIP("114.114.114.114")
	for {
		reply := dnsQuery(name, nameserver)
		if ip := getAnswer(reply); ip != nil {
			return ip
		} else if nsIP := getGlue(reply); nsIP != nil {
			nameserver = nsIP
		} else if domain := getNS(reply); domain != "" {
			nameserver = resolve(domain)
		} else {
			panic("something wrong")
		}
	}
}

func getAnswer(reply *dns.Msg) net.IP {
	for _, record := range reply.Answer {
		if record.Header().Rrtype == dns.TypeA {
			fmt.Println("  ", record)
			return record.(*dns.A).A
		}
	}
	return nil
}

func getGlue(reply *dns.Msg) net.IP {
	for _, record := range reply.Extra {
		if record.Header().Rrtype == dns.TypeA {
			fmt.Println("  ", record)
			return record.(*dns.A).A
		}
	}
	return nil
}

func getNS(reply *dns.Msg) string {
	for _, record := range reply.Ns {
		if record.Header().Rrtype == dns.TypeNS {
			fmt.Println("  ", record)
			return record.(*dns.NS).Ns
		}
	}
	return ""
}

func dnsQuery(domainName string, server net.IP) *dns.Msg {
	fmt.Printf("dig -r @%s %s\n", server.String(), domainName)
	msg := new(dns.Msg)
	msg.SetQuestion(domainName, dns.TypeA)
	c := new(dns.Client)
	reply, _, _ := c.Exchange(msg, server.String()+":53")
	return reply
}

func main() {
	name := os.Args[1]
	if !strings.HasSuffix(name, ".") {
		name = name + "."
	}
	fmt.Println("Result:", resolve(name))
}
