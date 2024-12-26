package domains_without_automation

import (
	"bytes"
	"context"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bccron"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"strings"
)

//
// This worker fetch once a week all domains without factor automation and sends an email
//

const fetchTokenUrl = "http://10.166.10.36:8080/api/auth/token"

type Worker struct {
	DatabaseEnv string `json:"dbenv"`
	Cron        string `json:"cron"`
}

type TransactionIds struct {
	TransactionId string `json:"transaction_id"`
}

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {

	worker.DatabaseEnv = conf.GetStringValueWithDefault(config.DBEnvKey, "local")
	worker.Cron, _ = conf.GetStringValue("cron")

	err := bcdb.InitDB(worker.DatabaseEnv)
	if err != nil {
		return eris.Wrapf(err, "failed to initalize DB")
	}
	return nil
}

func (w *Worker) Do(ctx context.Context) error {
	log.Info().Msg("Starting domains_without_automation worker")

	domains, err := fetchDomainsWithoutAutomation(ctx)
	if err != nil {
		return err
	}

	token, err := fetchToken(ctx)
	print(token)
	//fetchPubImps()
	//calculatePubImpsPerDomain()
	//sendEmail()

	print(domains)

	return nil
}

func fetchToken(ctx context.Context) (string, error) {
	tokenRequest := TokenRequest{
		Login:    "",
		Password: "",
	}

	publisherData, statusCode, err := worker.HttpClient.Do(ctx, http.MethodPost, fetchTokenUrl, bytes.NewBuffer(body))

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d", statusCode)
	}

	var token Token
	print(token)
	//if err := json.Unmarshal(publisherData, &publishers); err != nil {
	//	return nil, fmt.Errorf("error parsing publisher data  from API")
	//}
	return "", nil
}

func wrapTransactions(transactions []*TransactionIds) string {
	var wrappedTransactionIds []string
	for _, transactionId := range transactions {
		wrappedTransactionIds = append(wrappedTransactionIds, fmt.Sprintf(`'%s'`, transactionId.TransactionId))
	}
	transactionsIds := strings.Join(wrappedTransactionIds, ",")
	return transactionsIds
}

func fetchDomainsWithoutAutomation(ctx context.Context) ([]string, error) {

	var domains []string
	results, err := models.PublisherDomains(models.PublisherDomainWhere.Automation.EQ(false),
		qm.Select(models.PublisherDomainColumns.Domain)).All(ctx, bcdb.DB())

	if err != nil {
		return nil, eris.Wrapf(err, "failed to fetch domains without Automation from DB")
	}

	for _, result := range results {
		domains = append(domains, result.Domain)
	}
	return domains, nil
}

func (worker *Worker) GetSleep() int {
	if worker.Cron != "" {
		return bccron.Next(worker.Cron)
	}
	return 0
}
