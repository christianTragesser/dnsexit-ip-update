package dnsexit

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func CLIArgs() {
	cliDomain := flag.String("domain", "", "DNSExit domain name")
	cliKey := flag.String("key", "", "DNSExit API key")
	cliInterval := flag.Int("interval", 10, "Time interval in minutes")
	cliIPAddr := flag.String("ip", "", "Desired A record IP address")

	flag.Parse()

	updateData := updateRecord{
		Type:    recordType,
		Name:    *cliDomain,
		Content: *cliIPAddr,
		TTL:     recordTTL,
	}

	cliEvent := Event{
		URL:      apiURL,
		APIKey:   *cliKey,
		Record:   updateData,
		Interval: *cliInterval,
	}

	CLIWorkflow(cliEvent)

}

func CLIWorkflow(cliEvent Event) {
	//check for env vars and compare with flag values
	if cliEvent.Record.Name == "" {
		cliEvent.Record.Name = os.Getenv("DOMAIN")
	}

	if cliEvent.APIKey == "" {
		cliEvent.APIKey = os.Getenv("API_KEY")
	}

	if cliEvent.Record.Content == "" {
		cliEvent.Record.Content = os.Getenv("IP_ADDR")
	}

	if cliEvent.Interval == 10 {
		env, varSet := os.LookupEnv("CHECK_INTERVAL")
		if varSet {
			cliEvent.Interval, _ = strconv.Atoi(env)
		}
	}

	if hasDepencies(cliEvent) {
		domains := strings.Split(cliEvent.Record.Name, ",")
		wg := new(sync.WaitGroup)

		wg.Add(len(domains))

		for _, d := range domains {
			instance := cliEvent
			instance.Record.Name = d

			go setUpdate(wg, instance)
		}

		wg.Wait()

	} else {
		os.Exit(0)
	}

	time.Sleep(time.Duration(cliEvent.Interval) * time.Minute)

	CLIWorkflow(cliEvent)
}
