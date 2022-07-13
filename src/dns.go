package main

import (
	"errors"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
)

var dns_config *dns.ClientConfig

func init() {
	var err error
	dns_config, err = dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		log.Fatal("Unable to find dns servers")
	}
}

func getIPs(hostname string) (ips []net.IP) {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.RecursionDesired = true
	ipmap := map[string]struct{}{}

	for i := 0; i < 2; i++ {
		for _, server := range dns_config.Servers {
			m.SetQuestion(strings.TrimSuffix(hostname, ".")+".", dns.TypeA)
			r, _, err := c.Exchange(m, server+":"+dns_config.Port)
			if err == nil {
				for _, a := range r.Answer {
					if mx, ok := a.(*dns.A); ok {
						str := mx.A.String()
						if _, ok := ipmap[str]; !ok {
							ips = append(ips, mx.A)
							ipmap[str] = struct{}{}
						}
					}
				}
				if r.Rcode != dns.RcodeSuccess {
					err = errors.New("A lookup error " + os.Args[2])
				}
			} else {
				log.Println("Error querying DNS server for", hostname, err)
			}

			m.SetQuestion(strings.TrimSuffix(hostname, ".")+".", dns.TypeAAAA)
			r, _, err = c.Exchange(m, server+":"+dns_config.Port)
			if err == nil {
				for _, a := range r.Answer {
					if mx, ok := a.(*dns.AAAA); ok {
						str := mx.AAAA.String()
						if _, ok := ipmap[str]; !ok {
							ips = append(ips, mx.AAAA)
							ipmap[str] = struct{}{}
						}
					}
				}
				if r.Rcode != dns.RcodeSuccess {
					err = errors.New("A lookup error " + os.Args[2])
				}
			} else {
				log.Println("Error querying DNS server for", hostname, err)
			}
		}
		time.Sleep(70 * time.Millisecond)
	}
	return ips
}
