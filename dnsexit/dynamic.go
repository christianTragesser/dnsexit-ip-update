package dnsexit

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type apiResponse struct {
	code    int
	details []string
	message string
}

type dynamicUpdateEvent struct {
	resp apiResponse
	err  error
}

type dnsExitAPI interface {
	callDynamicUpdate(url string) (http.Response, error)
}

func (r apiResponse) callDynamicUpdate(url string) (http.Response, error) {
	var apiResponse http.Response
	var err error

	return apiResponse, err
}

func DynamicUpdate(api dnsExitAPI) (dynamicUpdateEvent, error) {
	var err error

	event := dynamicUpdateEvent{}

	return event, err

}
