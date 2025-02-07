package daily_alerts

type GPDecreaseReport struct{}

func (g *GPDecreaseReport) Aggregate(report Report) map[string]interface{} {
	return map[string]interface{}{}

}
