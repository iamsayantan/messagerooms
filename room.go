package messagerooms

// Room represents a single messaging room.
type Room struct {
	ID        string `json:"id"`
	RoomName  string `json:"room_name"`
	UserID    string `json:"-"`
	CreatedBy *User  `json:"created_by" gorm:"foreignkey:UserID"`
	Users     []User `json:"users,omitempty" gorm:"many2many:room_users"`
}

// RoomRepository provides interface methods for interacting with rooms data store.
type RoomRepository interface {
	Create(name string, user User) (*Room, error)
	Find(id string) (*Room, error)
	FindAll() ([]*Room, error)
	AddUserToRoom(room Room, user User) error
	CheckUserExistsInRoom(room Room, user User) bool
}
