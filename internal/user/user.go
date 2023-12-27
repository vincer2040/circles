package user

import "golang.org/x/crypto/bcrypt"

type User struct {
	First    string
	Last     string
	Email    string
	Password string
}

func New(first, last, email, password string) User {
	user := User{first, last, email, password}
	user.HashPassword()
	return user
}

func (user *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}
