package dnsexit

import (
	"errors"
	"net"
	"os"
	"strconv"
)

type CLICommand struct {
	domain   string
	key      string
	interval int
	address  string
}

func (cmd *CLICommand) setUpdateData() (updateRecord, error) {
	record := updateRecord{
		Type:    recordType,
		TTL:     recordTTL,
		Name:    cmd.domain,
		Content: cmd.address,
	}

	if record.Name == "" {
		name, varSet := os.LookupEnv("DOMAIN")
		if !varSet {
			return record, errors.New("Missing DNSExit domain name(s).")
		}

		record.Name = name
	}

	if record.Content == "" {
		address, err := getUpdateIP(record)
		if err != nil {
			log.WithFields(cliLogFields).Error(err)

			return record, err
		}

		record.Content = address
	} else {
		log.WithFields(cliLogFields).Info("Using IP flag value for update status.")
	}

	// test for valid IP address
	if net.ParseIP(record.Content) == nil {
		return record, errors.New("Invalid IP address provided to client.")
	}

	return record, nil
}

func (cmd *CLICommand) setClient(updateData updateRecord) (client, error) {
	client := client{
		URL:      apiURL,
		Record:   updateData,
		APIKey:   cmd.key,
		Interval: cmd.interval,
	}

	if client.APIKey == "" {
		key, varSet := os.LookupEnv("API_KEY")
		if !varSet {
			return client, errors.New("Missing DNSExit API Key.")
		}

		client.APIKey = key
	}

	if client.Interval == defaultInterval {
		interval, varSet := os.LookupEnv("CHECK_INTERVAL")
		if varSet {
			client.Interval, _ = strconv.Atoi(interval)
		}
	}

	return client, nil
}
