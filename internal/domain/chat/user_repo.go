package chat

type UserRepository interface {
	CreateUser(user *User) error
	GetUserById(userId string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	DeleteUserByUsername(username string) error
	DeleteUserById(userId string) error
	UpdateUser(user *User) error
	UserExistsById(userId string) (bool, error)
	UserExistsByUsername(username string) (bool, error)
	GetUserByIdUnscoped(userID string) (*User, error)
	RestoreUser(userID string) error
}
