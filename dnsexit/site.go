package dnsexit

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Configures a DNSExit site

type site struct {
	domains  string
	key      string
	address  string
	interval int
}

func (s *site) GetDomains() ([]string, error) {
	var envVarSet bool

	if s.domains == "" {
		s.domains, envVarSet = os.LookupEnv("DOMAINS")
		if !envVarSet {
			log.Error("Missing DNSExit domain name(s).")
			return nil, errors.New("domain name(s) not found")
		}
	}

	domains := strings.Split(s.domains, ",")

	return domains, nil
}

func (s *site) GetAPIKey() (string, error) {
	var envVarSet bool

	if s.key == "" {
		s.key, envVarSet = os.LookupEnv("API_KEY")
		if !envVarSet {
			log.Error("Missing DNSExit API Key.")
			return "", errors.New("API key not found")
		}
	}

	return s.key, nil
}

func (s *site) GetIPAddr() (string, error) {
	var envVarSet bool

	if s.address == "" {
		s.address, envVarSet = os.LookupEnv("IP_ADDR")
		if !envVarSet {
			log.Info("Using network egress IP address for update record.")
			type responseData struct {
				IP string `json:"ip"`
			}

			data := responseData{}

			url := "https://ifconfig.co"

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				log.Error("Failed to create egress IP address request.")
				return "", errors.New("failed to create egress IP address HTTP request")
			}

			req.Header.Set("Accept", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Error("Egress HTTP request failed.")
				return "", errors.New("HTTP request for egress IP failed")
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Error("Egress IP address not found.")
				return "", errors.New("egress IP response body failed")
			}
			defer resp.Body.Close()

			err = json.Unmarshal(body, &data)
			if err != nil {
				log.Error("Failed to parse egress IP address.")
				return "", errors.New("egress IP info not found in response body")
			}

			return data.IP, nil
		}
	}

	// test for valid IP address provided by user
	if net.ParseIP(s.address) == nil {
		log.Error("Invalid IP address provided: " + s.address + ".")
		return "", errors.New(s.address + " is an invalid IP address")
	}

	return s.address, nil
}

func (s *site) GetInterval() int {
	interval, envVarSet := os.LookupEnv("CHECK_INTERVAL")

	if s.interval != defaultInterval {
		return s.interval
	} else if envVarSet {
		// valid integer check, failure sets i to 0
		i, _ := strconv.Atoi(interval)
		if i != 0 {
			return i
		} else {
			log.Info("Invalid interval value was provided, defaulting to " + strconv.Itoa(defaultInterval) + " minutes.")
		}
	}

	return defaultInterval
}
