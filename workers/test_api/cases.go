package testapi

type testCase struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	Method   string `json:"method"`
	Payload  string `json:"payload"`
	Want     string `json:"want"`
}
