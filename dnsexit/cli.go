package dnsexit

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	apiURL          string = "https://api.dnsexit.com/dns/"
	defaultTTL      int    = 5
	defaultInterval int    = 10
	minInterval     int    = 5
)

var log = getLogger()

func getLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func CLI() {
	cliDomains := flag.String("domains", "", "DNSExit domain name(s), comma-separated")
	cliKey := flag.String("key", "", "DNSExit API key")
	cliIPAddr := flag.String("ip", "", "Desired record IP address (auto-discovered if omitted)")
	cliRecordType := flag.String("record-type", "", "DNS record type: A, AAAA, or SELF (default: A)")
	cliTTL := flag.Int("ttl", defaultTTL, "Record TTL in minutes")
	cliInterval := flag.Int("interval", defaultInterval, "Update check interval in minutes (minimum 5)")

	flag.Parse()

	intervalSet, ttlSet := false, false
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "interval":
			intervalSet = true
		case "ttl":
			ttlSet = true
		}
	})

	httpClient := &http.Client{}

	s := site{
		domains:     *cliDomains,
		key:         *cliKey,
		interval:    *cliInterval,
		intervalSet: intervalSet,
		address:     *cliIPAddr,
		recordType:  *cliRecordType,
		ttl:         *cliTTL,
		ttlSet:      ttlSet,
		httpClient:  httpClient,
	}

	domains, err := s.GetDomains()
	if err != nil {
		os.Exit(1)
	}

	apiKey, err := s.GetAPIKey()
	if err != nil {
		os.Exit(1)
	}

	recordType, err := s.GetRecordType()
	if err != nil {
		os.Exit(1)
	}

	ttl := s.GetTTL()
	interval := s.GetInterval()

	// SELF type delegates IP detection to DNSExit server-side; no local IP needed.
	var ipAddr string
	var staticIP bool
	if recordType != recordTypeSelf {
		ipAddr, staticIP, err = s.GetIPAddr()
		if err != nil {
			os.Exit(1)
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	clients := make([]client, 0, len(domains))

	for _, d := range domains {
		clients = append(clients, client{
			url:    apiURL,
			apiKey: apiKey,
			record: update{
				Update: updateRecord{
					Type:      recordType,
					TTL:       ttl,
					Name:      d,
					Content:   ipAddr,
					Overwrite: true,
				},
			},
			interval: interval,
			staticIP: staticIP || recordType == recordTypeSelf,
			http:     httpClient,
			resolver: netResolver{},
		})
	}

	channel := make(chan client)

	for _, c := range clients {
		go keepCurrent(ctx, c, channel)
	}

	for {
		select {
		case <-ctx.Done():
			log.Info("Shutting down.")
			return
		case c := <-channel:
			go func(c client) {
				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Duration(c.interval) * time.Minute):
					keepCurrent(ctx, c, channel)
				}
			}(c)
		}
	}
}
