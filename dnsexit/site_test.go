package dnsexit

import (
	"reflect"
	"strings"
	"testing"
)

func TestGetDomains(t *testing.T) {
	tests := []struct {
		name   string
		site   site
		envs   string
		err    string
		expect []string
	}{
		{
			name:   "single site provided by flags",
			site:   site{domains: "example.com"},
			expect: []string{"example.com"},
		},
		{
			name:   "multiple sites provided by flags",
			site:   site{domains: "example.com,test.io"},
			expect: []string{"example.com", "test.io"},
		},
		{
			name:   "env var provided, no flags",
			site:   site{},
			envs:   "example.com,test.io",
			expect: []string{"example.com", "test.io"},
		},
		{
			name:   "domain name not found",
			site:   site{},
			err:    "domain name(s) not found",
			expect: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envs != "" {
				t.Setenv("DOMAINS", tc.envs)
			}

			got, err := tc.site.GetDomains()
			if err != nil {
				if !strings.Contains(err.Error(), tc.err) {
					t.Errorf("GetDomains() error = %q, want %q", err, tc.err)
				}
				return
			}

			if !reflect.DeepEqual(got, tc.expect) {
				t.Errorf("got %v, want %v", got, tc.expect)
			}
		})
	}
}

func TestGetAPIKey(t *testing.T) {
	tests := []struct {
		name   string
		site   site
		envs   string
		err    string
		expect string
	}{
		{
			name:   "key provided by flags",
			site:   site{key: "12345"},
			expect: "12345",
		},
		{
			name:   "key provided by env var",
			site:   site{},
			envs:   "12345",
			expect: "12345",
		},
		{
			name:   "no key provided",
			site:   site{},
			err:    "API key not found",
			expect: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envs != "" {
				t.Setenv("API_KEY", tc.envs)
			}

			got, err := tc.site.GetAPIKey()
			if err != nil {
				if !strings.Contains(err.Error(), tc.err) {
					t.Errorf("GetAPIKey() error = %q, want %q", err, tc.err)
				}
				return
			}

			if got != tc.expect {
				t.Errorf("got %q, want %q", got, tc.expect)
			}
		})
	}
}

func TestGetIPAddr(t *testing.T) {
	tests := []struct {
		name         string
		site         site
		envAddr      string
		err          string
		expectIP     string
		expectStatic bool
	}{
		{
			name:         "valid address provided by flag",
			site:         site{address: "1.1.1.1"},
			expectIP:     "1.1.1.1",
			expectStatic: true,
		},
		{
			name:         "valid IPv6 address provided by flag",
			site:         site{address: "2001:db8::1"},
			expectIP:     "2001:db8::1",
			expectStatic: true,
		},
		{
			name:         "valid address provided by env var",
			site:         site{},
			envAddr:      "1.1.1.1",
			expectIP:     "1.1.1.1",
			expectStatic: true,
		},
		{
			name: "invalid address provided by flag",
			site: site{address: "1.1.1.256"},
			err:  "invalid IP address",
		},
		{
			name:    "invalid address provided by env var",
			site:    site{},
			envAddr: "1.1.1.256",
			err:     "invalid IP address",
		},
		{
			name: "auto-discover egress IP",
			site: site{
				httpClient: mockHTTP(`{"ip":"203.0.113.5"}`),
			},
			expectIP:     "203.0.113.5",
			expectStatic: false,
		},
		{
			name: "auto-discover fails",
			site: site{
				httpClient: &mockHTTPClient{err: errConnRefused},
			},
			err: "egress IP request failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envAddr != "" {
				t.Setenv("IP_ADDR", tc.envAddr)
			}

			gotIP, gotStatic, err := tc.site.GetIPAddr()
			if err != nil {
				if !strings.Contains(err.Error(), tc.err) {
					t.Errorf("GetIPAddr() error = %q, want %q", err, tc.err)
				}
				return
			}

			if gotIP != tc.expectIP {
				t.Errorf("IP: got %q, want %q", gotIP, tc.expectIP)
			}
			if gotStatic != tc.expectStatic {
				t.Errorf("static: got %v, want %v", gotStatic, tc.expectStatic)
			}
		})
	}
}

// errConnRefused is a sentinel error for tests that need a failing HTTP client.
var errConnRefused = &testError{"connection refused"}

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }

func TestGetRecordType(t *testing.T) {
	tests := []struct {
		name    string
		site    site
		envs    string
		expect  string
		wantErr bool
	}{
		{
			name:   "default record type is A",
			site:   site{},
			expect: recordTypeA,
		},
		{
			name:   "A type provided by flag",
			site:   site{recordType: "A"},
			expect: recordTypeA,
		},
		{
			name:   "AAAA type provided by flag",
			site:   site{recordType: "AAAA"},
			expect: recordTypeAAAA,
		},
		{
			name:   "SELF type provided by flag",
			site:   site{recordType: "SELF"},
			expect: recordTypeSelf,
		},
		{
			name:   "AAAA type provided by env var",
			site:   site{},
			envs:   "AAAA",
			expect: recordTypeAAAA,
		},
		{
			name:   "flag value is case insensitive",
			site:   site{recordType: "aaaa"},
			expect: recordTypeAAAA,
		},
		{
			name:   "env var is case insensitive",
			site:   site{},
			envs:   "self",
			expect: recordTypeSelf,
		},
		{
			name:    "unsupported type MX returns error",
			site:    site{recordType: "MX"},
			wantErr: true,
		},
		{
			name:    "unsupported type via env var returns error",
			site:    site{},
			envs:    "TXT",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envs != "" {
				t.Setenv("RECORD_TYPE", tc.envs)
			}

			got, err := tc.site.GetRecordType()
			if (err != nil) != tc.wantErr {
				t.Errorf("GetRecordType() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if got != tc.expect {
				t.Errorf("got %q, want %q", got, tc.expect)
			}
		})
	}
}

func TestGetTTL(t *testing.T) {
	tests := []struct {
		name   string
		site   site
		envs   string
		expect int
	}{
		{
			name:   "TTL provided by flag",
			site:   site{ttl: 30, ttlSet: true},
			expect: 30,
		},
		{
			name:   "flag equals default but ttlSet, flag wins over env var",
			site:   site{ttl: defaultTTL, ttlSet: true},
			envs:   "60",
			expect: defaultTTL,
		},
		{
			name:   "TTL provided by env var",
			site:   site{ttl: defaultTTL, ttlSet: false},
			envs:   "30",
			expect: 30,
		},
		{
			name:   "invalid TTL in env var falls back to default",
			site:   site{ttl: defaultTTL, ttlSet: false},
			envs:   "bad",
			expect: defaultTTL,
		},
		{
			name:   "TTL of 0 is valid for Let's Encrypt",
			site:   site{ttl: 0, ttlSet: true},
			expect: 0,
		},
		{
			name:   "no TTL provided anywhere returns default",
			site:   site{ttl: defaultTTL, ttlSet: false},
			expect: defaultTTL,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envs != "" {
				t.Setenv("RECORD_TTL", tc.envs)
			}

			got := tc.site.GetTTL()
			if got != tc.expect {
				t.Errorf("got %d, want %d", got, tc.expect)
			}
		})
	}
}

func TestGetInterval(t *testing.T) {
	tests := []struct {
		name   string
		site   site
		envs   string
		expect int
	}{
		{
			name:   "interval provided by flag",
			site:   site{interval: 20, intervalSet: true},
			expect: 20,
		},
		{
			name:   "flag equals default but intervalSet, flag wins over env var",
			site:   site{interval: defaultInterval, intervalSet: true},
			envs:   "5",
			expect: defaultInterval,
		},
		{
			name:   "interval provided by env var",
			site:   site{interval: defaultInterval, intervalSet: false},
			envs:   "15",
			expect: 15,
		},
		{
			name:   "invalid interval in env var falls back to default",
			site:   site{interval: defaultInterval, intervalSet: false},
			envs:   "2.2",
			expect: defaultInterval,
		},
		{
			name:   "no interval provided anywhere returns default",
			site:   site{interval: defaultInterval, intervalSet: false},
			expect: defaultInterval,
		},
		{
			name:   "interval below minimum is clamped to minimum",
			site:   site{interval: 2, intervalSet: true},
			expect: minInterval,
		},
		{
			name:   "interval exactly at minimum is accepted",
			site:   site{interval: minInterval, intervalSet: true},
			expect: minInterval,
		},
		{
			name:   "env var interval below minimum is clamped",
			site:   site{interval: defaultInterval, intervalSet: false},
			envs:   "2",
			expect: minInterval,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envs != "" {
				t.Setenv("CHECK_INTERVAL", tc.envs)
			}

			got := tc.site.GetInterval()
			if got != tc.expect {
				t.Errorf("got %d, want %d", got, tc.expect)
			}
		})
	}
}
