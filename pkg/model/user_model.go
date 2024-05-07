package model

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

const (
	UserExistsError = "Error 1062"

	HPWD_TAG = "H:"
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

func (u *User) PwdHash(tagged bool) string {
	hpass := u.CtPass

	if u.HasTextPwd() {
		if tagged {
			hpass = fmt.Sprintf("%v%v", HPWD_TAG, TextDigestOf(u.CtPass))
		} else {
			hpass = TextDigestOf(u.CtPass)
		}
	} else {
		if !tagged {
			hpass = u.CtPass[2:]
		}
	}

	return hpass
}

func (u *User) HasTextPwd() bool {
	return !strings.HasPrefix(u.CtPass, HPWD_TAG)
}

func (u *User) Equals(user *User) bool {
	fmt.Printf("Users: %v, %v\n", u.CtUser, user.CtUser)
	fmt.Printf("Profs: %v, %v\n", u.CtProf, user.CtProf)
	fmt.Printf("Email: %v, %v\n", u.UEmail, user.UEmail)
	fmt.Printf("PPics: %v, %v\n", u.CtPpic, user.CtPpic)
	eq := (u.CtUser == user.CtUser) &&
		(u.CtProf == user.CtProf) &&
		(u.UEmail == user.UEmail) &&
		(u.CtPpic == user.CtPpic)

	upass := u.CtPass
	if u.HasTextPwd() {
		upass = fmt.Sprintf("%v%v", HPWD_TAG, TextDigestOf(u.CtPass))
	}
	userp := user.CtPass
	if user.HasTextPwd() {
		userp = fmt.Sprintf("%v%v", HPWD_TAG, TextDigestOf(user.CtPass))
	}
	fmt.Printf("HPass: %v, %v\n", upass, userp)
	eq = eq && (upass == userp)

	return eq
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

func TextDigestOf(data string) string {
	return fmt.Sprintf("%x", DigestOf(data))
}

func DigestOf(data string) []byte {
	hashValue := sha256.Sum256([]byte(data))
	return hashValue[:]
}
