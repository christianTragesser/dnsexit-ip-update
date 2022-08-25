package dnsexit

import (
	"net"
	"testing"
)

type mockRecordStatus struct {
	currentRecord net.IP
	currentIP     string
}

func (c mockRecordStatus) getRecords(domain string) []net.IP {
	return []net.IP{c.currentRecord}
}

func (d mockRecordStatus) getLocationIP() string {
	return d.currentIP
}

func TestRecordCheck(t *testing.T) {
	tests := []struct {
		name          string
		currentIP     string
		currentRecord net.IP
		desiredIP     string
		domain        string
		expect        bool
	}{
		{
			name:          "A record update",
			domain:        "test.io",
			currentIP:     "1.1.1.1",
			currentRecord: net.IP{1, 2, 3, 4},
			desiredIP:     "2.2.2.2",
			expect:        false,
		},
		{
			name:          "A record current",
			domain:        "test.io",
			currentIP:     "1.1.1.1",
			currentRecord: net.IP{1, 1, 1, 1},
			desiredIP:     "1.1.1.1",
			expect:        true,
		},
		{
			name:          "No desired IP, no update",
			domain:        "test.io",
			currentIP:     "1.1.1.1",
			currentRecord: net.IP{1, 1, 1, 1},
			expect:        true,
		},
		{
			name:          "No desired IP, update",
			domain:        "test.io",
			currentIP:     "1.1.1.1",
			currentRecord: net.IP{1, 2, 3, 4},
			expect:        false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			statusAPI := mockRecordStatus{
				currentRecord: tc.currentRecord,
				currentIP:     tc.currentIP,
			}

			testRecord := updateRecord{
				Name:    tc.domain,
				Content: tc.desiredIP,
			}

			testEvent := Event{
				Record: testRecord,
			}

			got := recordIsCurrent(statusAPI, testEvent)

			if got != tc.expect {
				t.Errorf("recordIsCurrent unit test failure '%s'\n got: '%v'\n want: '%v'", tc.name, got, tc.expect)
			}
		})
	}
}
