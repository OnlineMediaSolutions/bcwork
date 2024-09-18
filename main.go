package main

import (
	"github.com/m6yf/bcwork/cmd"
	"github.com/m6yf/bcwork/structs"
	"github.com/m6yf/bcwork/workers/alerts"
	"github.com/m6yf/bcwork/workers/ansible/inventory"
	"github.com/m6yf/bcwork/workers/dns"
	factors_autmation "github.com/m6yf/bcwork/workers/factors/automation"
	factors_monitor "github.com/m6yf/bcwork/workers/factors/monitor"
	"github.com/m6yf/bcwork/workers/hello"
	"github.com/m6yf/bcwork/workers/metadata"
	"github.com/m6yf/bcwork/workers/questclean"
	"github.com/m6yf/bcwork/workers/report/demand"
	"github.com/m6yf/bcwork/workers/report/iiq"
	"github.com/m6yf/bcwork/workers/report/ip"
	"github.com/m6yf/bcwork/workers/report/logs"
	"github.com/m6yf/bcwork/workers/report/nbdemand"
	"github.com/m6yf/bcwork/workers/report/nbsupply"
	"github.com/m6yf/bcwork/workers/report/revenue"
	"github.com/m6yf/bcwork/workers/sellers"
	"github.com/m6yf/bcwork/workers/sync/publisher"
	"github.com/rs/zerolog/log"
	"strings"
)

var buildtime string
var githash string
var gittag string
var modelver string
var version = "v1.1.1"

func main() {

	register()

	// Model Version string
	if modelver != "" {
		if strings.Contains(modelver, " ") {
			toks := strings.Split(modelver, " ")
			log.Logger = log.With().Str("model.version", toks[1]).Logger()
		}
	}

	log.Info().Str("worker.version", gittag).Msg("worker starting up")

	cmd.Execute()

}

func register() {
	structs.RegsiterName("hello", hello.Worker{})
	structs.RegsiterName("nbsupply", nbsupply.Worker{})
	structs.RegsiterName("nbdemand", nbdemand.Worker{})
	structs.RegsiterName("revenue", revenue.Worker{})
	structs.RegsiterName("demand", demand.Worker{})
	structs.RegsiterName("iiq", iiq.Worker{})
	structs.RegsiterName("qdbclean", questclean.Worker{})
	structs.RegsiterName("report.iiq", iiq.Worker{})
	structs.RegsiterName("dns", dns.Worker{})
	structs.RegsiterName("logs", logs.Worker{})
	structs.RegsiterName("metadata", metadata.Worker{})
	structs.RegsiterName("inventory", inventory.Worker{})
	structs.RegsiterName("ip", ip.Worker{})
	structs.RegsiterName("sync.publisher", publisher.Worker{})
	structs.RegsiterName("factors", factors_autmation.Worker{})
	structs.RegsiterName("factors.monitor", factors_monitor.Worker{})
	structs.RegsiterName("alerts", alerts.Worker{})
	structs.RegsiterName("sellers", sellers.Worker{})
}
