package constant

type CompetitorUpdateRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"  validate:"required,url"`
}
