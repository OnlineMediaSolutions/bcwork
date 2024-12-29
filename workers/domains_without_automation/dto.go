package domains_without_automation

type Token struct {
	Data Data `json:"data"`
}

type Data struct {
	Token string `boil:"token" json:"token" toml:"token" yaml:"token"`
}

type ReportResults struct {
	Data ResultsData `json:"data"`
}

type ResultsData struct {
	Result []Result `json:"result"`
}

type Result struct {
	Publisher      string  `json:"publisher"`
	PublisherId    int     `json:"PublisherId"`
	Domain         string  `json:"domain"`
	AccountManager string  `json:"AM"`
	PubImps        int     `json:"PubImps"`
	LoopingRatio   float64 `json:"nbLR"`
	Cost           float64 `json:"Cost"`
	CPM            float64 `json:"nbCpm"`
	Revenue        float64 `json:"Revenue"`
	RPM            float64 `json:"nbRpm"`
	DpRPM          float64 `json:"nbDpRpm"`
	GP             float64 `json:"nbGp"`
	GPP            float64 `json:"nbGpp"`
}

type EmailProperties struct {
	TO   string `json:"TO"`
	BCC  string `json:"BCC"`
	FROM string `json:"FROM"`
}
