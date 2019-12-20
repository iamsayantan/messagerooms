package messagerooms

// Room represents a single messaging room.
type Room struct {
	ID        string
	RoomName  string
	UserID    string
	CreatedBy User
	Users     []User
}

// RoomRepository provides interface methods for interacting with rooms data store.
type RoomRepository interface {
	Create(name string, user User) (*Room, error)
	Find(id string) (*Room, error)
	AddUserToRoom(room Room, user User) error
	CheckUserExistsInRoom(room Room, user User) bool
}
