package users

type UserRepository interface {
	GetUserByID(userID int) (*User, error)
}
