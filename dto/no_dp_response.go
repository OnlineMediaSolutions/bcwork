package dto

import "github.com/m6yf/bcwork/models"

type NoDPResponseReport struct {
	Time          string  `boil:"time" json:"time"`
	DPID          string  `boil:"dpid" json:"dpid"`
	PubID         string  `boil:"pubid" json:"pubid"`
	PublisherName string  `boil:"publisher_name" json:"publisher_name"`
	Domain        string  `boil:"domain" json:"domain"`
	BidRequests   float64 `boil:"bid_requests" json:"bid_requests"`
	BidResponses  float64 `boil:"bid_responses" json:"bid_responses"`
}

func (n NoDPResponseReport) ToModel() *models.NoDPResponseReport {
	return &models.NoDPResponseReport{
		Time:            n.Time,
		DemandPartnerID: n.DPID,
		PublisherID:     n.PubID,
		Domain:          n.Domain,
		BidRequests:     n.BidRequests,
		BidResponses:    n.BidResponses,
	}
}

func (n *NoDPResponseReport) FromModel(mod *models.NoDPResponseReport) {
	n.Time = mod.Time
	n.DPID = mod.DemandPartnerID
	n.PubID = mod.PublisherID
	if mod.R != nil && mod.R.Publisher != nil {
		n.PublisherName = mod.R.Publisher.Name
	}
	n.Domain = mod.Domain
	n.BidRequests = mod.BidRequests
	n.BidResponses = mod.BidResponses
}

func (n NoDPResponseReport) BuildKey() string {
	return n.Time + ":" + n.DPID + ":" + n.PubID + ":" + n.Domain
}
