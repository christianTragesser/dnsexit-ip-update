package dnsexit

import (
	"context"
	"math/rand"
	"net"
	"time"
)

var nameservers = [3]string{"198.204.241.154", "204.27.62.66", "69.197.184.202"}

func setRecordIP(event *updateEvent) string {
	// check for IP flag
	// use current egress IP if no IP flag provided
	if event.Record.Content == "" {
		event.Record.Content = event.Record.getLocationIP()
		recordLogFields["IP"] = event.Record.Content
		log.WithFields(recordLogFields).Info("Using location determined IP address.")
	} else {
		recordLogFields["IP"] = event.Record.Content
		log.WithFields(recordLogFields).Info("IP address argument provided.")
	}

	return event.Record.Content
}

func getARecord(domain string) ([]string, error) {
	// take in domain name
	// use dnsexit nameservers to resolve current A record IP address
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

func recordIsCurrent(event *updateEvent) bool {
	recordLogFields["domain"] = event.Record.Name

	currentRecords := event.Record.getCurrentARecord(event.Record.Name)

	if len(currentRecords) > 0 {
		for _, record := range currentRecords {
			if event.Record.Content == record {
				log.Infof("A record for %s domain is up to date.", event.Record.Name)

				return true
			}
		}

		return false
	}

	return false
}
