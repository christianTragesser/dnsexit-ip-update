package dnsexit

type recordStatusAPI interface {
	getRecordIP(domain string) string
	getDesiredIP() string
}

type recordStatus struct {
	current   bool
	currentIP string
	desiredIP string
	domain    string
}

func (c recordStatus) getRecordIP(domain string) string {
	return c.currentIP
}

func (d recordStatus) getDesiredIP() string {
	return d.desiredIP
}

func getRecordStatus(api recordStatusAPI, status recordStatus) bool {
	return status.current
}
