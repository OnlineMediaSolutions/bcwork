package metadata_clean

import (
	"bufio"
	"context"
	"fmt"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"os"
	"strings"
)

type Worker struct {
	DatabaseEnv string `json:"dbenv"`
}

type TransactionIds struct {
	TransactionId string `json:"transaction_id"`
}

const fetch_query = `WITH ranked_records AS (
    select transaction_id, key, updated_at,
        row_number() over (Partition by key Order by updated_at) AS row_num,
        count(*) over (partition by key) as total_count
    from metadata_queue)
select transaction_id
from
    ranked_records
where
    row_num <=5 and total_count >= 5
ORDER BY
    key,
    updated_at;`

const copy_query = `
insert into metadata_queue_temp
	select * from metadata_queue
where transaction_id in (%s);`

const delete_query = `DELETE from metadata_queue_temp where transaction_id in (%s);`

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {

	w.DatabaseEnv = conf.GetStringValueWithDefault(config.DBEnvKey, "local")
	err := bcdb.InitDB(w.DatabaseEnv)
	if err != nil {
		return eris.Wrapf(err, "failed to initalize DB")
	}
	return nil
}

func (w *Worker) Do(ctx context.Context) error {
	fmt.Println("Start to Remove old rows from Metadata_queue")
	err := startWorker()
	if err != nil {
		return err
	}

	transactions, err := fetchRowsFromDB(ctx)
	if err != nil {
		return err
	}

	transactionsIds := wrapTransactions(transactions)

	err = copyTransactionToTempTable(ctx, transactionsIds)
	if err != nil {
		return err
	}

	err = deleteTransactionsInMetaData(ctx, transactionsIds)
	if err != nil {
		return err
	}
	fmt.Println("Finished Metadata_queue clean Worker")
	return nil
}

func wrapTransactions(transactions []*TransactionIds) string {
	var wrappedTransactionIds []string
	for _, transactionId := range transactions {
		wrappedTransactionIds = append(wrappedTransactionIds, fmt.Sprintf(`'%s'`, transactionId.TransactionId))
	}
	transactionsIds := strings.Join(wrappedTransactionIds, ",")
	return transactionsIds
}

func deleteTransactionsInMetaData(ctx context.Context, transactionIds string) error {

	delete_query := fmt.Sprintf(delete_query, transactionIds)
	_, err := queries.Raw(delete_query).ExecContext(ctx, bcdb.DB())
	if err != nil {
		return fmt.Errorf("error deleting data from metadata_queue table: %w", err)
	}
	return nil
}

func copyTransactionToTempTable(ctx context.Context, transactionIds string) error {

	copy_query := fmt.Sprintf(copy_query, transactionIds)
	_, err := queries.Raw(copy_query).ExecContext(ctx, bcdb.DB())
	if err != nil {
		return fmt.Errorf("error copy data to metadata_queue_temp table: %w", err)
	}
	return nil
}

func fetchRowsFromDB(ctx context.Context) ([]*TransactionIds, error) {

	var transactionIds []*TransactionIds
	err := queries.Raw(fetch_query).Bind(ctx, bcdb.DB(), &transactionIds)
	if err != nil {
		return nil, fmt.Errorf("error fetching transaction ids from metadataQueue: %w", err)
	}
	return transactionIds, nil
}

func startWorker() error {
	reader := bufio.NewReader(os.Stdin)

	log.Info().Msg("Are you sure that you want to remove old rows from Metadata_queue?")
	log.Info().Msg("Press Y to continue")

	var input, _ = reader.ReadString('\n')
	choice := strings.TrimSpace(input)

	if choice != "Y" && choice != "y" {
		return fmt.Errorf("\"Worker stopped")
	}
	return nil
}

func (w *Worker) GetSleep() int {
	return 0
}
