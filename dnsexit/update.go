package dnsexit

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type IPAddrAPI interface {
	getContent() string
	egressIP() (string, error)
}

type updateRecord struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

func (u updateRecord) getContent() string {
	return u.Content
}

func (u updateRecord) egressIP() (string, error) {
	log.WithFields(updateRecordLogFields).Info("Using network egress IP address for update status.")
	type responseData struct {
		IP string `json:"ip"`
	}

	data := responseData{}

	url := "https://ifconfig.co"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.WithFields(updateRecordLogFields).Error("Failed to create IP address request client.")

		return "", err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithFields(updateRecordLogFields).Error("Failed to determine egress IP address.")

		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(updateRecordLogFields).Error("Failed to read site IP address response body.")

		return "", err
	}
	defer resp.Body.Close()

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error(err)
		log.WithFields(updateRecordLogFields).Errorln("Failed to retrieve site IP address.")

		return "", err
	}

	return data.IP, nil
}

func getUpdateIP(ip IPAddrAPI) (string, error) {
	flagIP := ip.getContent()

	if flagIP == "" {
		address, varSet := os.LookupEnv("IP_ADDR")
		if varSet {
			log.WithFields(updateRecordLogFields).Info("Using environment variable IP_ADDR value for update status.")

			return address, nil
		}

		address, err := ip.egressIP()
		if err != nil {
			return "", err
		}

		return address, nil
	}

	return flagIP, nil
}
