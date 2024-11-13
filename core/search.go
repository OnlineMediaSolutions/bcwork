package core

import (
	"context"

	"github.com/m6yf/bcwork/dto"
)

type SearchService struct{}

func NewSearchService() *SearchService {
	return &SearchService{}
}

type SearchRequest struct {
	Query   string
	Section string
}

func (s *SearchService) Search(ctx context.Context, ops *SearchRequest) ([]*dto.SearchResult, error) {
	return nil, nil
}
