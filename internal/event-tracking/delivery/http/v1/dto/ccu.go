package dto

import (
	"isling-be/internal/event-tracking/usecase/request"
	"strconv"
	"time"
)

type CountCCUReq struct {
	Timestamp  string `validate:"required,datetime"`
	WindowSize string `validate:"number,max=2"`
}

func (r *CountCCUReq) ToReq() (*request.CountCCUReq, error) {
	timestamp, err := time.Parse(time.RFC3339, r.Timestamp)
	if err != nil {
		return nil, err
	}

	windowSize, err := strconv.Atoi(r.WindowSize)
	if err != nil {
		return nil, err
	}

	req := &request.CountCCUReq{
		Timestamp:  timestamp,
		WindowSize: uint(windowSize),
	}

	return req, nil
}
