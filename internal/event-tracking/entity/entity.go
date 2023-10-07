package entity

import (
	"time"
)

type UserActivity[T any] struct {
	UserID    string    `json:"userId,omitempty"`
	EventName string    `json:"eventName"`
	Data      T         `json:"data"`
	Device    string    `json:"device,omitempty"` // desktop, mobile, tablet
	OS        string    `json:"os,omitempty"`     // Windows, Android, iOS etc.
	App       string    `json:"app,omitempty"`    // play
	Timestamp time.Time `json:"timestamp"`
}

type CreateUserActivityDTO[T any] struct {
	EventName string `json:"eventName"`
	Data      T      `json:"data,omitempty"`
	App       string `json:"app,omitempty"`
}

type ActOnItemData struct {
	ItemID string `json:"itemId"`
}
