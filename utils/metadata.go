package utils

type MetadataKey struct {
	Publisher string `json:"publisher"`
	Domain    string `json:"domain"`
	Device    string `json:"device"`
	Country   string `json:"country"`
}

func CreateMetadataKey(data MetadataKey, prefix string) string {
	key := prefix + ":" + data.Publisher
	if data.Domain != "" {
		key = key + ":" + data.Domain
	}
	if data.Country != "" {
		key = key + ":" + data.Country
	}
	if data.Device == "mobile" {
		key = "mobile:" + key
	}
	return key
}
