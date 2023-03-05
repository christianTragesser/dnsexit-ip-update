package dnsexit

type clientAPI interface {
	getDomain() string
	currentRecords() ([]string, error)
}

type client struct {
	URL      string
	APIKey   string
	Record   updateRecord
	Interval int
}

func (c client) getDomain() string {
	return c.Record.Name
}

func (c client) currentRecords() ([]string, error) {
	ips, err := resolve(c)
	if err != nil {
		clientLogFields["domain"] = c.getDomain()
		log.WithFields(clientLogFields).Error(err)

		return []string{}, err
	}

	return ips, nil
}

func (c client) current(currentRecords []string, address string) bool {
	if len(currentRecords) > 0 {
		for _, ip := range currentRecords {
			if ip == address {
				return true
			}
		}
	}

	return false
}

/*
func (r *client) postUpdate(event client) (updateResponse, error) {
	var response updateResponse

	updatePayload := map[string]updateRecord{"update": event.Record}
	jsonPayload, _ := json.Marshal(updatePayload)
	data := bytes.NewReader([]byte(jsonPayload))

	req, err := http.NewRequest("POST", event.URL, data)
	if err != nil {
		log.Errorln("Failed to create HTTP POST to DNSExit API.")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", event.APIKey)
	req.Header.Set("domain", event.Record.Name)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		log.WithFields(updateRecordLogFields).Error("API POST failed for dynamic update.")

		return response, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		log.Errorln("Failed to read API response body.")
	}

	err = json.Unmarshal(body, &response)

	if response.Code != 0 {
		updateRecordLogFields["status"] = response.Code

		log.WithFields(updateRecordLogFields).Error(response.Message)
	}

	return response, err
}
*/
