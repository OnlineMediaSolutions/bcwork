package constant

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
