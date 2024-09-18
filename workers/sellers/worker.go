package sellers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/rotisserie/eris"
	"io"
	"log"
	"net/http"
	"sync"
)

type Competitor struct {
	Name string
	URL  string
}

type Worker struct {
	DatabaseEnv string            `json:"dbenv"`
	TaskCrons   map[string]string `json:"task_crons"`
}

func (worker *Worker) Init(ctx context.Context, conf config.StringMap) error {
	worker.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local")
	err := bcdb.InitDB(worker.DatabaseEnv)
	if err != nil {
		return eris.Wrapf(err, "Failed to initialize DB")
	}
	return nil
}

func (worker *Worker) Do(ctx context.Context) error {
	db := bcdb.DB()

	competitors, _ := fetchCompetitors(ctx, db)
	const numWorkers = 5
	var wg sync.WaitGroup
	jobs := make(chan Competitor, len(competitors))

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go request(i, jobs, &wg)
	}

	for _, competitor := range competitors {
		jobs <- competitor
	}
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()

	fmt.Println("All jobs completed")
	return nil
}

func fetchCompetitors(ctx context.Context, db *sqlx.DB) ([]Competitor, error) {

	query := `SELECT competitor_name, url FROM competitors;`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var competitors []Competitor
	for rows.Next() {
		var comp Competitor
		err := rows.Scan(&comp.Name, &comp.URL)
		if err != nil {
			return nil, err
		}
		competitors = append(competitors, comp)
	}

	return competitors, nil
}

func request(id int, jobs <-chan Competitor, results chan<- map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for competitor := range jobs {
		resp, err := http.Get(competitor.URL)
		if err != nil {
			log.Printf("Worker %d: Failed to get %s (%s): %v\n", id, competitor.Name, competitor.URL, err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Worker %d: Failed to read response for %s (%s): %v\n", id, competitor.Name, competitor.URL, err)
			continue
		}
		var jsonData interface{}
		if err := json.Unmarshal(body, &jsonData); err != nil {
			log.Printf("Worker %d: Failed to parse JSON for %s (%s): %v\n", id, competitor.Name, competitor.URL, err)
			continue
		}

		result := map[string]interface{}{
			competitor.Name: jsonData,
		}

		results <- result
	}
}

func (worker *Worker) GetSleep() int {
	return 0
}
