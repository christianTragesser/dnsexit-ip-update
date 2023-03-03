package dnsexit

import (
	"flag"
	"fmt"
)

const (
	apiURL          string = "https://api.dnsexit.com/dns/"
	recordType      string = "A"
	recordTTL       int    = 480
	defaultInterval int    = 10
)

type Client struct {
	URL      string
	APIKey   string
	Record   UpdateRecord
	Interval int
}

func CLI() {
	// set dynamic dns client options
	cliDomain := flag.String("domain", "", "DNSExit domain name")
	cliKey := flag.String("key", "", "DNSExit API key")
	cliInterval := flag.Int("interval", defaultInterval, "Time interval in minutes")
	cliIPAddr := flag.String("ip", "", "Desired A record IP address")

	flag.Parse()

	cmd := CLICommand{
		domain:   *cliDomain,
		key:      *cliKey,
		interval: *cliInterval,
		address:  *cliIPAddr,
	}

	recordData, err := cmd.setUpdateData()
	if err != nil {
		log.Fatal(err)
	}

	client, err := cmd.setClient(recordData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(client)
}
