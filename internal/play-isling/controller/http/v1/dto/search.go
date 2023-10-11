package dto

import (
	"isling-be/internal/play-isling/usecase"
	"strconv"
)

type SearchReqDTO struct {
	Query  string `json:"query" validate:"required,min=1,max=64"`
	Limit  string `json:"limit,omitempty" validate:"omitempty,number,max=2"`
	Offset string `json:"offset,omitempty" validate:"omitempty,number"`
}

func (r *SearchReqDTO) ToReq() *usecase.SearchRequest {
	limit, _ := strconv.Atoi(r.Limit)
	offset, _ := strconv.Atoi(r.Offset)

	return &usecase.SearchRequest{
		Query:  r.Query,
		Limit:  limit,
		Offset: offset,
	}
}
