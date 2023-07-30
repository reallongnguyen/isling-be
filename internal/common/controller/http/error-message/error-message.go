package errormessage

import (
	play_entity "isling-be/internal/play-isling/entity"
	"net/http"
)

type HTTPError struct {
	HTTPCode int    `json:"-"`
	Code     int    `json:"code"`
	Message  string `json:"message"`
}

var ErrorMap = make(map[error]*HTTPError)

func registerErr(err error, httpCode, code int) {
	ErrorMap[err] = &HTTPError{
		HTTPCode: httpCode,
		Code:     code,
		Message:  err.Error(),
	}
}

func init() {
	registerErr(play_entity.ErrInvalidRoomSlug, http.StatusBadRequest, 1020000)
	registerErr(play_entity.ErrMissingUpdatePerm, http.StatusForbidden, 1020001)
	registerErr(play_entity.ErrMissingDeletePerm, http.StatusForbidden, 1020002)
	registerErr(play_entity.ErrRoomNotFound, http.StatusNotFound, 1020003)
}
