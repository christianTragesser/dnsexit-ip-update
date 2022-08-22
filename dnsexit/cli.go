package dnsexit

import (
	"flag"
	"os"
)

var logc = GetLogger("cli")

func CLIArgs() {
	cliDomain := flag.String("domain", "", "DNSExit domain name")
	cliIPAddr := flag.String("ip", "", "IP address")
	cliKey := flag.String("key", "", "DNS API key")

	flag.Parse()

	updateData := updateRecord{
		Type:    recordType,
		Name:    *cliDomain,
		Content: *cliIPAddr,
		TTL:     recordTTL,
	}

	cliEvent := Event{
		URL:    apiURL,
		APIKey: *cliKey,
		Record: updateData,
	}

	eventResp, _ := CLIWorkflow(cliEvent)

	if eventResp.Message != "" {
		logc.Infoln(eventResp)
	}
}

func CLIWorkflow(cliEvent Event) (Event, error) {
	var response Event
	var err error
	statusAPI := recordStatus{}

	if !recordIsCurrent(statusAPI, cliEvent) && dynamicUpdateDepencies(cliEvent) {
		response, err = dynamicUpdate(response, cliEvent)
		if err != nil {
			logc.Errorln("Dynamic IP update failed.")
		}
	} else {
		os.Exit(0)
	}

	return response, err
}
