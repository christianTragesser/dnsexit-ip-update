package dnsexit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPClient implements HTTPClient for testing.
// Responses are returned sequentially; requests and bodies are captured for inspection.
type mockHTTPClient struct {
	responses []*http.Response
	err       error
	calls     int
	requests  []*http.Request
	bodies    [][]byte
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Capture request and body before anything else.
	if req.Body != nil {
		body, _ := io.ReadAll(req.Body)
		m.bodies = append(m.bodies, body)
		req.Body = io.NopCloser(bytes.NewReader(body))
	} else {
		m.bodies = append(m.bodies, nil)
	}
	m.requests = append(m.requests, req)

	if m.err != nil {
		return nil, m.err
	}
	if m.calls >= len(m.responses) {
		return nil, fmt.Errorf("unexpected HTTP call %d", m.calls+1)
	}
	resp := m.responses[m.calls]
	m.calls++
	return resp, nil
}

func jsonResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func mockHTTP(bodies ...string) *mockHTTPClient {
	responses := make([]*http.Response, len(bodies))
	for i, body := range bodies {
		responses[i] = jsonResponse(body)
	}
	return &mockHTTPClient{responses: responses}
}

// decodePayload unmarshals a captured request body into an update struct.
func decodePayload(t *testing.T, raw []byte) update {
	t.Helper()
	var u update
	if err := json.Unmarshal(raw, &u); err != nil {
		t.Fatalf("failed to decode request body: %v", err)
	}
	return u
}

// mockResolver implements DomainResolver for testing.
type mockResolver struct {
	addr string
	err  error
}

func (m mockResolver) Resolve(_ context.Context, _ string) (string, error) {
	return m.addr, m.err
}

// ── splitDomain ────────────────────────────────────────────────────────────

func TestSplitDomain(t *testing.T) {
	tests := []struct {
		host string
		base string
		sub  string
	}{
		{"example.com", "example.com", ""},
		{"mysite.io", "mysite.io", ""},
		{"sub.example.com", "example.com", "sub"},
		{"api.mysite.io", "mysite.io", "api"},
		{"deep.sub.example.com", "example.com", "deep.sub"},
	}

	for _, tc := range tests {
		t.Run(tc.host, func(t *testing.T) {
			base, sub := splitDomain(tc.host)
			if base != tc.base || sub != tc.sub {
				t.Errorf("splitDomain(%q) = (%q, %q), want (%q, %q)", tc.host, base, sub, tc.base, tc.sub)
			}
		})
	}
}

// ── fetchEgressIP ──────────────────────────────────────────────────────────

func TestFetchEgressIP(t *testing.T) {
	tests := []struct {
		name    string
		client  *mockHTTPClient
		want    string
		wantErr bool
	}{
		{
			name:   "successful discovery",
			client: mockHTTP(`{"ip":"203.0.113.5"}`),
			want:   "203.0.113.5",
		},
		{
			name:    "HTTP request fails",
			client:  &mockHTTPClient{err: fmt.Errorf("connection refused")},
			wantErr: true,
		},
		{
			name:    "invalid JSON response",
			client:  mockHTTP(`not json`),
			wantErr: true,
		},
		{
			name:    "empty IP in response",
			client:  mockHTTP(`{"ip":""}`),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := fetchEgressIP(tc.client)
			if (err != nil) != tc.wantErr {
				t.Errorf("fetchEgressIP() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

// ── postUpdate ─────────────────────────────────────────────────────────────

func newTestClient(hc *mockHTTPClient, name string) client {
	return client{
		url:    "http://test.invalid",
		apiKey: "key123",
		record: update{Update: updateRecord{
			Type:      recordTypeA,
			Name:      name,
			Content:   "1.1.1.1",
			TTL:       5,
			Overwrite: true,
		}},
		http: hc,
	}
}

func TestPostUpdate(t *testing.T) {
	tests := []struct {
		name    string
		domain  string
		hc      *mockHTTPClient
		wantErr bool
		errMsg  string
	}{
		{
			name:   "code 0 is success",
			domain: "example.com",
			hc:     mockHTTP(`{"code":0,"message":"OK"}`),
		},
		{
			name:   "code 1 is warning, not error",
			domain: "example.com",
			hc:     mockHTTP(`{"code":1,"message":"partial","details":["record A updated","record B skipped"]}`),
		},
		{
			name:    "code 2 is auth error",
			domain:  "example.com",
			hc:      mockHTTP(`{"code":2,"message":"bad key"}`),
			wantErr: true,
			errMsg:  "DNSExit API error 2",
		},
		{
			name:    "non-zero code includes message",
			domain:  "example.com",
			hc:      mockHTTP(`{"code":3,"message":"missing requirements"}`),
			wantErr: true,
			errMsg:  "missing requirements",
		},
		{
			name:    "HTTP request fails",
			domain:  "example.com",
			hc:      &mockHTTPClient{err: fmt.Errorf("connection refused")},
			wantErr: true,
		},
		{
			name:    "invalid JSON response",
			domain:  "example.com",
			hc:      mockHTTP(`not json`),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := newTestClient(tc.hc, tc.domain)
			err := c.postUpdate()
			if (err != nil) != tc.wantErr {
				t.Errorf("postUpdate() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if tc.errMsg != "" && (err == nil || !strings.Contains(err.Error(), tc.errMsg)) {
				t.Errorf("error %q does not contain %q", err, tc.errMsg)
			}
		})
	}
}

func TestPostUpdateDomainHeader(t *testing.T) {
	tests := []struct {
		fullDomain  string
		wantHeader  string
		wantPayload string
	}{
		{"example.com", "example.com", ""},
		{"sub.example.com", "example.com", "sub"},
		{"deep.sub.example.com", "example.com", "deep.sub"},
	}

	for _, tc := range tests {
		t.Run(tc.fullDomain, func(t *testing.T) {
			hc := mockHTTP(`{"code":0,"message":"OK"}`)
			c := newTestClient(hc, tc.fullDomain)

			if err := c.postUpdate(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify domain header contains base domain only.
			if got := hc.requests[0].Header.Get("domain"); got != tc.wantHeader {
				t.Errorf("domain header: got %q, want %q", got, tc.wantHeader)
			}

			// Verify name field in payload contains subdomain only.
			payload := decodePayload(t, hc.bodies[0])
			if payload.Update.Name != tc.wantPayload {
				t.Errorf("payload name: got %q, want %q", payload.Update.Name, tc.wantPayload)
			}
		})
	}
}

func TestPostUpdateOverwrite(t *testing.T) {
	hc := mockHTTP(`{"code":0,"message":"OK"}`)
	c := newTestClient(hc, "example.com")
	c.record.Update.Overwrite = true

	if err := c.postUpdate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	payload := decodePayload(t, hc.bodies[0])
	if !payload.Update.Overwrite {
		t.Errorf("expected overwrite=true in API payload, got false")
	}
}

// ── keepCurrent ────────────────────────────────────────────────────────────

func TestKeepCurrent(t *testing.T) {
	const apiOK = `{"code":0,"message":"OK"}`

	tests := []struct {
		name      string
		client    client
		wantCalls int
	}{
		{
			name: "static IP, record is current, no update sent",
			client: client{
				record:   update{Update: updateRecord{Type: recordTypeA, Name: "example.com", Content: "1.1.1.1"}},
				staticIP: true,
				http:     &mockHTTPClient{},
				resolver: mockResolver{addr: "1.1.1.1"},
			},
			wantCalls: 0,
		},
		{
			name: "static IP, record is stale, update sent",
			client: client{
				record:   update{Update: updateRecord{Type: recordTypeA, Name: "example.com", Content: "2.2.2.2"}},
				staticIP: true,
				http:     mockHTTP(apiOK),
				resolver: mockResolver{addr: "1.1.1.1"},
			},
			wantCalls: 1,
		},
		{
			name: "DNS resolution error, update skipped",
			client: client{
				record:   update{Update: updateRecord{Type: recordTypeA, Name: "example.com", Content: "1.1.1.1"}},
				staticIP: true,
				http:     &mockHTTPClient{},
				resolver: mockResolver{err: fmt.Errorf("lookup failed")},
			},
			wantCalls: 0,
		},
		{
			name: "dynamic IP, egress matches DNS, no update sent",
			client: client{
				record:   update{Update: updateRecord{Type: recordTypeA, Name: "example.com"}},
				staticIP: false,
				http:     mockHTTP(`{"ip":"1.1.1.1"}`),
				resolver: mockResolver{addr: "1.1.1.1"},
			},
			wantCalls: 1,
		},
		{
			name: "dynamic IP, egress differs from DNS, update sent",
			client: client{
				record:   update{Update: updateRecord{Type: recordTypeA, Name: "example.com"}},
				staticIP: false,
				http:     mockHTTP(`{"ip":"2.2.2.2"}`, apiOK),
				resolver: mockResolver{addr: "1.1.1.1"},
			},
			wantCalls: 2,
		},
		{
			name: "dynamic IP, egress fetch fails, update skipped",
			client: client{
				record:   update{Update: updateRecord{Type: recordTypeA, Name: "example.com"}},
				staticIP: false,
				http:     &mockHTTPClient{err: fmt.Errorf("no network")},
				resolver: mockResolver{addr: "1.1.1.1"},
			},
			wantCalls: 0,
		},
		{
			name: "SELF type always sends update, no DNS resolution",
			client: client{
				record:   update{Update: updateRecord{Type: recordTypeSelf, Name: "example.com"}},
				staticIP: true,
				http:     mockHTTP(apiOK),
				resolver: mockResolver{err: fmt.Errorf("should not be called")},
			},
			wantCalls: 1,
		},
		{
			name: "AAAA type behaves like A type when record is stale",
			client: client{
				record:   update{Update: updateRecord{Type: recordTypeAAAA, Name: "example.com", Content: "2001:db8::1"}},
				staticIP: true,
				http:     mockHTTP(apiOK),
				resolver: mockResolver{addr: "2001:db8::2"},
			},
			wantCalls: 1,
		},
		{
			name: "AAAA type behaves like A type when record is current",
			client: client{
				record:   update{Update: updateRecord{Type: recordTypeAAAA, Name: "example.com", Content: "2001:db8::1"}},
				staticIP: true,
				http:     &mockHTTPClient{},
				resolver: mockResolver{addr: "2001:db8::1"},
			},
			wantCalls: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ch := make(chan client, 1)
			keepCurrent(context.Background(), tc.client, ch)
			<-ch

			hc := tc.client.http.(*mockHTTPClient)
			if hc.calls != tc.wantCalls {
				t.Errorf("HTTP calls = %d, want %d", hc.calls, tc.wantCalls)
			}
		})
	}
}

func TestKeepCurrentSelfSkipsDNS(t *testing.T) {
	// Verify that SELF type sends the update without ever touching the resolver.
	resolverCalled := false
	resolver := &callTrackingResolver{onResolve: func() { resolverCalled = true }}

	c := client{
		record:   update{Update: updateRecord{Type: recordTypeSelf, Name: "example.com"}},
		staticIP: true,
		http:     mockHTTP(`{"code":0,"message":"OK"}`),
		resolver: resolver,
	}

	ch := make(chan client, 1)
	keepCurrent(context.Background(), c, ch)
	<-ch

	if resolverCalled {
		t.Error("resolver was called for SELF type; expected no DNS resolution")
	}
}

// callTrackingResolver records whether Resolve was invoked.
type callTrackingResolver struct {
	onResolve func()
}

func (r *callTrackingResolver) Resolve(_ context.Context, _ string) (string, error) {
	r.onResolve()
	return "1.1.1.1", nil
}
