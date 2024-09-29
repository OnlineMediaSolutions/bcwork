package testapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
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

	messager   messager.Messager
	httpClient httpclient.HttpClient
}

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {
	const (
		baseURLDefault     = "https://api.nanoook.com"
		cronDefault        = "0 * * * *" // every hour
		logSeverityDefault = 2
	)

	w.BaseURL = conf.GetStringValueWithDefault(config.BaseURLKey, baseURLDefault) // TODO: add key to config
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

	w.httpClient = &http.Client{}

	return nil
}

func (w *Worker) Do(ctx context.Context) error {
	log.Info().Msg("starting testing API")

	errReport := make(map[string]string, 0)
	for _, testCase := range testCases {
		err := w.processTestCase(testCase)
		if err != nil {
			log.Error().Msgf("FAIL [%v]: %v", testCase.name, err)
			name := fmt.Sprintf("%v [%v %v]", testCase.name, testCase.method, testCase.endpoint)
			errReport[name] = err.Error()
		}
	}

	if len(errReport) > 0 {
		message := prepareMessage(errReport)
		log.Print(message)
		// TODO: send report
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

func (w *Worker) processTestCase(testCase testCase) error {
	payload := strings.NewReader(testCase.payload)
	req, err := http.NewRequest(testCase.method, w.BaseURL+testCase.endpoint, payload)
	if err != nil {
		return err
	}
	req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	res, err := w.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	got := prepareData(data)

	if got != testCase.want {
		return fmt.Errorf("not equal: got = %v, want = %v", got, testCase.want)
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

func prepareMessage(report map[string]string) string {
	message := "Test API worker. Failed tests:\n"

	var i int = 1
	for name, err := range report {
		message += fmt.Sprintf("%v. %v: %v\n", i, name, err)
		i++
	}

	return message
}
