package dnsexit

import (
	"os"
	"strings"
	"testing"
)

type MockIPAddrAPI interface {
	getContent() string
	egressIP() (string, error)
}

type mockUpdateRecord struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

func (m mockUpdateRecord) getContent() string {
	return m.Content
}

func (m mockUpdateRecord) egressIP() (string, error) {
	return "3.3.3.3", nil
}

func TestGetUpdateIP(t *testing.T) {
	tests := []struct {
		name   string
		record mockUpdateRecord
		envs   map[string]string
		err    string
		expect string
	}{
		{
			name: "Flags",
			record: mockUpdateRecord{
				Type:    recordType,
				TTL:     recordTTL,
				Name:    "testing",
				Content: "1.1.1.1",
			},
			envs:   map[string]string{},
			expect: "1.1.1.1",
		},
		{
			name: "Env Vars",
			record: mockUpdateRecord{
				Type:    recordType,
				TTL:     recordTTL,
				Name:    "testing",
				Content: "",
			},
			envs:   map[string]string{"IP_ADDR": "2.2.2.2"},
			expect: "2.2.2.2",
		},
		{
			name: "Dynamic IP",
			record: mockUpdateRecord{
				Type:    recordType,
				TTL:     recordTTL,
				Name:    "testing",
				Content: "",
			},
			envs:   map[string]string{},
			expect: "3.3.3.3",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envs {
				os.Setenv(k, v)
			}
			defer os.Unsetenv("IP_ADDR")

			got, err := getUpdateIP(tc.record)
			if err != nil {
				if !strings.Contains(err.Error(), tc.err) {
					t.Errorf("getUpdateIP unit test failure %s:\n got error: '%v'\nwant error: '%v'", tc.name, err, tc.err)
				}
			}

			if got != tc.expect {
				t.Errorf("getUpdateIP unit test failure %s:\n got:'%v'\nwant: '%v'", tc.name, got, tc.expect)
			}

		})
	}
}
