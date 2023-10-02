package dnsexit

import (
	"context"
	"math/rand"
	"net"
	"time"
)

func resolve(c clientAPI) ([]string, error) {
	var ip []string
	var nsErr error
	domain := c.getDomain()

	// get DNS Exit nameservers
	nameServers, _ := net.LookupNS(c.getDomain())

	// randomize nameservers list
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(nameServers), func(i, j int) { nameServers[i], nameServers[j] = nameServers[j], nameServers[i] })

	for _, ns := range nameServers {
		nsIP, err := net.LookupIP(ns.Host)
		if err != nil {
			resolverLogFields["nameserver"] = ns.Host
			log.WithFields(resolverLogFields).Warn("Failure to resolve Nameserver IP address.")

			continue
		}

		r := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: 5000 * time.Millisecond,
				}
				return d.DialContext(ctx, "tcp", ns.Host+":53")
			},
		}

		ip, nsErr = r.LookupHost(context.Background(), domain)
		if nsErr != nil {
			resolverLogFields["nameserver"] = ns.Host + "(" + nsIP[0].String() + ")"
			log.WithFields(resolverLogFields).Warn("DNSExit Nameserver failure.")

			continue
		}

		resolverLogFields["nameserver"] = ns.Host
	}

	resolverLogFields["current"] = ip
	resolverLogFields["domain"] = domain
	log.WithFields(resolverLogFields).Info("DNS resolution.")

	return ip, nsErr
}
