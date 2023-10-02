package dnsexit

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"sync"
)

type DNSExitResponse struct {
	Code    int      `json:"code"`
	Details []string `json:"details"`
	Message string   `json:"message"`
}

type clientAPI interface {
	currentRecords() ([]string, error)
	getDomain() string
}

type client struct {
	URL      string
	APIKey   string
	Record   updateRecord
	Interval int
}

func (c client) getDomain() string {
	return c.Record.Name
}

func (c client) setUpdateIP() (string, error) {
	var err error

	if c.Record.Content == "" {
		c.Record.Content, err = getUpdateIP(c.Record)
		if err != nil {
			log.WithFields(clientLogFields).Error(err)

			return c.Record.Content, err
		}
	} else {
		log.WithFields(cliLogFields).Info("Using IP flag value for update status.")
	}

	// test for valid IP address
	if net.ParseIP(c.Record.Content) == nil {
		return c.Record.Content, errors.New("Invalid IP address provided to client.")
	}

	return c.Record.Content, err
}

func (c client) currentRecords() ([]string, error) {
	ips, err := resolve(c)
	if err != nil {
		clientLogFields["domain"] = c.Record.Name
		log.WithFields(clientLogFields).Error(err)
	}

	return ips, err
}

func (c client) current(currentRecords []string, address string) bool {
	if len(currentRecords) > 0 {
		for _, ip := range currentRecords {
			if ip == address {
				return true
			}
		}
	}

	return false
}

func (c client) postUpdate() (DNSExitResponse, error) {
	var response DNSExitResponse

	updatePayload := map[string]updateRecord{"update": c.Record}
	jsonPayload, _ := json.Marshal(updatePayload)
	data := bytes.NewReader([]byte(jsonPayload))

	req, err := http.NewRequest("POST", c.URL, data)
	if err != nil {
		log.WithFields(clientLogFields).Error("Failed to create HTTP POST to DNSExit API.")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", c.APIKey)
	req.Header.Set("domain", c.Record.Name)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		log.WithFields(clientLogFields).Error("API POST failed for dynamic update.")

		return response, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		log.WithFields(clientLogFields).Errorln("Failed to read API response body.")
	}

	err = json.Unmarshal(body, &response)

	if response.Code != 0 {
		clientLogFields["status"] = response.Code
		log.WithFields(clientLogFields).Error(response.Message)

		return response, err
	}

	return response, err
}

func (c client) update(wg *sync.WaitGroup) {
	defer wg.Done()

	clientLogFields["domain"] = c.Record.Name
	log.WithFields(clientLogFields).Info("Checking Dynamic DNS status.")

	currentIPs, err := c.currentRecords()
	if err != nil {
		log.Error("Unable to resolve the provided domain name.")
	}

	if c.current(currentIPs, c.Record.Content) {
		clientLogFields["domain"] = c.Record.Name
		clientLogFields["IP"] = c.Record.Content
		clientLogFields["type"] = c.Record.Type
		log.WithFields(clientLogFields).Info("Dynamic DNS record is up to date.")
	} else {
		response, err := c.postUpdate()
		if err != nil {
			log.WithFields(clientLogFields).Error("Dynamic DNS update failed.")
		}

		if response.Code == 0 && response.Message != "" {
			clientLogFields["domain"] = c.Record.Name
			clientLogFields["IP"] = c.Record.Content
			clientLogFields["type"] = c.Record.Type
			log.WithFields(clientLogFields).Infoln("Dynamic DNS successfully updated.")
		}
	}
}
