package ip

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/bcdb"
	"github.com/m6yf/bcwork/config"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/quest"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"sort"
	"time"
)

type Worker struct {
	Sleep       time.Duration `json:"sleep"`
	Hours       int           `json:"hours"`
	Limit       int           `json:"limit"`
	Start       string        `json:"start"`
	DatabaseEnv string        `json:"dbenv"`
	Debug       bool          `json:"debug"`
	Domain      string        `json:"domain"`
}

func (w *Worker) Init(ctx context.Context, conf config.StringMap) error {

	var err error
	w.Sleep, _ = conf.GetDurationValueWithDefault("sleep", time.Duration(5*time.Minute))
	w.Hours, err = conf.GetIntValueWithDefault("hours", 1)
	if err != nil {
		log.Warn().Err(err).Msg("failed to fetch hours config value (will user default)")
	}

	w.Limit, err = conf.GetIntValueWithDefault("limit", 200)
	if err != nil {
		log.Warn().Err(err).Msg("failed to fetch limit config value (will user default)")
	}

	w.Domain = conf.GetStringValueWithDefault("domain", "postgresqltutorial.com")

	w.DatabaseEnv = conf.GetStringValueWithDefault("dbenv", "local_prod")
	err = bcdb.InitDB(w.DatabaseEnv)
	if err != nil {
		return errors.Wrapf(err, "failed to initalize DB")
	}

	err = quest.InitDB("quest2")
	if err != nil {
		return errors.Wrapf(err, "failed to initalize DB")
	}

	w.Start, _ = conf.GetStringValue("start")
	if w.Start != "" {
		w.Sleep = 0
	}

	w.Debug = conf.GetBoolValueWithDefault("debug", false)
	if w.Debug {
		boil.DebugMode = true
	}

	return nil

}

func (w *Worker) Do(ctx context.Context) error {

	var err error
	ips, err := w.fetchIPs(ctx)
	if err != nil {
		return err
	}

	b, err := json.Marshal(ips)
	if err != nil {
		return err
	}

	mdata := models.MetadataQueue{
		Key:           "throttle.ips",
		Value:         b,
		TransactionID: bcguid.NewFrom("throttle.ips" + time.Now().String()),
	}
	err = mdata.Insert(ctx, bcdb.DB(), boil.Infer())
	if err != nil {
		return err
	}

	log.Info().Msg("DONE")
	return nil
}

func (w *Worker) GetSleep() int {
	return int(w.Sleep.Seconds())
}

type record struct {
	IP   string `json:"ip"`
	Imps string `json:"imps"`
}

func (w *Worker) fetchIPs(ctx context.Context) ([]string, error) {

	log.Info().Msg("fetch ips")

	var records []record

	//	q := fmt.Sprintf(`select * from
	//(select ip,sum(1) imps from impression where timestamp >= dateadd('h', -1 * %d, now()) group by ip) sub
	//where imps>%d order by imps desc`, w.Hours, w.Limit)

	q := fmt.Sprintf(`select * from 
(select ip,sum(1) imps from impression where timestamp >= dateadd('h', -1 * %d, now()) and domain='%s' group by ip) sub 
where imps>=%d`, w.Hours, w.Domain, w.Limit)

	//log.Info().Str("q", q).Msg("processBidRequestCounters")
	err := queries.Raw(q).Bind(ctx, quest.DB(), &records)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query impressions from questdb")
	}

	log.Info().Interface("ips", records).Int("len", len(records)).Msg("IPS to block")

	res := []string{}
	for _, r := range records {
		res = append(res, r.IP)
	}

	sort.Strings(res)

	return res, nil
}
