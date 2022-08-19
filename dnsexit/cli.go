package dnsexit

import (
	"flag"
	"fmt"
	"os"
)

var logc = GetLogger("cli")

func CLIArgs() {
	cliArgs := flag.NewFlagSet("update", flag.ExitOnError)
	cliDomain := cliArgs.String("domain", "", "DNSExit domain name")
	cliIPAddr := cliArgs.String("ip", "", "IP address")
	cliKey := cliArgs.String("key", "", "DNS API key")

	err := cliArgs.Parse(os.Args[2:])
	if err != nil {
		logc.Errorln("Failed to parse update CLI args.")
	}

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
	action := os.Args[1]

	switch action {
	case "update":
		if dynamicUpdateDepencies(cliEvent) {
			response, err = dynamicUpdate(response, cliEvent)
			if err != nil {
				logc.Errorln("Dynamic IP update failed.")
			}
		} else {
			os.Exit(0)
		}
	case "-h":
		fmt.Println("dnsexit options:\n dnsexit update -h")
	default:
		fmt.Printf("Invalid argument: '%s'", action)
		fmt.Println("\n 'dnsexit -h' to list valid options")
	}

	return response, err
}
