package real_time_report

import (
	"fmt"
	"github.com/m6yf/bcwork/utils/helpers"
	"time"
)

func GenerateReportDetails(worker *Worker) (string, string, string) {
	body := fmt.Sprintf("Full Publisher Requests between %s - %s\n",
		helpers.FormatDate(worker.Start.Format(time.RFC3339)),
		helpers.FormatDate(worker.End.Format(time.RFC3339)))
	subject := fmt.Sprintf("Full Publisher Requests %s", helpers.FormatDate(worker.End.Format(time.RFC3339)))
	reportName := fmt.Sprintf("Full Publisher Requests_%s.csv", helpers.FormatDate(worker.End.Format(time.RFC3339)))

	return body, subject, reportName
}
