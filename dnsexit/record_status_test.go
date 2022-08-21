package dnsexit

import (
	"testing"
)

type mockRecordStatus struct {
	current   bool
	currentIP string
	desiredIP string
	domain    string
}

func (c mockRecordStatus) getRecordIP(domain string) string {
	return c.currentIP
}

func (d mockRecordStatus) getDesiredIP() string {
	return d.desiredIP
}

func TestRecordCheck(t *testing.T) {
	tests := []struct {
		name      string
		currentIP string
		desiredIP string
		domain    string
		expect    bool
	}{
		{
			name:      "A record update",
			domain:    "test.io",
			currentIP: "1.1.1.1",
			desiredIP: "2.2.2.2",
			expect:    false,
		},
		{
			name:      "A record current",
			domain:    "test.io",
			currentIP: "1.1.1.1",
			desiredIP: "1.1.1.1",
			expect:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			statusAPI := mockRecordStatus{}
			recordStatus := recordStatus{
				currentIP: tc.currentIP,
				desiredIP: tc.desiredIP,
				domain:    tc.domain,
			}

			got := getRecordStatus(statusAPI, recordStatus)

			if got != tc.expect {
				t.Errorf("getRecordStatus unit test failure '%s'\n got: '%v'\n want: '%v'", tc.name, got, tc.expect)
			}
		})
	}
}
