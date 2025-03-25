package bid_cache_report

type BidCacheData struct {
	Time                 string  `boil:"time" json:"time" toml:"time" yaml:"time"`
	PublisherID          string  `boil:"publisher_id" json:"publisher_id" toml:"publisher_id" yaml:"publisher_id"`
	Domain               string  `boil:"domain" json:"domain" toml:"domain" yaml:"domain"`
	Target               string  `boil:"tgrp" json:"tgrp" toml:"tgrp" yaml:"tgrp"`
	Revenue              float64 `boil:"revenue" json:"revenue" toml:"revenue" yaml:"revenue"`
	Cost                 float64 `boil:"cost" json:"cost" toml:"cost" yaml:"cost"`
	DemandPartnerFee     float64 `boil:"demand_partner_fee" json:"demand_partner_fee" toml:"demand_partner_fee" yaml:"demand_partner_fee"`
	SoldImpressions      int     `boil:"sold_impressions" json:"sold_impressions" toml:"sold_impressions" yaml:"sold_impressions"`
	PublisherImpressions int     `boil:"publisher_impressions" json:"publisher_impressions" toml:"publisher_impressions" yaml:"publisher_impressions"`
	DataFee              float64 `boil:"data_fee" json:"data_fee" toml:"data_fee" yaml:"data_fee"`
	GP                   float64 `boil:"gp" json:"gp" toml:"gp" yaml:"gp"`
	GPperPubImp          float64 `boil:"gp_per_pub_imp" json:"gp_per_pub_imp" toml:"gp_per_pub_imp" yaml:"gp_per_pub_imp"`
}
