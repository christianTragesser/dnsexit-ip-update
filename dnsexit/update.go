package dnsexit

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
)

const (
	apiURL     string = "https://api.dnsexit.com/dns/"
	recordType string = "A"
	recordTTL  int    = 480
)

type updateRecord struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

type Event struct {
	Code     int      `json:"code"`
	Details  []string `json:"details"`
	Message  string   `json:"message"`
	URL      string
	APIKey   string
	Record   updateRecord
	Interval int
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
		log.Errorln("Failed to create HTTP POST to DNSExit API.")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", event.APIKey)
	req.Header.Set("domain", event.Record.Name)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		log.WithFields(updateLogFields).Error("API POST failed for dynamic update.")

		return responseData, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		log.Errorln("Failed to read API response body.")
	}

	err = json.Unmarshal(body, &responseData)

	if responseData.Code != 0 {
		updateLogFields["status"] = responseData.Code

		log.WithFields(updateLogFields).Error(responseData.Message)
	}

	return responseData, err
}

func dynamicUpdate(api dnsExitAPI, event Event) (Event, error) {
	eventResponse, err := api.setUpdate(event)
	if err != nil {
		log.WithFields(updateLogFields).Error("Failed to set A record update.")
	}

	return eventResponse, err
}

func hasDepencies(event Event) bool {
	var eventReady = true

	if event.APIKey == "" {
		log.Errorln("Missing API Key.")
		eventReady = false
	}

	if event.Record.Name == "" {
		log.Errorln("Missing DNSExit domain name.")
		eventReady = false
	}

	if event.Record.Content != "" && net.ParseIP(event.Record.Content) == nil {
		log.Errorf("Invalid A record content provided: %s", event.Record.Content)
		eventReady = false
	}

	return eventReady
}
