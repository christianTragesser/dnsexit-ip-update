package dnsexit

import (
	"flag"
	"os"
	"time"
)

var logc = GetLogger("cli")

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
			response, err = dynamicUpdate(response, cliEvent)
			if err != nil {
				logc.Errorln("Dynamic IP update failed.")
			}

			if response.Message != "" {
				logc.Infoln(response)
			}
		}
	} else {
		os.Exit(0)
	}

	time.Sleep(time.Duration(cliEvent.Interval) * time.Minute)

	CLIWorkflow(cliEvent)
}
