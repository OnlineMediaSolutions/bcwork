package rtb_house_report

import "time"

func PrepareRequestData() RequestData {
	requestData := RequestData{
		Data: RequestDetails{
			Date: Date{
				Range: []string{
					time.Now().AddDate(0, 0, -3).Format("2006-01-02 15:04:05"),
					time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05"),
				},
				Interval: "none",
			},
			Dimensions: []string{
				"Publisher",
				"Domain",
			},
			Metrics: []string{
				"SoldImps",
			},
		},
	}

	return requestData
}
