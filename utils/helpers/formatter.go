package helpers

import (
	"fmt"
	"strconv"
	"strings"
)

type FormatValues struct{}

func formatWithCommas(num interface{}) string {
	var str string

	switch v := num.(type) {
	case int:
		str = strconv.Itoa(v)
	case float64:
		str = fmt.Sprintf("%.2f", v) // Format float to 2 decimal places
	default:
		return "Invalid type"
	}

	if strings.Contains(str, ".") {
		parts := strings.Split(str, ".")
		intPart := parts[0]
		decimalPart := parts[1]

		return addCommas(intPart) + "." + decimalPart
	}

	return addCommas(str)
}

func addCommas(numStr string) string {
	n := len(numStr)
	if n <= 3 {
		return numStr
	}

	var result strings.Builder
	for i, digit := range numStr {
		if (n-i)%3 == 0 && i != 0 {
			result.WriteString(",")
		}
		result.WriteRune(digit)
	}

	return result.String()
}

func (f *FormatValues) CPM(cpm float64) string {
	return fmt.Sprintf("$%.2f", cpm)
}

func (f *FormatValues) RPM(rpm float64) string {
	return fmt.Sprintf("$%.2f", rpm)
}

func (f *FormatValues) DPRPM(dpRpm float64) string {
	return fmt.Sprintf("$%.2f", dpRpm)
}

func (f *FormatValues) GP(gp float64) string {
	return fmt.Sprintf("$%.2f", gp)
}

func (f *FormatValues) GPP(gpp float64) string {
	return fmt.Sprintf("%.f%%", gpp*100)
}

func (f *FormatValues) PubImps(pubImps int) string {
	return formatWithCommas(pubImps)
}

func (f *FormatValues) SoldImps(soldImps int) string {
	return formatWithCommas(soldImps)
}

func (f *FormatValues) Revenue(revenue float64) string {
	return fmt.Sprintf("$%.2f", revenue)
}

func (f *FormatValues) Cost(cost float64) string {
	return fmt.Sprintf("$%.2f", cost)
}

func (f *FormatValues) FillRate(fillRate float64) string {
	return fmt.Sprintf("%.2f%%", fillRate*100)
}

func (f *FormatValues) BidRequests(bidRequests float64) string {
	return formatWithCommas(bidRequests)
}
func (f *FormatValues) BidResponses(bidResponses float64) string {
	return formatWithCommas(bidResponses)
}
