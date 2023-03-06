package dnsexit

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

type DNSExitResponse struct {
	Code    int      `json:"code"`
	Details []string `json:"details"`
	Message string   `json:"message"`
}

type clientAPI interface {
	getDomain() string
	currentRecords() ([]string, error)
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

func (c client) currentRecords() ([]string, error) {
	ips, err := resolve(c)
	if err != nil {
		clientLogFields["domain"] = c.Record.Name
		log.WithFields(clientLogFields).Error(err)

		return []string{}, err
	}

	return ips, nil
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
		updateRecordLogFields["status"] = response.Code
		log.WithFields(clientLogFields).Error(response.Message)

		return response, err
	}

	return response, err
}

func (c client) update(wg *sync.WaitGroup) {
	defer wg.Done()

	cliLogFields["domain"] = c.Record.Name
	log.WithFields(cliLogFields).Info("Checking Dynamic DNS status.")

	currentIPs, err := c.currentRecords()
	if err != nil {
		log.Fatal("Unable to resolve the provided domain name.")
	}

	if c.current(currentIPs, c.Record.Content) {
		cliLogFields["domain"] = c.Record.Name
		cliLogFields["IP"] = c.Record.Content
		cliLogFields["Type"] = c.Record.Type
		log.WithFields(cliLogFields).Info("Dynamic DNS record is up to date.")
	} else {
		response, err := c.postUpdate()
		if err != nil {
			log.WithFields(cliLogFields).Error("Dynamic DNS update failed.")
		}

		if response.Code == 0 && response.Message != "" {
			cliLogFields["domain"] = c.Record.Name
			cliLogFields["IP"] = c.Record.Content
			cliLogFields["type"] = c.Record.Type
			cliLogFields["code"] = response.Code
			cliLogFields["status"] = response.Message
			log.WithFields(cliLogFields).Infoln("Dynamic DNS update completed.")
		}
	}
}
