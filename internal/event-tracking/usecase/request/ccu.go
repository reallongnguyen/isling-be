package request

import "time"

type CountCCUReq struct {
	Timestamp  time.Time
	WindowSize uint
}
