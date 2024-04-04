package dnsexit

import (
	"flag"
	"log/slog"
	"os"
	"time"
)

const (
	apiURL          string = "https://api.dnsexit.com/dns/"
	recordType      string = "A"
	recordTTL       int    = 480
	defaultInterval int    = 10
)

var log = getLogger()

func getLogger() *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return logger
}

func CLI() {
	// read in CLI parameters
	cliDomains := flag.String("domains", "", "DNSExit domain name(s)")
	cliKey := flag.String("key", "", "DNSExit API key")
	cliIPAddr := flag.String("ip", "", "Desired A record IP address")
	cliInterval := flag.Int("interval", 10, "Time interval in minutes")

	flag.Parse()

	s := site{
		domains:  *cliDomains,
		key:      *cliKey,
		interval: *cliInterval,
		address:  *cliIPAddr,
	}

	// set domain name(s)
	domains, err := s.GetDomains()
	if err != nil {
		os.Exit(1)
	}

	// set API key
	apiKey, err := s.GetAPIKey()
	if err != nil {
		os.Exit(1)
	}

	// set IP address
	ipAddr, err := s.GetIPAddr()
	if err != nil {
		os.Exit(1)
	}

	// set update interval
	interval := s.GetInterval()

	// create a dynamic update client for every domain provided
	clients := make([]client, 0)

	for _, d := range domains {
		updateRecordData := updateRecord{
			Type:    recordType,
			TTL:     recordTTL,
			Name:    d,
			Content: ipAddr,
		}

		updateData := update{
			Update: updateRecordData,
		}

		client := client{
			url:      apiURL,
			apiKey:   apiKey,
			record:   updateData,
			interval: interval,
		}

		clients = append(clients, client)
	}

	// run each client continuously in seperate goroutines
	channel := make(chan client)

	for _, client := range clients {
		go keepCurrent(client, channel)
	}

	for i := range channel {
		go func(c client) {
			time.Sleep(time.Duration(c.interval) * time.Minute)
			keepCurrent(c, channel)
		}(i)
	}
}
