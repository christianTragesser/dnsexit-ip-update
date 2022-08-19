package dnsexit

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	apiURL     string = "https://api.dnsexit.com/dns/"
	recordType string = "A"
	recordTTL  int    = 480
)

var logd = GetLogger("dynamic")

type updateRecord struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

type Event struct {
	Code    int      `json:"code"`
	Details []string `json:"details"`
	Message string   `json:"message"`
	URL     string
	APIKey  string
	Record  updateRecord
}

type dnsExitAPI interface {
	setUpdate(event Event) (Event, error)
}

func (r Event) setUpdate(event Event) (Event, error) {
	var responseData Event

	updatePayload := map[string]updateRecord{"update": event.Record}
	jsonPayload, _ := json.Marshal(updatePayload)
	data := bytes.NewReader([]byte(jsonPayload))

	req, err := http.NewRequest("POST", event.URL, data)
	if err != nil {
		logd.Errorln("Failed to create HTTP request.")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", event.APIKey)
	req.Header.Set("domain", event.Record.Name)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logd.Error(err)
		logd.Errorln("HTTP POST failed for dynamic update.")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logd.Error(err)
		logd.Errorln("Failed to read API response body.")
	}

	err = json.Unmarshal(body, &responseData)

	return responseData, err
}

func dynamicUpdate(api dnsExitAPI, event Event) (Event, error) {
	eventResponse, err := api.setUpdate(event)
	if err != nil {
		logd.Errorln("Failed to set A record update.")
	}

	return eventResponse, err
}

func dynamicUpdateDepencies(event Event) bool {
	var eventReady = true

	if event.APIKey == "" {
		logd.Errorln("Missing API Key.")
		eventReady = false
	}
	if event.Record.Name == "" {
		logd.Errorln("Missing DNSExit domain name.")
		eventReady = false
	}
	if event.Record.Content == "" {
		logd.Errorln("Missing IP address.")
		eventReady = false
	}

	return eventReady
}
