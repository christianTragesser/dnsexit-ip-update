package dnsexit

import (
	"context"
	"math/rand"
	"net"
	"time"
)

var nameservers = [4]string{"162.244.82.74", "198.204.241.154", "204.27.62.66", "69.197.184.202"}

func dnsLookup(domain string) ([]string, error) {
	srvNum := rand.Intn(len(&nameservers))
	nameserver := nameservers[srvNum]

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, nameserver+":53")
		},
	}

	ip, err := r.LookupHost(context.Background(), domain)

	resolverLogFields["nameserver"] = nameserver
	resolverLogFields["domain"] = domain
	resolverLogFields["result"] = ip
	log.WithFields(resolverLogFields).Info("DNS resolution.")

	return ip, err
}
