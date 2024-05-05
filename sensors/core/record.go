package core

type SensorRecord struct {
	Key    string            `json:"key"`
	Tags   map[string]string `json:"tags,omitempty"`
	Values Counters          `json:"values"`
}

func NewRecord(k string, v Counters) *SensorRecord {
	return &SensorRecord{
		Key:    ExtractKey(k),
		Tags:   ExtractTags(k),
		Values: v,
	}
}

type SensorsTable struct {
	Hour   string          `json:"hour"`
	Closed bool            `json:"closed"`
	Data   []*SensorRecord `json:"data"`
}

func SensorTableFromFile(sensors *HourlySensors) SensorsTable {
	table := SensorsTable{
		Closed: sensors.Closed,
		Hour:   sensors.Hour,
	}
	for k, v := range sensors.Sensors {
		table.Data = append(table.Data, NewRecord(k, v))
	}

	return table
}

func (st SensorsTable) Count(key string) int64 {
	var res int64
	//for _, sensor := range st.Data {
	//	if sensor.Key == key {
	//		res += sensor.Values.Count
	//	}
	//}

	return res
}

func (st SensorsTable) Select(key string, limit int) []*SensorRecord {
	res := make([]*SensorRecord, 0)
	for _, rec := range st.Data {
		if key == rec.Key {
			res = append(res, rec)
		}
		if len(res) == limit {
			break
		}
	}

	return res
}

func (st SensorsTable) SumCount(key string) float64 {
	var count float64
	for _, rec := range st.Data {
		if key == rec.Key {
			count += rec.Values["count"]
		}
	}

	return count
}
