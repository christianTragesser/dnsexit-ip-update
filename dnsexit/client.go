package dnsexit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

// HTTPClient is satisfied by *http.Client and allows injection in tests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// DomainResolver resolves a domain name to an IP address.
type DomainResolver interface {
	Resolve(ctx context.Context, domain string) (string, error)
}

// Supported record types.
const (
	recordTypeA    = "A"
	recordTypeAAAA = "AAAA"
	recordTypeSelf = "SELF"
)

type DNSExitResponse struct {
	Code    int      `json:"code"`
	Details []string `json:"details"`
	Message string   `json:"message"`
}

type updateRecord struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	TTL       int    `json:"ttl"`
	Overwrite bool   `json:"overwrite"`
}

type update struct {
	Update updateRecord `json:"update"`
}

type client struct {
	url      string
	apiKey   string
	record   update
	interval int
	staticIP bool
	http     HTTPClient
	resolver DomainResolver
}

// netResolver implements DomainResolver using real DNS lookups against DNSExit nameservers.
type netResolver struct{}

func (netResolver) Resolve(ctx context.Context, domain string) (string, error) {
	nameServers, err := net.LookupNS(domain)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve nameservers for %s: %w", domain, err)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	host := nameServers[rng.Intn(len(nameServers))]
	ns := host.Host[:len(host.Host)-1]

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: 5 * time.Second}
			return d.DialContext(ctx, "tcp", ns+":53")
		},
	}

	log.Info("Using " + ns + " to resolve " + domain + ".")

	addrs, err := r.LookupHost(ctx, domain)
	if err != nil {
		return "", fmt.Errorf("%s failed to resolve %s: %w", ns, domain, err)
	}
	if len(addrs) == 0 {
		return "", fmt.Errorf("no addresses returned for %s", domain)
	}

	log.Info("Resolved domain " + domain + " to " + addrs[0] + ".")
	return addrs[0], nil
}

// fetchEgressIP discovers the current public egress IP address.
func fetchEgressIP(hc HTTPClient) (string, error) {
	type responseData struct {
		IP string `json:"ip"`
	}

	req, err := http.NewRequest(http.MethodGet, "https://ifconfig.co", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create egress IP request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := hc.Do(req)
	if err != nil {
		return "", fmt.Errorf("egress IP request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read egress IP response: %w", err)
	}

	var data responseData
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("failed to parse egress IP response: %w", err)
	}

	if data.IP == "" {
		return "", fmt.Errorf("empty IP in egress response")
	}

	return data.IP, nil
}

// splitDomain returns the base registrar domain and any subdomain prefix.
//
//	"sub.example.com"      → ("example.com", "sub")
//	"deep.sub.example.com" → ("example.com", "deep.sub")
//	"example.com"          → ("example.com", "")
//
// Note: does not handle multi-part public suffixes (e.g. .co.uk).
func splitDomain(host string) (base, sub string) {
	parts := strings.Split(host, ".")
	if len(parts) <= 2 {
		return host, ""
	}
	return strings.Join(parts[len(parts)-2:], "."), strings.Join(parts[:len(parts)-2], ".")
}

func (c client) postUpdate() error {
	base, hostname := splitDomain(c.record.Update.Name)

	// Build the API payload using only the hostname/subdomain in the name field.
	// The base domain goes in the request header per the DNSExit API spec.
	payload := update{
		Update: updateRecord{
			Type:      c.record.Update.Type,
			Name:      hostname,
			Content:   c.record.Update.Content,
			TTL:       c.record.Update.TTL,
			Overwrite: c.record.Update.Overwrite,
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal update payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewReader(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create POST request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("domain", base)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP POST failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read API response: %w", err)
	}

	var response DNSExitResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse API response: %w", err)
	}

	switch response.Code {
	case 0:
		log.Info("Successfully updated " + c.record.Update.Name + " " + c.record.Update.Type + " record.")
		return nil
	case 1:
		// Partial execution: at least some actions succeeded; log as warning, not error.
		log.Warn("Partial update for " + c.record.Update.Name + ": " + response.Message)
		for _, d := range response.Details {
			log.Warn(d)
		}
		return nil
	default:
		for _, d := range response.Details {
			log.Error(d)
		}
		return fmt.Errorf("DNSExit API error %d: %s", response.Code, response.Message)
	}
}

func keepCurrent(ctx context.Context, c client, p chan client) {
	send := func() {
		select {
		case p <- c:
		case <-ctx.Done():
		}
	}

	// SELF type delegates IP detection to DNSExit server-side; no local comparison needed.
	if c.record.Update.Type == recordTypeSelf {
		if err := c.postUpdate(); err != nil {
			log.Error(err.Error())
		}
		send()
		return
	}

	desiredIP := c.record.Update.Content

	if !c.staticIP {
		ip, err := fetchEgressIP(c.http)
		if err != nil {
			log.Error("Failed to fetch egress IP: " + err.Error())
			send()
			return
		}
		desiredIP = ip
		c.record.Update.Content = ip
	}

	currentAddr, err := c.resolver.Resolve(ctx, c.record.Update.Name)
	if err != nil {
		log.Error(err.Error())
		send()
		return
	}

	if currentAddr == desiredIP {
		log.Info(c.record.Update.Name + " " + c.record.Update.Type + " record is up to date.")
		send()
		return
	}

	log.Info("Updating " + c.record.Update.Name + " " + c.record.Update.Type + " record from " + currentAddr + " to " + desiredIP + ".")
	if err := c.postUpdate(); err != nil {
		log.Error(err.Error())
	}
	send()
}
