package dnsexit

import (
	"os"
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
			name: "single site provided by flags",
			site: site{
				domains: "example.com",
			},
			envs:   "",
			err:    "",
			expect: []string{"example.com"},
		},
		{
			name: "multiple sites provided by flags",
			site: site{
				domains: "example.com,test.io",
			},
			envs:   "",
			err:    "",
			expect: []string{"example.com", "test.io"},
		},
		{
			name:   "env vars provided but not flags",
			site:   site{},
			envs:   "example.com,test.io",
			err:    "",
			expect: []string{"example.com", "test.io"},
		},
		{
			name:   "domain name not found",
			site:   site{},
			envs:   "",
			err:    "domain name(s) not found",
			expect: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envs != "" {
				os.Setenv("DOMAINS", tc.envs)
			}
			defer os.Unsetenv("DOMAINS")

			got, err := tc.site.GetDomains()
			if err != nil {
				if !strings.Contains(err.Error(), tc.err) {
					t.Errorf("GetDomains unit test failure %s:\n got: '%v'\nwant: '%v'", tc.name, err, tc.err)
				}
			}

			if !reflect.DeepEqual(got, tc.expect) {
				t.Errorf("GetDomains unit test failure %s:\n got: '%v'\nwant: '%v'", tc.name, got, tc.expect)
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
			envs:   "",
			err:    "",
			expect: "12345",
		},
		{
			name:   "key provided by env vars",
			site:   site{key: ""},
			envs:   "12345",
			err:    "",
			expect: "12345",
		},
		{
			name:   "no key provided",
			site:   site{key: ""},
			envs:   "",
			err:    "API key not found",
			expect: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envs != "" {
				os.Setenv("API_KEY", tc.envs)
			}
			defer os.Unsetenv("API_KEY")

			got, err := tc.site.GetAPIKey()
			if err != nil {
				if !strings.Contains(err.Error(), tc.err) {
					t.Errorf("GetAPIKey unit test failure %s:\n got: '%v'\nwant: '%v'", tc.name, err, tc.err)
				}
			}

			if got != tc.expect {
				t.Errorf("GetAPIKey unit test failure %s:\n got: '%v'\nwant: '%v'", tc.name, got, tc.expect)
			}
		})
	}
}

func TestGetIPAddr(t *testing.T) {
	tests := []struct {
		name   string
		site   site
		envs   string
		err    string
		expect string
	}{
		{
			name:   "valid address provided by flags",
			site:   site{address: "1.1.1.1"},
			envs:   "",
			err:    "",
			expect: "1.1.1.1",
		},
		{
			name:   "valid address provided by env vars",
			site:   site{address: ""},
			envs:   "1.1.1.1",
			err:    "",
			expect: "1.1.1.1",
		},
		{
			name:   "invalid address provided by flags",
			site:   site{address: "1.1.1.256"},
			envs:   "",
			err:    "1.1.1.256 is an invalid IP address",
			expect: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envs != "" {
				os.Setenv("IP_ADDR", tc.envs)
			}
			defer os.Unsetenv("IP_ADDR")

			got, err := tc.site.GetIPAddr()
			if err != nil {
				if !strings.Contains(err.Error(), tc.err) {
					t.Errorf("GetIPAddr unit test failure %s:\n got: '%v'\nwant: '%v'", tc.name, err, tc.err)
				}
			}

			if got != tc.expect {
				t.Errorf("GetIPAddr unit test failure %s:\n got: '%v'\nwant: '%v'", tc.name, got, tc.expect)
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
			name:   "interval provided by flags",
			site:   site{interval: 2},
			envs:   "",
			expect: 2,
		},
		{
			name:   "interval provided by env vars",
			site:   site{interval: defaultInterval},
			envs:   "2",
			expect: 2,
		},
		{
			name:   "invalid interval provided",
			site:   site{interval: defaultInterval},
			envs:   "2.2",
			expect: defaultInterval,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envs != "" {
				os.Setenv("CHECK_INTERVAL", tc.envs)
			}
			defer os.Unsetenv("CHECK_INTERVAL")

			got := tc.site.GetInterval()
			if got != tc.expect {
				t.Errorf("GetInterval unit test failure %s:\n got: '%v'\nwant: '%v'", tc.name, got, tc.expect)
			}
		})
	}
}
