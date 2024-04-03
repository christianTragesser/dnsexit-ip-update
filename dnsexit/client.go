package dnsexit

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"
)

type updateRecord struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

type DNSExitResponse struct {
	Code    int      `json:"code"`
	Details []string `json:"details"`
	Message string   `json:"message"`
}

type client struct {
	url      string
	apiKey   string
	record   updateRecord
	interval int
}

func (c client) ResolveDomain() (string, error) {
	// retrieve DNSExit nameservers
	nameServers, _ := net.LookupNS(c.record.Name)

	// select DNSExit nameserver
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := rng.Intn(len(nameServers))
	host := nameServers[i]
	ns := host.Host[:len(host.Host)-1]

	// build custom DNS resolver to query DNSExit nameserver
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 5000 * time.Millisecond,
			}
			return d.DialContext(ctx, "tcp", ns+":53")
		},
	}

	log.Info("Using " + ns + " to resolve " + c.record.Name + ".")

	// perform DNS lookup for site domain
	resolvedAddrs, err := r.LookupHost(context.Background(), c.record.Name)
	if err != nil {
		log.Error(ns + "failed to resolve " + c.record.Name)
	}

	if len(resolvedAddrs) == 0 {
		log.Error("Failed to resolve " + c.record.Name)
		return "", err
	} else {
		log.Info("Resolved domain " + c.record.Name + " to " + resolvedAddrs[0] + ".")
		return resolvedAddrs[0], nil
	}
}

func (c client) postUpdate() {
	// create POST request, send to DNSExit API, check the response
	var response DNSExitResponse

	jsonPayload, _ := json.Marshal(c.record)
	data := bytes.NewReader([]byte(jsonPayload))

	req, err := http.NewRequest("POST", c.url, data)
	if err != nil {
		log.Error("Failed to create POST method.")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("domain", c.record.Name)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("HTTP POST for dynamic update failed.")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read API response body.")
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Error("Failed to parse API response.")
	}

	if response.Code != 0 {
		log.Error(response.Message)
	} else {
		log.Info("Successfully updated " + c.record.Name + " A record.")
	}
}

func keepCurrent(c client, p chan client) {
	currentAddr, err := c.ResolveDomain()
	if err != nil {
		log.Error(err.Error())
	}

	if currentAddr == c.record.Content {
		log.Info(c.record.Name + " site IP address is up to date.")
		p <- c
	} else {
		log.Info("Updating " + c.record.Name + " A record from " + currentAddr + " to " + c.record.Content + ".")
		c.postUpdate()
		p <- c
	}
}
