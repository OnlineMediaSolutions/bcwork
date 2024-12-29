package dto

type PublisherDomainRequest struct {
	DemandParnerId string `json:"demand_partner_id"`
	Data           []Data `json:"data"`
}

type Data struct {
	PubId        string `json:"pubId"`
	Domain       string `json:"domain"`
	AdsTxtStatus bool   `json:"ads_txt_status"`
}
