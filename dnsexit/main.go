package dnsexit

import (
	"flag"
)

const (
	apiURL          string = "https://api.dnsexit.com/dns/"
	recordType      string = "A"
	recordTTL       int    = 480
	defaultInterval int    = 10
)

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

	updateRecordData, err := cmd.setUpdateData()
	if err != nil {
		log.Fatal(err)
	}

	client, err := cmd.setClient(updateRecordData)
	if err != nil {
		log.Fatal(err)
	}

	cliLogFields["domain"] = cmd.domain
	log.WithFields(cliLogFields).Info("Checking Dynamic DNS status.")

	currentIPs, err := client.currentRecords()
	if err != nil {
		log.Fatal("Unable to resolve the provided domain name.")
	}

	if client.current(currentIPs, client.Record.Content) {
		cliLogFields["domain"] = client.Record.Name
		cliLogFields["IP"] = client.Record.Content
		cliLogFields["Type"] = client.Record.Type
		log.WithFields(cliLogFields).Info("Dynamic DNS record is up to date.")
	}
}
