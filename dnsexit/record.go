package dnsexit

import (
	"encoding/json"
	"io"
	"net/http"
)

type recordStatusAPI interface {
	getRecords(domain string) []string
	getLocationIP() string
}

type recordStatus struct{}

func (c recordStatus) getRecords(domain string) []string {
	ips, err := dnsLookup(domain)
	if err != nil {
		recordLogFields["domain"] = domain

		log.WithFields(recordLogFields).Error(err)

		return []string{}
	}

	return ips
}

func (d recordStatus) getLocationIP() string {
	type responseData struct {
		IP string `json:"ip"`
	}

	data := responseData{}

	url := "https://ifconfig.co"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Errorln("Failed create HTTP request client.")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err)
		log.WithFields(recordLogFields).Error("Failed to determine location IP address.")

		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		log.Errorln("Failed to read site IP address response body.")
	}
	defer resp.Body.Close()

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error(err)
		log.Errorln("Failed to retrieve site IP address.")
	}

	return data.IP
}

func setRecordIP(api recordStatusAPI, event Event) string {
	if event.Record.Content == "" {
		event.Record.Content = api.getLocationIP()
		recordLogFields["IP"] = event.Record.Content
		log.WithFields(recordLogFields).Info("Using location determined IP address.")
	} else {
		recordLogFields["IP"] = event.Record.Content
		log.WithFields(recordLogFields).Info("IP address argument provided.")
	}

	return event.Record.Content
}

func recordIsCurrent(api recordStatusAPI, event Event) bool {
	recordLogFields["domain"] = event.Record.Name

	currentRecords := api.getRecords(event.Record.Name)

	if len(currentRecords) > 0 {
		for _, record := range currentRecords {
			if event.Record.Content == record {
				log.Infof("A record for %s domain is up to date.", event.Record.Name)

				return true
			}
		}

		return false
	}

	return false
}
