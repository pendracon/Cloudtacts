package model

// User information related data
type UserList struct {
	Users []User	`json: "users"`
}

type User struct {
	CtUser	string	`json: "ctuser"`
	CtPass	string	`json: "ctpass"`
	CtProf	string	`json: "ctprof"`
	CtPpic	string	`json: "ctppic"`
	UEmail	string	`json: "uemail"`
	AToken	string	`json: "atoken"`
	LLogin	string	`json: "llogin"`
	UValid	string	`json: "uvalid"`
}
