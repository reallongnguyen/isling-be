package entity

import "github.com/mileusna/useragent"

type UserAgent struct {
	Device string
	OS     string
}

func (r *UserAgent) From(ua useragent.UserAgent) *UserAgent {
	r.Device = ua.Device
	r.OS = ua.OS

	return r
}
