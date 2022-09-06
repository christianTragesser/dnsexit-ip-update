package dnsexit

import (
	"flag"
	"os"
	"strconv"
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
	var response Event
	var err error
	statusAPI := recordStatus{}

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
		env, ok := os.LookupEnv("CHECK_INTERVAL")
		if ok {
			cliEvent.Interval, _ = strconv.Atoi(env)
		}
	}

	if hasDepencies(cliEvent) {
		if !recordIsCurrent(statusAPI, cliEvent) {
			cliLogFields["domain"] = cliEvent.Record.Name
			cliLogFields["A record"] = cliEvent.Record.Content

			response, err = dynamicUpdate(response, cliEvent)
			if err != nil {
				log.WithFields(cliLogFields).Error("Dynamic IP update failed.")
			}

			if response.Code == 0 && response.Message != "" {
				log.WithFields(cliLogFields).Infoln(response.Message)
			}
		}
	} else {
		os.Exit(0)
	}

	time.Sleep(time.Duration(cliEvent.Interval) * time.Minute)

	CLIWorkflow(cliEvent)
}
