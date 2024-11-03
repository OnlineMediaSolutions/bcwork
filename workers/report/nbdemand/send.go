package nbdemand

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/friendsofgo/errors"
	"github.com/m6yf/bcwork/core"
	"github.com/m6yf/bcwork/models"
	"github.com/m6yf/bcwork/utils/bcguid"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CompassNewDemandRecord struct {
	CombinationId          string  `json:"CombinationId"`
	DateStamp              int64   `json:"DateStamp"`
	DemandPartner          string  `json:"DemandPartner"`
	DemandPlacement        string  `json:"DemandPlacement"`
	IntegrationType        string  `json:"IntegrationType"`
	PublisherId            int64   `json:"PublisherId"`
	Domain                 string  `json:"Domain"`
	Country                string  `json:"Country"`
	Device                 string  `json:"Device"`
	Os                     string  `json:"Os"`
	Size                   string  `json:"Size,omitempty"`
	PaymentType            string  `json:"PaymentType,omitempty"`
	PublisherTagName       string  `json:"PublisherTagName,omitempty"`
	PublisherPlacementType string  `json:"PublisherPlacementType"`
	BCMBidRequests         int64   `json:"BCMBidRequests"`
	DPBidResponses         int64   `json:"DPBidResponses"`
	DPAvgBidPrice          float64 `json:"DPAvgBidPrice"`
	DPFee                  float64 `json:"DPFee"`
	DPAuctionWins          int64   `json:"DPAuctionWins"`
	Auction                float64 `json:"Auction"`
	SoldImps               int64   `json:"SoldImps"`
	Revenue                float64 `json:"Revenue"`
	UpdatedAt              string  `json:"updatedAt"`
	DataImpressions        int64   `json:"DataImpressions"`
	DataFee                float64 `json:"DataFee"`
}

var loc *time.Location

func ConvertToCompass(modSlice models.NBDemandHourlySlice) []*CompassNewDemandRecord {

	var err error
	if loc == nil {
		loc, err = time.LoadLocation("EST")
		if err != nil {
			log.Fatal().Err(err).Msg("failed to find EST loc")
		}
	}

	var res []*CompassNewDemandRecord
	for _, mod := range modSlice {
		if mod.BidResponses == 0 && mod.SoldImpressions == 0 {
			continue
		}
		compassDemandPartnerID := strings.ToLower(core.DemandPartnerMap[mod.DemandPartnerID])
		if compassDemandPartnerID == "" {
			log.Warn().Str("dpid", mod.DemandPartnerID).Msg("demand partner id not found")
			compassDemandPartnerID = "other"
		}

		if mod.DemandPartnerPlacementID == "" {
			mod.DemandPartnerPlacementID = "other"
		}
		val := &CompassNewDemandRecord{
			DemandPartner:          compassDemandPartnerID,
			DemandPlacement:        mod.DemandPartnerPlacementID,
			Domain:                 mod.Domain,
			Country:                mod.Country,
			IntegrationType:        "ortb",
			Os:                     mod.Os,
			PublisherPlacementType: mod.PlacementType,
			Device:                 mod.DeviceType,
			BCMBidRequests:         mod.BidRequests,
			DPBidResponses:         mod.BidResponses,
			Revenue:                mod.Revenue,
			SoldImps:               mod.SoldImpressions,
			DPFee:                  mod.DPFee,
			Auction:                mod.Auction,
			DPAuctionWins:          mod.AuctionWins,
			DPAvgBidPrice:          mod.AvgBidPrice * float64(mod.BidResponses),
			UpdatedAt:              time.Now().In(loc).Format("2006-01-02 15:04:00"),
			DataImpressions:        mod.DataImpressions,
			DataFee:                mod.DataFee,
		}

		//if false {
		val.Size = mod.Size
		if mod.RequestType != "js" {
			val.Size = "-"
		}

		val.PaymentType = "NP HB"
		if mod.RequestType == "js" {
			val.PaymentType = "NP CPM"
		} else if mod.RequestType == "tam" {
			val.PaymentType = "NP TAM"
		}

		val.PublisherTagName = "Header Bidding"
		if mod.RequestType == "js" {
			val.PublisherTagName = val.Domain + "_" + val.Device + "_" + val.Size
		}

		//}

		val.DateStamp = mod.Time.In(loc).Unix() / 100
		val.CombinationId = bcguid.NewFromf(mod.DemandPartnerID, mod.PublisherID, mod.DemandPartnerPlacementID, mod.Domain, mod.Country, mod.DeviceType, mod.Os, mod.PlacementType, mod.RequestType, mod.Size, mod.PaymentType, val.DateStamp, mod.Datacenter)
		val.PublisherId, err = strconv.ParseInt(mod.PublisherID, 10, 64)
		if err != nil {
			log.Error().Err(err).Interface("mod", mod).Msgf("illegal publisher id when parsing to int, 0 will be used")
		}

		res = append(res, val)
	}

	return res
}

func Send(vals []*CompassNewDemandRecord) error {

	b, err := json.Marshal(vals)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal compass values")
	}

	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)
	if err != nil {
		return errors.Wrapf(err, "failed to create gzip writer")
	}

	if _, err := g.Write(b); err != nil {
		return errors.Wrapf(err, "failed to write to gzip writer")
	}
	if err = g.Close(); err != nil {
		return errors.Wrapf(err, "failed to close gzip writer")
	}

	log.Info().Int("payload.records", len(vals)).Int("payload.bytes", len(b)).Int("payload.compressed", buf.Len()).Msg("sending to compass")
	req, err := http.NewRequest("POST", "https://nb-reports.ministerial5.com/demand-reports", &buf)
	//req, err := http.NewRequest("POST", "https://staging-nb-reports.ministerial5.com/supply-reports", &buf)
	req.Header.Add("Content-Encoding", "gzip")
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "network error")
	}

	defer resp.Body.Close()
	if err != nil {
		return errors.Wrapf(err, "failed to send post requests")
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "failed to read post body")
	}

	if resp.StatusCode != http.StatusOK {
		log.Error().Str("rec", string(b)).Msg("post request returned http error")
		return errors.Wrapf(err, "http error")
	}
	//i += pageSize
	//log.Info().Msgf("page %d sent", i/pageSize)
	//
	return nil

}
