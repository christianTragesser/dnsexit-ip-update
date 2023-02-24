package dnsexit

import (
	"flag"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	apiURL     string = "https://api.dnsexit.com/dns/"
	recordType string = "A"
	recordTTL  int    = 480
)

func hasDependencies(event *updateEvent) bool {
	const eventReady = true

	if event.APIKey == "" {
		log.Fatal("Missing DNSExit API Key.")
	}

	if event.Record.Name == "" {
		log.Fatal("Missing DNSExit domain name.")
	}

	if event.Record.Content != "" && net.ParseIP(event.Record.Content) == nil {
		log.Fatal("Invalid value for IP address provided.")
	}

	return eventReady
}

func CLIArgs() {
	// set dynamic dns options
	cliDomain := flag.String("domain", "", "DNSExit domain name")
	cliKey := flag.String("key", "", "DNSExit API key")
	cliInterval := flag.Int("interval", 10, "Time interval in minutes")
	cliIPAddr := flag.String("ip", "", "Desired A record IP address")

	flag.Parse()

	updateData := UpdateRecord{
		Type:    recordType,
		Name:    *cliDomain,
		Content: *cliIPAddr,
		TTL:     recordTTL,
	}

	cliUpdateEvent := updateEvent{
		URL:      apiURL,
		APIKey:   *cliKey,
		Record:   updateData,
		Interval: *cliInterval,
	}

	CLIWorkflow(&cliUpdateEvent)

}

func CLIWorkflow(event *updateEvent) {
	//check for env vars and compare with flag values
	if event.Record.Name == "" {
		event.Record.Name = os.Getenv("DOMAIN")
	}

	if event.APIKey == "" {
		event.APIKey = os.Getenv("API_KEY")
	}

	if event.Record.Content == "" {
		event.Record.Content = os.Getenv("IP_ADDR")
	}

	if event.Interval == 10 {
		env, varSet := os.LookupEnv("CHECK_INTERVAL")
		if varSet {
			event.Interval, _ = strconv.Atoi(env)
		}
	}

	if hasDependencies(event) {
		// set domain list
		domains := strings.Split(event.Record.Name, ",")

		cliLogFields["domains"] = domains
		log.WithFields(cliLogFields).Info("Checking Dynamic DNS status.")

		// set IP address value for new A record
		recordIP := setRecordIP(event)

		wg := new(sync.WaitGroup)

		wg.Add(len(domains))

		// dynamic dns update for each domain provided
		for _, d := range domains {
			instance := event
			instance.Record.Content = recordIP
			instance.Record.Name = d

			// idempotent A record updates
			go update(wg, instance)
		}

		wg.Wait()
	}

	if event.Interval > 0 {
		time.Sleep(time.Duration(event.Interval) * time.Minute)

		CLIWorkflow(event)
	}
}
