package entity

type (
	RoomCollection struct {
		ID    int64         `json:"id"`
		Name  string        `json:"name"`
		Rooms []*RoomPublic `json:"rooms"`
	}
)
