package db

// UserDB is the folder containing all the user's data.
// Everything in their is encrypted using his password
type User struct {
	// root is the user's own folder (see api.go)
	*Store
	ID    int
	Email string
}

func NewUser(id int, email, root string) *User {
	return &User{
		Email: email,
		ID:    id,
		Store: NewStore(root),
	}
}
