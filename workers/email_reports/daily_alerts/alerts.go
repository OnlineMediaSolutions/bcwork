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
	Request(*Worker) ([]AggregatedReport, error)
	GetRequestData(*Worker) RequestData
	Aggregate(report []AggregatedReport) map[string][]AggregatedReport
	ComputeAverage(map[string][]AggregatedReport, *Worker) map[string][]AlertsEmailRepo
	PrepareAndSendEmail(map[string][]AlertsEmailRepo, *Worker) error
	SendCustomHTMLEmail(to, bcc, subject string, body string, report []AlertsEmailRepo) error
	GenerateHTMLTableWithTemplate(report []AlertsEmailRepo, body string) (string, error)
}

func GetAlerts(reportType string) Alerts {
	switch reportType {
	case "LOOPING_RATIO_DECREASE":
		return &LoopingRationDecreaseReport{}
	case "RPM_DECREASE":
		return &RPMDecreaseReport{}
	default:
		return nil
	}
}
