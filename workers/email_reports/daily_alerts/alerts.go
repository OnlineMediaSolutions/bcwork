package daily_alerts

type RequestData struct {
	Data RequestDetails `json:"data"`
}
type RequestDetails struct {
	Date       Date     `json:"date"`
	Dimensions []string `json:"dimensions"`
	Metrics    []string `json:"metrics"`
}
type Date struct {
	Range    []string `json:"range"`
	Interval string   `json:"interval"`
}

type Alerts interface {
	Request() ([]AggregatedReport, error)
	GetRequestData() RequestData
	Aggregate(report []AggregatedReport) map[string][]AggregatedReport
	ComputeAverage(map[string][]AggregatedReport) map[string][]AlertsEmails
	PrepareAndSendEmail(map[string][]AlertsEmails, *Worker) error
	SendCustomHTMLEmail(to, bcc, subject string, body string, report []AlertsEmails) error
	GenerateHTMLTableWithTemplate(report []AlertsEmails, body string) (string, error)
}

func GetAlerts(alertType string) Alerts {
	switch alertType {
	case "LOOPING_RATIO_DECREASE":
		return &LoopingRationDecreaseReport{}
	case "RPM_DECREASE":
		return &RPMDecreaseReport{}
	default:
		return nil
	}
}
