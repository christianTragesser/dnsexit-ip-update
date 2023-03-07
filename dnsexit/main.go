package dnsexit

import (
	"flag"
	"strings"
	"sync"
	"time"
)

const (
	apiURL          string = "https://api.dnsexit.com/dns/"
	recordType      string = "A"
	recordTTL       int    = 480
	defaultInterval int    = 10
)

func CLI() {
	// read in CLI parameters
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

	// construct DNSExit dynamic update record
	updateRecordData, err := cmd.setUpdateDomains()
	if err != nil {
		log.Fatal(err)
	}

	// create an update client for every domain provided in CLI command
	domains := strings.Split(updateRecordData.Name, ",")
	clients := make([]client, len(domains))

	for i, d := range domains {
		update := updateRecordData
		update.Name = d

		client, err := cmd.setClient(update)
		if err != nil {
			log.Fatal(err)
		}

		clients[i] = client
	}

	// run clients in persistent loop
	CLIUpdate(clients)
}

func CLIUpdate(clients []client) {
	if len(clients) > 0 {
		var interval int

		wg := new(sync.WaitGroup)
		wg.Add(len(clients))

		for _, c := range clients {
			var err error

			interval = c.Interval

			c.Record.Content, err = c.setUpdateIP()
			if err != nil {
				log.Fatal(err)
			}

			go c.update(wg)
		}

		wg.Wait()

		if interval > 0 {
			time.Sleep(time.Duration(interval) * time.Minute)

			CLIUpdate(clients)
		}
	}
}
