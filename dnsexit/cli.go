package dnsexit

import (
	"flag"
	"os"
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
