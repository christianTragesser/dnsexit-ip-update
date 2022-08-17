package dnsexit

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type mockAPIResponse struct {
	code    int
	details []string
	message string
}

func (r mockAPIResponse) callDynamicUpdate(url string) (http.Response, error) {
	var err error

	respBody := []byte(fmt.Sprintf("%v", r))

	mockResponse := http.Response{
		Status:     "OK",
		StatusCode: 200,
		Proto:      "HTTP/1.0",
		Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
	}

	return mockResponse, err
}

func TestDynamicUpdate(t *testing.T) {
	// calls DNSexit Update API with host list
	// if update event, log and returns true
	// if non-update event, log and returns false
	// if response code is 2-7, logs and returns false + error

	tests := []struct {
		name     string
		resp     mockAPIResponse
		expected dynamicUpdateEvent
	}{
		{
			name: "Update event true",
			resp: mockAPIResponse{
				code:    0,
				details: []string{"Update happened."},
				message: "Success",
			},
			expected: dynamicUpdateEvent{
				resp: apiResponse{
					code:    0,
					details: []string{"Update happened."},
					message: "Success",
				},
				err: nil,
			},
		},
		{
			name: "Update event false",
			resp: mockAPIResponse{
				code:    1,
				details: []string{"Update did not happened."},
				message: "No Op.",
			},
			expected: dynamicUpdateEvent{
				resp: apiResponse{
					code:    1,
					details: []string{"Update did not happened."},
					message: "No Op.",
				},
				err: nil,
			},
		},
		{
			name: "Update event failuer",
			resp: mockAPIResponse{
				code:    2,
				details: []string{"Update exception."},
				message: "Exception",
			},
			expected: dynamicUpdateEvent{
				resp: apiResponse{
					code:    2,
					details: []string{"Update exception."},
					message: "Exception",
				},
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dnsEvent := tc.resp
			got, _ := DynamicUpdate(dnsEvent)

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("DynamicUpdateAPI unit test failure '%s'\n got: '%v'\n want: '%v'", tc.name, got, tc.expected)
			}
		})
	}
}
