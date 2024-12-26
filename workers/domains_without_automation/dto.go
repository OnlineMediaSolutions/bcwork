package domains_without_automation

type Token struct {
	Data Data `json:"data"`
}

type Data struct {
	Token string `boil:"token" json:"token" toml:"token" yaml:"token"`
}

type TokenRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
