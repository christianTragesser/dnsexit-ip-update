package dnsexit

import (
	"errors"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type site struct {
	domains     string
	key         string
	address     string
	recordType  string
	ttl         int
	ttlSet      bool
	interval    int
	intervalSet bool
	httpClient  HTTPClient
}

func (s *site) GetDomains() ([]string, error) {
	if s.domains == "" {
		var ok bool
		s.domains, ok = os.LookupEnv("DOMAINS")
		if !ok {
			log.Error("Missing DNSExit domain name(s).")
			return nil, errors.New("domain name(s) not found")
		}
	}

	return strings.Split(s.domains, ","), nil
}

func (s *site) GetAPIKey() (string, error) {
	if s.key == "" {
		var ok bool
		s.key, ok = os.LookupEnv("API_KEY")
		if !ok {
			log.Error("Missing DNSExit API Key.")
			return "", errors.New("API key not found")
		}
	}

	return s.key, nil
}

// GetIPAddr returns the desired IP, whether it is static (user-provided), and any error.
// A non-static IP means keepCurrent should re-discover the egress IP each cycle.
func (s *site) GetIPAddr() (string, bool, error) {
	if s.address != "" {
		if net.ParseIP(s.address) == nil {
			log.Error("Invalid IP address provided: " + s.address + ".")
			return "", false, errors.New(s.address + " is an invalid IP address")
		}
		return s.address, true, nil
	}

	if addr, ok := os.LookupEnv("IP_ADDR"); ok {
		if net.ParseIP(addr) == nil {
			log.Error("Invalid IP address in IP_ADDR: " + addr + ".")
			return "", false, errors.New(addr + " is an invalid IP address")
		}
		return addr, true, nil
	}

	hc := s.httpClient
	if hc == nil {
		hc = http.DefaultClient
	}

	ip, err := fetchEgressIP(hc)
	if err != nil {
		log.Error("Failed to discover egress IP: " + err.Error())
		return "", false, err
	}

	log.Info("Using network egress IP address (" + ip + ") for update record.")
	return ip, false, nil
}

// GetRecordType returns the DNS record type to manage. Accepts A, AAAA, and SELF.
func (s *site) GetRecordType() (string, error) {
	rt := s.recordType
	if rt == "" {
		rt, _ = os.LookupEnv("RECORD_TYPE")
	}
	if rt == "" {
		return recordTypeA, nil
	}

	switch strings.ToUpper(rt) {
	case recordTypeA, recordTypeAAAA, recordTypeSelf:
		return strings.ToUpper(rt), nil
	default:
		log.Error("Unsupported record type: " + rt + ". Supported types: A, AAAA, SELF.")
		return "", errors.New("unsupported record type: " + rt)
	}
}

// GetTTL returns the record TTL in minutes. TTL 0 is valid (used for Let's Encrypt).
func (s *site) GetTTL() int {
	if s.ttlSet {
		return s.ttl
	}

	if val, ok := os.LookupEnv("RECORD_TTL"); ok {
		i, err := strconv.Atoi(val)
		if err == nil && i >= 0 {
			return i
		}
		log.Info("Invalid RECORD_TTL value, defaulting to " + strconv.Itoa(defaultTTL) + " minutes.")
	}

	return defaultTTL
}

func (s *site) GetInterval() int {
	interval := defaultInterval

	if s.intervalSet {
		interval = s.interval
	} else if val, ok := os.LookupEnv("CHECK_INTERVAL"); ok {
		i, err := strconv.Atoi(val)
		if err == nil && i > 0 {
			interval = i
		} else {
			log.Info("Invalid CHECK_INTERVAL value, defaulting to " + strconv.Itoa(defaultInterval) + " minutes.")
		}
	}

	if interval < minInterval {
		log.Warn("Interval " + strconv.Itoa(interval) + " is below the DNSExit minimum of " + strconv.Itoa(minInterval) + " minutes, using " + strconv.Itoa(minInterval) + ".")
		return minInterval
	}

	return interval
}
