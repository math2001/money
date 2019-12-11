package db

// UserDB is the folder containing all the user's data.
// Everything in their is encrypted using his password
type User struct {
	privroot string
    cryptor: cyptor
}

func (u *User) Save(filename string, plaintext []byte) error {
	panic("not implemented")
}

func (u *User) Load(filename string) ([]byte, error) {
	panic("not implemented")
}

func NewUser(privroot string) *User {
	return &User{
		privroot: privroot,
	}
}
