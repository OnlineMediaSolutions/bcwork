package real_time_report

import (
	"fmt"
	"github.com/m6yf/bcwork/utils/helpers"
	"time"
)

func GenerateReportDetails(worker *Worker) (string, string, string) {
	body := fmt.Sprintf("Bid requests daily numbers - UTC timezone %s - %s\n",
		helpers.FormatDate(worker.Start.Format(time.RFC3339)),
		helpers.FormatDate(worker.End.Format(time.RFC3339)))
	subject := fmt.Sprintf("Real time reports %s", helpers.FormatDate(worker.End.Format(time.RFC3339)))
	reportName := fmt.Sprintf("Real time report_%s.csv", helpers.FormatDate(worker.End.Format(time.RFC3339)))

	return body, subject, reportName
}
