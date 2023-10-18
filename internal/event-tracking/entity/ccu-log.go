package entity

import (
	"time"
)

type CCULog struct {
	UserID    string
	Device    string
	OS        string
	IP        string
	Timestamp time.Time
}

func (r *CCULog) GetKey() string {
	return r.UserID + "-" + r.IP + "-" + r.Device + "-" + r.OS
}

func (r *CCULog) GetUnixMin() int64 {
	return r.Timestamp.Unix() / 60
}
