package constant

import "github.com/m6yf/bcwork/utils/helpers"

type GPParameters struct {
	Cost             float64
	Revenue          float64
	DemandPartnerFee float64
	DataFee          float64
	BidRequests      float64
	PublisherID      string
	Fees             map[string]float64
	ConsultantFees   map[string]float64
}

type GPResults struct {
	TamFee        float64
	TechFee       float64
	ConsultantFee float64
	GP            float64
	GPP           float64
}

func PubFillRate(pubImps, bidRequests int64) float64 {
	if bidRequests == 0 {
		return 0
	}
	return float64(pubImps) / float64(bidRequests)
}

func CPM(cost, pubImps float64) float64 {
	if pubImps == 0 {
		return 0
	}
	return (cost / pubImps) * 1000
}

func RPM(revenue, pubImps float64) float64 {
	if pubImps == 0 {
		return 0
	}
	return (revenue / pubImps) * 1000
}

func DpRPM(revenue, soldImps float64) float64 {
	if soldImps == 0 {
		return 0
	}
	return (revenue / soldImps) * 1000
}

func GP(revenue, fees, consultantFees float64) float64 {
	return revenue - fees - consultantFees
}

func CalculateGP(params GPParameters) GPResults {
	results := GPResults{}

	// Calculate Tam Fee
	results.TamFee = helpers.RoundFloat(params.Fees["tam_fee"] * params.Cost)

	// Calculate Tech Fee
	results.TechFee = helpers.RoundFloat(params.Fees["tech_fee"] * params.BidRequests / 1000000)

	// Calculate Consultant Fee
	results.ConsultantFee = 0.0
	if value, exists := params.ConsultantFees[params.PublisherID]; exists {
		results.ConsultantFee = params.Cost * value
	}

	// Calculate Gross Profit (GP)
	results.GP = helpers.RoundFloat(params.Revenue - params.Cost - params.DemandPartnerFee - params.DataFee - results.TamFee - results.TechFee - results.ConsultantFee)

	// Calculate Gross Profit Percentage (GPP)
	results.GPP = 0
	if params.Revenue != 0 {
		results.GPP = helpers.RoundFloat(results.GP / params.Revenue)
	}

	return results
}
