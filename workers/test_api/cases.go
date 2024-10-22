package testapi

import (
	"net/http"
)

const testPublisherID = "9995"

type testCase struct {
	name     string
	endpoint string
	method   string
	payload  string
	want     string
}

var testCases []testCase = []testCase{
	{
		name:     "Test",
		endpoint: "/",
		method:   http.MethodGet,
		payload:  "",
		want:     "UP",
	},
	{
		name:     "TestPing",
		endpoint: "/ping",
		method:   http.MethodGet,
		payload:  "",
		want:     `{"status":"OK","message":"Service is UP!!!"}`,
	},
	{
		name:     "TestPublisherGet",
		endpoint: "/publisher/get",
		method:   http.MethodPost,
		payload:  `{"filter":{"publisher_id":["` + testPublisherID + `"]}}`,
		want:     `[{"publisher_id":"9995","name":"Test_New_API","office_location":"IL","domains":["_test.com"],"integration_type":["JS Tags (Compass)","Prebid.js","Prebid Server","oRTB EP","JS Tags (NP)"],"status":"Paused","confiant":{},"pixalate":{}}]`,
	},
	{
		name:     "TestPublisherDomainGet",
		endpoint: "/publisher/domain/get",
		method:   http.MethodPost,
		payload:  `{"filter":{"publisher_id":["` + testPublisherID + `"]}}`,
		want:     `[{"publisher_id":"9995","domain":"_test.com","automation":false,"gpp_target":0,"integration_type":[],"confiant":{},"pixalate":{}}]`,
	},
	{
		name:     "TestFactorGet",
		endpoint: "/factor/get",
		method:   http.MethodPost,
		payload:  `{"filter":{"publisher":["` + testPublisherID + `"]}}`,
		want:     `[{"publisher":"9995","domain":"_test.com","country":null,"device":"desktop","factor":1}]`,
	},
	{
		name:     "TestFloorGet",
		endpoint: "/floor/get",
		method:   http.MethodPost,
		payload:  `{"filter":{"publisher":["` + testPublisherID + `"]}}`,
		want:     `[{"rule_id":"_test_rule_id","publisher":"9995","publisher_name":"Test_New_API","domain":"_test.com","country":"all","device":"mobile","floor":0.2,"browser":"","os":"","placement_type":""}]`,
	},
	{
		name:     "TestDPOGet",
		endpoint: "/dpo/get",
		method:   http.MethodPost,
		payload:  `{"filter":{"publisher":["` + testPublisherID + `"]}}`,
		want:     `[{"rule_id":"_test_rule_id","demand_partner_id":"_test","publisher":"9995","domain":"_test.com","country":"il","os":"","device_type":"","placement_type":"","browser":"","factor":90,"name":"Test_New_API","demand_partner_name":"_test"}]`,
	},
	{
		name:     "TestGlobalFactorGet",
		endpoint: "/global/factor/get",
		method:   http.MethodPost,
		payload:  `{"filter":{"key":["consultant_fee"],"publisher_id":["` + testPublisherID + `"]}}`,
		want:     `[{"key":"consultant_fee","publisher_id":"9995","value":0.11}]`,
	},
}

// app.Get("/report/daily/revenue", report.DailyRevenue)
// app.Get("/report/hourly/revenue", report.HourlyRevenue)
// app.Get("/report/demand", rest.DemandReportGetHandler)
// app.Get("/report/demand/hourly", rest.DemandHourlyReportGetHandler)
// app.Get("/report/publisher", rest.PublisherReportGetHandler)
// app.Get("/report/publisher/hourly", rest.PublisherHourlyReportGetHandler)
// app.Get("/report/iiq/hourly", rest.IiqTestingGetHandler)

// app.Post("/metadata/update", rest.MetadataPostHandler)

// app.Get("/had/price/set", rest.HouseAdPriceSetHandler)
// app.Get("/had/price/get", rest.HouseAdPriceGetHandler)
// app.Get("/had/price/get/all", rest.HouseAdPriceGetAllHandler)

// app.Post("/demand/factor", rest.DemandFactorPostHandler)
// app.Get("/demand/factor/set", rest.DemandFactorSetHandler)
// app.Get("/demand/factor/get", rest.DemandFactorGetHandler)
// app.Get("/demand/factor/get/all", rest.DemandFactorGetAllHandler)

// app.Get("/price/floor/set", rest.PriceFloorSetHandler)
// app.Get("/price/floor/get", rest.PriceFloorGetHandler)
// app.Get("/price/floor/get/all", rest.PriceFloorGetAllHandler)
// app.Post("/price/fixed", rest.FixedPricePostHandler)
// app.Get("/price/fixed", rest.FixedPriceGetAllHandler)

// app.Post("/confiant", validations.ValidateConfiant, rest.ConfiantPostHandler)
// app.Post("/confiant/get", rest.ConfiantGetAllHandler)

// app.Post("/global/factor", validations.ValidateGlobalFactor, rest.GlobalFactorPostHandler)

// app.Post("/pixalate", validations.ValidatePixalate, rest.PixalatePostHandler)
// app.Post("/pixalate/get", rest.PixalateGetAllHandler)
// app.Delete("/pixalate/delete", rest.PixalateDeleteHandler)

// app.Post("/block", rest.BlockPostHandler)
// app.Post("/block/get", rest.BlockGetAllHandler)
// app.Post("/dp/get", rest.DemandPartnerGetHandler)

// app.Get("/dpo/update", dpo.ValidateQueryParams, rest.DemandPartnerOptimizationUpdateHandler)

// app.Post("/publisher/new", validations.PublisherValidation, rest.PublisherNewHandler)
// app.Post("/publisher/update", rest.PublisherUpdateHandler)
// app.Post("/publisher/count", rest.PublisherCountHandler)
// app.Post("/publisher/details/get", rest.PublisherDetailsGetHandler)

// app.Post("/publisher/domain", validations.PublisherDomainValidation, rest.PublisherDomainPostHandler)

// app.Post("/factor", validations.ValidateFactor, rest.FactorPostHandler)

// app.Post("/floor", validations.ValidateFloors, rest.FloorPostHandler)

// app.Post("/bulk/factor", validations.ValidateBulkFactors, bulk.FactorBulkPostHandler)
// app.Post("/bulk/floor", validations.ValidateBulkFloor, bulk.FloorBulkPostHandler)
// app.Post("/bulk/dpo", validations.ValidateDPOInBulk, bulk.DemandPartnerOptimizationBulkPostHandler)
// app.Post("/bulk/global/factor", validations.ValidateBulkGlobalFactor, bulk.GlobalFactorBulkPostHandler)

// app.Post("/config/get", rest.ConfigurationGetHandler)
// app.Post("/config", validations.ValidateConfig, rest.ConfigurationPostHandler)
