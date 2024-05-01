package auth

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"Cloudtacts/pkg/model"
	"Cloudtacts/pkg/util"
)

const (
	SELECT_USER_INFO string = "SELECT ctpass, ctppic, atoken, llogin, uvalid FROM user WHERE ctuser = ? AND ctprof = ? AND uemail = ?"
	INSERT_USER_STMT string = "INSERT INTO user (ctuser, ctpass, ctprof, uemail, ctppic) VALUES(?, ?, ?, ?, ?)"
	DELETE_USER_STMT string = "DELETE FROM user WHERE ctuser = ? AND ctprof = ? AND uemail = ?"
	UPDATE_USER_STMT string = "UPDATE user SET ctpass = ?, ctppic = ?, atoken = ?, llogin = ?, uvalid = ? WHERE ctuser = ? AND ctprof = ? AND uemail = ?"
)

type UserError struct {
	mesg	string
}

type UserDBClient interface {
	// Updates the referenced user instance with information from the database.
	UserInfo(user *model.User) error

	// Adds the referenced user information to the database.
	AddUser(user *model.User) error

	// Deletes the referenced user information from the database.
	DelUser(user *model.User) error

	// Updates the referenced user information in the database.
	UpdateUser(user *model.User) error

	// Return host URL of the database.
	HostUrl() string
}

func (err UserError) Error() string {
	return err.mesg
}

func (uc *userClient) UserInfo(user *model.User) error {
	if ok, err := validateUserKey(user); !ok {
		return err
	}

	rows, err := uc.dbClient.Query(SELECT_USER_INFO, user.CtUser, user.CtProf, user.UEmail)
	if err != nil {
		util.LogError("Error querying user info.", err)
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&user.CtPass, &user.CtPpic, &user.AToken, &user.LLogin, &user.UValid)
		if err != nil {
			util.LogError("Error scanning user info.", err)
		}
	}

	err = rows.Err()
	if err != nil {
		util.LogError("Got unknown results error.", err)
	}

	return err
}

func (uc *userClient) AddUser(user *model.User) error {
	if ok, err := validateUserKey(user); !ok {
		return err
	}

	stmtIns, err := uc.dbClient.Prepare(INSERT_USER_STMT)
	if err != nil {
		util.LogError("Error preparing add user statement.", err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(user.CtUser, user.CtPass, user.CtProf, user.UEmail, user.CtPpic)
	if err != nil {
		util.LogError("Error executing add user.", err)
	}

	return nil
}

func (uc *userClient) DelUser(user *model.User) error {
	if ok, err := validateUserKey(user); !ok {
		return err
	}

	stmtDel, err := uc.dbClient.Prepare(DELETE_USER_STMT)
	if err != nil {
		util.LogError("Error preparing delete user statement.", err)
	}
	defer stmtDel.Close()

	_, err = stmtDel.Exec(user.CtUser, user.CtProf, user.UEmail)
	if err != nil {
		util.LogError("Error executing delete user.", err)
	}

	return nil
}

func (uc *userClient) UpdateUser(user *model.User) error {
	if ok, err := validateUserKey(user); !ok {
		return err
	}

	stmtUpd, err := uc.dbClient.Prepare(UPDATE_USER_STMT)
	if err != nil {
		util.LogError("Error preparing update user statement.", err)
	}
	defer stmtUpd.Close()

	_, err = stmtUpd.Exec(user.CtPass, user.CtPpic, user.AToken, user.LLogin, user.UValid, user.CtUser, user.CtProf, user.UEmail)
	if err != nil {
		util.LogError("Error executing update user.", err)
	}

	return nil
}

func (uc *userClient) HostUrl() string {
	return uc.hostUrl
}

func NewUserDBClient(host, port, database string) *userClient {
	hostUrl := fmt.Sprintf("%v:%v", host, port)
	client := new(userClient)
	client.hostUrl = hostUrl
	var err error

	if hostUrl != "" {
		client.dbClient, err = sql.Open("mysql", fmt.Sprintf("root:devpass@tcp(%v)/%v?tls=skip-verify&autocommit=true&parseTime=true", hostUrl, database))

		if err != nil {
			util.LogError(fmt.Sprintf("Error opening user database on host %v.", hostUrl), err)
		}
		util.LogIt(fmt.Sprintf("DB client using user database on host %v.", hostUrl))
	} else {
		client.dbClient, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/cloudtacts")

		if err != nil {
			util.LogError("Error opening user database on localhost.", err)
		}
		util.LogIt("DB client using user database at localhost.")
	}

	if err == nil {
		client.dbClient.SetConnMaxLifetime(time.Minute * 3)
		client.dbClient.SetMaxOpenConns(10)
		client.dbClient.SetMaxIdleConns(10)
	}

	return client
}

func CloseUserDBClient(client *userClient) {
	client.dbClient.Close()
}

func validateUserKey(user *model.User) (bool, error) {
	switch {
		case len(user.CtUser) == 0:
			return false, UserError{"User identifier is empty!"}
		case len(user.CtProf) == 0:
			return false, UserError{"User profile name is empty!"}
		case len(user.UEmail) == 0:
			return false, UserError{"User e-mail address is empty!"}
	}

	return true, nil
}

type userClient struct {
	hostUrl  string
	dbClient *sql.DB
}
