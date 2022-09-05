package dnsexit

import (
	"fmt"
	"reflect"
	"testing"
)

type mockAPIResponse struct {
	code    int
	details []string
	message string
}

func (m mockAPIResponse) setUpdate(event Event) (Event, error) {
	var err error

	mockEvent := Event{
		Code:    m.code,
		Details: m.details,
		Message: m.message,
	}

	if m.code > 7 {
		err = fmt.Errorf("This is a mocked exception.")
	}

	return mockEvent, err
}

func TestDynamicUpdate(t *testing.T) {
	// calls DNSexit Update API with host list
	// if update event, log and returns true
	// if non-update event, log and returns false
	// if response code is 2-7, logs and returns false + error

	tests := []struct {
		name     string
		resp     mockAPIResponse
		expected Event
	}{
		{
			name: "Update event true",
			resp: mockAPIResponse{
				code:    0,
				details: []string{"Update call success."},
				message: "Success",
			},
			expected: Event{
				Code:    0,
				Details: []string{"Update call success."},
				Message: "Success",
			},
		},
		{
			name: "Update event failure",
			resp: mockAPIResponse{
				code:    3,
				details: []string{"Update call failure."},
				message: "Exception",
			},
			expected: Event{
				Code:    3,
				Details: []string{"Update call failure."},
				Message: "Exception",
			},
		},
		{
			name: "Update event exception",
			resp: mockAPIResponse{
				code:    8,
				details: []string{"Update call exception."},
				message: "Exception",
			},
			expected: Event{
				Code:    8,
				Details: []string{"Update call exception."},
				Message: "Exception",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := dynamicUpdate(tc.resp, tc.expected)

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("DynamicUpdate unit test failure '%s'\n got: '%v'\n want: '%v'", tc.name, got, tc.expected)
			}
		})
	}
}
