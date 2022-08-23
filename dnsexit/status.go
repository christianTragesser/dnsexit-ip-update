package dnsexit

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
)

var logs = GetLogger("status")

type recordStatusAPI interface {
	getRecords(domain string) []net.IP
	getLocationIP() string
}

type recordStatus struct{}

func (c recordStatus) getRecords(domain string) []net.IP {
	ips, err := net.LookupIP(domain)
	if err != nil {
		logs.Error(err)
		logs.Errorf("Failed to resolve hostname %s.", domain)
	}

	return ips
}

func (d recordStatus) getLocationIP() string {
	type responseData struct {
		IP string `json:"ip"`
	}

	var data responseData

	url := "https://ifconfig.co"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logs.Error(err)
		logs.Errorln("Failed create HTTP request client.")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.Error(err)
		logs.Errorln("IP address HTTP request failed.")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err)
		logs.Errorln("Failed to read site IP address response body.")
	}
	defer resp.Body.Close()

	err = json.Unmarshal(body, &data)
	if err != nil {
		logs.Error(err)
		logs.Errorln("Failed to retrieve site IP address.")
	}

	return data.IP
}

func recordIsCurrent(api recordStatusAPI, event Event) bool {

	if event.Record.Content == "" {
		event.Record.Content = api.getLocationIP()
		logs.Infof("Using IP address %s for %s A record value.", event.Record.Content, event.Record.Name)
	} else {
		logs.Infof("Using preferred IP address %s for %s A Record value.", event.Record.Content, event.Record.Name)
	}

	currentRecords := api.getRecords(event.Record.Name)

	for _, record := range currentRecords {
		if event.Record.Content == record.String() {
			logs.Infof("A record for %s domain is up to date.", event.Record.Name)
			return true
		}
	}

	return false
}
