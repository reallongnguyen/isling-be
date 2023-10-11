package entity

type RoomOwnerSearchResult struct {
	ID        string  `json:"id"`
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
	FullName  *string `json:"fullName,omitempty"`
	AvatarURL *string `json:"avatarUrl,omitempty"`
}

type RoomSearchResult struct {
	RoomPublic
	ID        string                 `json:"id"`
	OwnerID   string                 `json:"ownerId"`
	Owner     *RoomOwnerSearchResult `json:"owner"`
	Highlight string                 `json:"highlight,omitempty"`
	Score     float32                `json:"score"`
}
