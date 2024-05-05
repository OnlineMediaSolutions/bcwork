package influx

import (
	"context"
	"github.com/friendsofgo/errors"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

//func Query(ctx context.Context, query string) (*api.QueryTableResult, error) {
//	token := "0b8MubUWHEDc3ykmjIQt7tEpp4B3QwPyV_rkwDo3nF5epe5UGqA1XZHy4SP7aHry6cXZsA1jwcPcvHnv4dK88w=="
//	//if token == "" {
//	//	log.Error().Msgf("influx db otken not found please set ENV INFLUX_TOKEN")
//	//	return nil, nil
//	//}
//
//	org := "brightcom.com"
//	client := influxdb2.NewClientWithOptions(
//		"https://us-east-1-1.aws.cloud2.influxdata.com",
//		token,
//		influxdb2.DefaultOptions().SetBatchSize(20).SetHTTPRequestTimeout(180))
//	// Get query client
//	queryAPI := client.QueryAPI(org)
//
//	// Get parser flux query result
//	result, err := queryAPI.Query(ctx, query)
//	if err != nil {
//		return nil, errors.Wrapf(err, "failed to query influx")
//	}
//
//	return result, nil
//}

func Query(ctx context.Context, query string) (*api.QueryTableResult, error) {

	//token := "UhGpkR3ZJ4tL5jO_Cbi9OwsLko2F4t-LTJjucPVvC4O8aDPgENeYOF-PnBHJjw29Qy73O_McoL3A1RuiCnpSAg=="
	token := "Gt07Fb1159iEdr0ZPx4gh6FaKIpCwhrJBIxj-fhGoCykeQFRc5-Pr0irzBPNJCFiJ_RjmzpAbfkcfqGh-gZ9GA=="
	//if token == "" {
	//	log.Error().Msgf("influx db otken not found please set ENV INFLUX_TOKEN")
	//	return nil, nil
	//}

	org := "brightcom.com"
	client := influxdb2.NewClientWithOptions(
		"http://bcinflux-nyc1-02-private:8086",
		token,
		influxdb2.DefaultOptions().SetBatchSize(20).SetHTTPRequestTimeout(500))
	// Get query client
	queryAPI := client.QueryAPI(org)

	// Get parser flux query result
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query influx")
	}

	return result, nil
}
