package dto

type BlockUpdateRequest struct {
	Publisher string   `json:"publisher" validate:"required"`
	Domain    string   `json:"domain"`
	BCAT      []string `json:"bcat"`
	BADV      []string `json:"badv"`
}

type BlockGetRequest struct {
	Types     []string `json:"types"`
	Publisher string   `json:"publisher"`
	Domain    string   `json:"domain"`
}
