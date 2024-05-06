package model

import "fmt"

const (
	UserExistsError = "Error 1062"
)

var (
	NoUserIdError       = UserError{"U01", "User identifier is empty!", nil}
	NoProfileIdError    = UserError{"U02", "User profile name is empty!", nil}
	NoEmailAddressError = UserError{"U03", "User e-mail address is empty!", nil}
)

// User information related data
type UserList struct {
	Users []User `json: "users"`
}

type User struct {
	CtUser string `json: "ctuser"`
	CtPass string `json: "ctpass"`
	CtProf string `json: "ctprof"`
	CtPpic string `json: "ctppic"`
	UEmail string `json: "uemail"`
	AToken string `json: "atoken"`
	LLogin string `json: "llogin"`
	UValid string `json: "uvalid"`
}

type UserError struct {
	Code    string
	Message string
	Cause   error
}

func (err UserError) Error() string {
	if err.Cause == nil {
		return fmt.Sprintf("%v: %v", err.Code, err.Message)
	} else {
		return fmt.Sprintf("%v: %v\n%v", err.Code, err.Message, err.Cause)
	}
}

func (err UserError) WithCause(src error) UserError {
	return UserError{err.Code, err.Message, src}
}
