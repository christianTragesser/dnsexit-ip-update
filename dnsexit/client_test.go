package dnsexit

import (
	"testing"
)

//nolint:deadcode,unused
type mockClientAPI interface {
	getDomain() string
	currentRecords() ([]string, error)
}

type mockClient struct {
	URL      string
	APIKey   string
	Record   updateRecord
	Interval int
}

func (c mockClient) getDomain() string {
	return c.Record.Name
}

func (c mockClient) currentRecords() ([]string, error) {
	return []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"}, nil
}

func TestCurrentRecords(t *testing.T) {
	tests := []struct {
		name   string
		client mockClient
		expect bool
	}{
		{
			name: "Not current",
			client: mockClient{
				URL:    apiURL,
				APIKey: "12345",
				Record: updateRecord{
					Type:    recordType,
					TTL:     recordTTL,
					Name:    "testing",
					Content: "4.4.4.4",
				},
				Interval: defaultInterval,
			},
			expect: false,
		},
		{
			name: "Current",
			client: mockClient{
				URL:    apiURL,
				APIKey: "12345",
				Record: updateRecord{
					Type:    recordType,
					TTL:     recordTTL,
					Name:    "testing",
					Content: "2.2.2.2",
				},
				Interval: defaultInterval,
			},
			expect: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ips, err := tc.client.currentRecords()
			if err != nil {
				t.Errorf("currentRecords unit test failure %s:\n %v", tc.name, err)
			}

			record := client{}

			got := record.current(ips, tc.client.Record.Content)

			if got != tc.expect {
				t.Errorf("currentRecords unit test failure %s:\n got:'%v'\nwant:'%v'", tc.name, got, tc.expect)
			}
		})
	}
}
