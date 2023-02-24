package dnsexit

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

type UpdateRecord struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

type updateEvent struct {
	URL      string
	APIKey   string
	Record   UpdateRecord
	Interval int
}

type UpdateResponse struct {
	Code    int      `json:"code"`
	Details []string `json:"details"`
	Message string   `json:"message"`
}

func (u *UpdateRecord) getLocationIP() string {
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

func (c *UpdateRecord) getCurrentARecord(domain string) []string {
	ips, err := getARecord(domain)
	if err != nil {
		recordLogFields["domain"] = domain

		log.WithFields(recordLogFields).Error(err)

		return []string{}
	}

	return ips
}

func (r *updateEvent) postUpdate(event updateEvent) (UpdateResponse, error) {
	var response UpdateResponse

	updatePayload := map[string]UpdateRecord{"update": event.Record}
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

		return response, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		log.Errorln("Failed to read API response body.")
	}

	err = json.Unmarshal(body, &response)

	if response.Code != 0 {
		updateLogFields["status"] = response.Code

		log.WithFields(updateLogFields).Error(response.Message)
	}

	return response, err
}

func update(wg *sync.WaitGroup, event *updateEvent) {
	defer wg.Done()

	if !recordIsCurrent(event) {
		updateLogFields["domain"] = event.Record.Name
		updateLogFields["A record"] = event.Record.Content

		if event.Record.Content == "" {
			event.Record.Content = event.Record.getLocationIP()
		}

		response, err := event.postUpdate(*event)
		if err != nil {
			log.WithFields(updateLogFields).Error("Failed to update A record.")
		}

		if response.Code == 0 && response.Message != "" {
			log.WithFields(updateLogFields).Infoln(response.Message)
		}
	}
}
