package testapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/m6yf/bcwork/config"
	httpclient "github.com/m6yf/bcwork/modules/http_client"
	"github.com/m6yf/bcwork/modules/messager"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	_ "github.com/sirupsen/logrus"
)

var (
	timePattern     = `\d{4}-[0-1]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d{1,9}Z`
	createdAtRegexp = regexp.MustCompile(`"created_at":"` + timePattern + `",|,"created_at":"` + timePattern + `"`)
	updatedAtRegexp = regexp.MustCompile(`"updated_at":"` + timePattern + `",|,"updated_at":"` + timePattern + `"`)
)

type Worker struct {
	BaseURL     string `json:"base_url"`
	LogSeverity int    `json:"logsev"`
	Cron        string `json:"cron"`

	cases []testCase

	messager   messager.Messager
	httpClient httpclient.Doer
}

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {
	const (
		baseURLDefault     = "https://api.nanoook.com"
		cronDefault        = "0 * * * *" // every hour
		logSeverityDefault = 2
	)

	w.BaseURL = conf.GetStringValueWithDefault(config.BaseURLKey, baseURLDefault)
	w.Cron = conf.GetStringValueWithDefault(config.CronExpressionKey, cronDefault)

	logSeverity, err := conf.GetIntValueWithDefault(config.LogSeverityKey, logSeverityDefault)
	if err != nil {
		return eris.Wrapf(err, "failed to parse log severity")
	}

	w.LogSeverity = logSeverity
	zerolog.SetGlobalLevel(zerolog.Level(w.LogSeverity))

	slackMod, err := messager.NewSlackModule()
	if err != nil {
		return eris.Wrapf(err, "failed to init slack module")
	}
	w.messager = slackMod

	w.httpClient = httpclient.New(true)

	file, err := os.Open("test_cases.json")
	if err != nil {
		return eris.Wrapf(err, "failed to open file with test cases")
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return eris.Wrapf(err, "failed to read file with test cases")
	}

	var testCases []testCase
	err = json.Unmarshal(data, &testCases)
	if err != nil {
		return eris.Wrapf(err, "failed to unmarshal test cases")
	}

	w.cases = testCases

	return nil
}

func (w *Worker) Do(ctx context.Context) error {
	log.Info().Msg("starting testing API")

	errReport := make([][]string, 0)
	for _, testCase := range w.cases {
		err := w.processTestCase(ctx, testCase)
		if err != nil {
			log.Err(err).Msgf("FAIL [%v]", testCase.Name)
			name := fmt.Sprintf("%v [%v %v]", testCase.Name, testCase.Method, testCase.Endpoint)
			errReport = append(errReport, []string{name, err.Error()})
		}
	}

	if len(errReport) > 0 {
		message := prepareMessage(errReport)
		err := w.messager.SendMessage(message)
		if err != nil {
			return fmt.Errorf("could not send error report message: %w", err)
		}
		return fmt.Errorf("testing API finished with errors, amount of errors: %v", len(errReport))
	}

	log.Info().Msg("testing API passed successfully")
	return nil
}

func (w *Worker) GetSleep() int {
	next := bccron.Next(w.Cron)
	log.Info().Msg(fmt.Sprintf("next run in: %v", time.Duration(next)*time.Second))
	if w.Cron != "" {
		return next
	}
	return 0
}

func (w *Worker) processTestCase(ctx context.Context, testCase testCase) error {
	data, _, err := w.httpClient.Do(ctx, testCase.Method, w.BaseURL+testCase.Endpoint, strings.NewReader(testCase.Payload))
	if err != nil {
		return fmt.Errorf("error while doing request: %w", err)
	}

	got := prepareData(data)

	if got != testCase.Want {
		return fmt.Errorf("not equal:\ngot  = %v\nwant = %v", got, testCase.Want)
	}

	return nil
}

func prepareData(data []byte) string {
	return updatedAtRegexp.ReplaceAllString(
		createdAtRegexp.ReplaceAllString(
			string(data),
			"",
		),
		"",
	)
}

func prepareMessage(report [][]string) string {
	message := "*Test API worker. Failed tests:*\n"
	var sep = "\n"

	for i, err := range report {
		if i+1 == len(report) {
			sep = ""
		}
		message += fmt.Sprintf("%v. _%v_: ```%v```%v", i+1, err[0], err[1], sep)
		i++
	}

	return message
}
