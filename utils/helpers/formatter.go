package helpers

import (
	"fmt"
)

type FormatValues struct{}

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
	return fmt.Sprintf("$%.1f", gp)
}

func (f *FormatValues) GPP(gpp float64) string {
	return fmt.Sprintf("%.0f%%", gpp*100)
}

func (f *FormatValues) PubImps(pubImps float64) string {
	return fmt.Sprintf("$%f", pubImps)
}

func (f *FormatValues) SoldImps(soldImps float64) string {
	return fmt.Sprintf("$%f", soldImps)
}

func (f *FormatValues) Revenue(revenue float64) string {
	return fmt.Sprintf("$%.0f", revenue)
}

func (f *FormatValues) Cost(cost float64) string {
	return fmt.Sprintf("$%.0f", cost)
}

func (f *FormatValues) FillRate(fillRate float64) string {
	return fmt.Sprintf("%.2f%%", fillRate*100)
}
