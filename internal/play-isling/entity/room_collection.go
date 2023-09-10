package entity

type (
	RoomCollection struct {
		ID    string        `json:"id"`
		Name  string        `json:"name"`
		Rooms []*RoomPublic `json:"rooms"`
	}
)
