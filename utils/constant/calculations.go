package constant

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
