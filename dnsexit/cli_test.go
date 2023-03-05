package dnsexit

import (
	"os"
	"strings"
	"testing"
)

func TestSetUpdateData(t *testing.T) {
	tests := []struct {
		name   string
		cmd    CLICommand
		envs   map[string]string
		err    string
		expect updateRecord
	}{
		{
			name: "Flags",
			cmd: CLICommand{
				domain:   "testing",
				key:      "12345",
				interval: 1,
				address:  "1.1.1.1",
			},
			envs: map[string]string{},
			expect: updateRecord{
				Type:    recordType,
				TTL:     recordTTL,
				Name:    "testing",
				Content: "1.1.1.1",
			},
		},
		{
			name: "Env Vars",
			cmd: CLICommand{
				domain:   "",
				key:      "",
				interval: 1,
				address:  "1.1.1.1",
			},
			envs: map[string]string{"DOMAIN": "testing"},
			expect: updateRecord{
				Type:    recordType,
				TTL:     recordTTL,
				Name:    "testing",
				Content: "1.1.1.1",
			},
		},
		{
			name: "Domain Name Not Found",
			cmd: CLICommand{
				domain:   "",
				key:      "12345",
				interval: 1,
				address:  "1.1.1.1",
			},
			envs: map[string]string{},
			err:  "Missing DNSExit domain name(s).",
			expect: updateRecord{
				Type:    recordType,
				TTL:     recordTTL,
				Name:    "",
				Content: "1.1.1.1",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envs {
				os.Setenv(k, v)
			}
			defer os.Unsetenv("DOMAIN")

			got, err := tc.cmd.setUpdateData()
			if err != nil {
				if !strings.Contains(err.Error(), tc.err) {
					t.Errorf("setUpdateData unit test failure %s:\n got error: '%v'\nwant error: '%v'", tc.name, err, tc.err)
				}
			}

			if got != tc.expect {
				t.Errorf("setUpdateData unit test failure %s:\n got:'%v'\nwant: '%v'", tc.name, got, tc.expect)
			}
		})
	}
}

func TestSetClient(t *testing.T) {
	tests := []struct {
		name   string
		cmd    CLICommand
		update updateRecord
		envs   map[string]string
		err    string
		expect client
	}{
		{
			name: "Flags",
			cmd: CLICommand{
				domain:   "test",
				key:      "12345",
				interval: 10,
				address:  "1.1.1.1",
			},
			update: updateRecord{
				Type:    recordType,
				TTL:     recordTTL,
				Name:    "test",
				Content: "1.1.1.1",
			},
			envs: map[string]string{},
			expect: client{
				URL:    apiURL,
				APIKey: "12345",
				Record: updateRecord{
					Type:    recordType,
					TTL:     recordTTL,
					Name:    "test",
					Content: "1.1.1.1",
				},
				Interval: 10,
			},
		},
		{
			name: "Env Vars",
			cmd: CLICommand{
				domain:   "test",
				key:      "",
				interval: 10,
				address:  "1.1.1.1",
			},
			update: updateRecord{
				Type:    recordType,
				TTL:     recordTTL,
				Name:    "test",
				Content: "1.1.1.1",
			},
			envs: map[string]string{"API_KEY": "12345", "CHECK_INTERVAL": "20"},
			expect: client{
				URL:    apiURL,
				APIKey: "12345",
				Record: updateRecord{
					Type:    recordType,
					TTL:     recordTTL,
					Name:    "test",
					Content: "1.1.1.1",
				},
				Interval: 20,
			},
		},
		{
			name: "API Key Not Found",
			cmd: CLICommand{
				domain:   "test",
				key:      "",
				interval: 10,
				address:  "1.1.1.1",
			},
			update: updateRecord{
				Type:    recordType,
				TTL:     recordTTL,
				Name:    "test",
				Content: "1.1.1.1",
			},
			envs: map[string]string{},
			err:  "Missing DNSExit API Key.",
			expect: client{
				URL:    apiURL,
				APIKey: "",
				Record: updateRecord{
					Type:    recordType,
					TTL:     recordTTL,
					Name:    "test",
					Content: "1.1.1.1",
				},
				Interval: 10,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envs {
				os.Setenv(k, v)
			}
			defer os.Unsetenv("API_KEY")
			defer os.Unsetenv("CHECK_INTERVAL")

			got, err := tc.cmd.setClient(tc.update)
			if err != nil {
				if !strings.Contains(err.Error(), tc.err) {
					t.Errorf("setClient unit test failure %s:\n got error: '%v'\nwant error: '%v'", tc.name, err, tc.err)
				}
			}

			if got != tc.expect {
				t.Errorf("setClient unit test failure %s:\n got:'%v'\nwant: '%v'", tc.name, got, tc.expect)
			}
		})
	}
}
