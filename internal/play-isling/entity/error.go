package entity

import "errors"

var (
	ErrInvalidRoomSlug   = errors.New("invalid room slug")
	ErrMissingDeletePerm = errors.New("missing delete permission")
	ErrMissingUpdatePerm = errors.New("missing update permission")
	ErrRoomNotFound      = errors.New("room not found")
)
