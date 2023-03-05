package dnsexit

import (
	"context"
	"math/rand"
	"net"
	"time"
)

func resolve(c clientAPI) ([]string, error) {
	// use dnsexit nameservers to resolve current A client IP address
	var nameservers = [3]string{"198.204.241.154", "204.27.62.66", "69.197.184.202"}

	srvNum := rand.Intn(len(&nameservers))

	nameserver := nameservers[srvNum]

	domain := c.getDomain()

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

	resolverLogFields["current"] = ip
	resolverLogFields["domain"] = domain
	resolverLogFields["nameserver"] = nameserver
	log.WithFields(resolverLogFields).Info("DNS resolution.")

	return ip, err
}
