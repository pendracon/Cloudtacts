package auth

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"Cloudtacts/pkg/config"
	"Cloudtacts/pkg/model"
	"Cloudtacts/pkg/util"
)

const (
	SELECT_USER_INFO string = "SELECT ctpass, ctppic, atoken, llogin, uvalid FROM user WHERE ctuser = ? AND ctprof = ? AND uemail = ?"
	INSERT_USER_STMT string = "INSERT INTO user (ctuser, ctpass, ctprof, uemail, ctppic) VALUES(?, ?, ?, ?, ?)"
	UPDATE_USER_STMT string = "UPDATE user SET ctpass = ?, ctppic = ?, atoken = ? WHERE ctuser = ? AND ctprof = ? AND uemail = ?"
	DELETE_USER_STMT string = "DELETE FROM user WHERE ctuser = ? AND ctprof = ? AND uemail = ?"

	UPDATE_USER_STMT_TMPL string = "UPDATE user SET %v = ? WHERE ctuser = ? AND ctprof = ? AND uemail = ?"
)

type UserDBClient interface {
	// Updates the referenced user instance with information from the database.
	UserInfo(*model.User) model.ServiceError

	// Adds the referenced user information to the database.
	AddUser(*model.User) model.ServiceError

	// Deletes the referenced user information from the database.
	DeleteUser(*model.User) model.ServiceError

	// Updates the referenced user information in the database.
	UpdateUser(*model.User) model.ServiceError

	// Return host URL of the database.
	HostUrl() string

	// Closes the client connection.
	Close()
}

func (uc *userClient) UserInfo(user *model.User) model.ServiceError {
	var ferr model.ServiceError

	if ok, err := validateUserKey(user); !ok {
		return model.InvalidKeyError.WithCause(err)
	}

	rows, err := uc.conn.Query(SELECT_USER_INFO, user.CtUser, user.CtProf, user.UEmail)
	if err != nil {
		ferr = model.DbQueryError.WithCause(err)
	} else {
		defer rows.Close()

		if rows.Next() {
			ctppic := make([]byte, 52)
			atoken := make([]byte, 20)
			llogin := make([]byte, 14)
			uvalid := make([]byte, 14)
			err := rows.Scan(&user.CtPass, &ctppic, &atoken, &llogin, &uvalid)
			if err != nil {
				ferr = model.DbScanError.WithCause(err)
			} else {
				user.CtPpic = string(ctppic[:])
				user.AToken = string(atoken[:])
				user.LLogin = string(llogin[:])
				user.UValid = string(uvalid[:])
			}
		}

		err = rows.Err()
		if err != nil {
			ferr = model.DbResultsError.WithCause(err)
		}
	}

	return ferr
}

func (uc *userClient) AddUser(user *model.User) model.ServiceError {
	var ferr = model.ServiceError{}

	if ok, err := validateUserKey(user); !ok {
		ferr = model.InvalidKeyError.WithCause(err)
	}

	stmtIns, err := uc.conn.Prepare(INSERT_USER_STMT)
	if err != nil {
		ferr = model.DbPrepareError.WithCause(err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(user.CtUser, user.CtPass, user.CtProf, user.UEmail, user.CtPpic)
	if err != nil {
		ferr = model.DbInsertError.WithCause(err)
	}

	return ferr
}

func (uc *userClient) DeleteUser(user *model.User) model.ServiceError {
	var ferr = model.ServiceError{}

	if ok, err := validateUserKey(user); !ok {
		ferr = model.InvalidKeyError.WithCause(err)
	}

	stmtDel, err := uc.conn.Prepare(DELETE_USER_STMT)
	if err != nil {
		ferr = model.DbPrepareError.WithCause(err)
	}
	defer stmtDel.Close()

	_, err = stmtDel.Exec(user.CtUser, user.CtProf, user.UEmail)
	if err != nil {
		ferr = model.DbExecuteError.WithCause(err)
	}

	return ferr
}

func (uc *userClient) UpdateUser(user *model.User) model.ServiceError {
	var ferr = model.NoError

	var err error
	if ok, err := validateUserKey(user); !ok {
		ferr = model.InvalidKeyError.WithCause(err)
	}

	var stmtUpd *sql.Stmt
	if ferr == model.NoError {
		stmtUpd, err = uc.conn.Prepare(UPDATE_USER_STMT)
		if err != nil {
			ferr = model.DbPrepareError.WithCause(err)
		}
	}

	if ferr == model.NoError {
		_, err = stmtUpd.Exec(user.CtPass, user.CtPpic, user.AToken, user.CtUser, user.CtProf, user.UEmail)
		if err != nil {
			ferr = model.DbExecuteError.WithCause(err)
		}
	}

	if ferr == model.NoError && len(user.LLogin) > 0 {
		ferr = updateDateTimeColumn(user, uc, "llogin", user.LLogin)
	}
	if ferr == model.NoError && len(user.UValid) > 0 {
		ferr = updateDateTimeColumn(user, uc, "uvalid", user.UValid)
	}

	return ferr
}

func (uc *userClient) HostUrl() string {
	return uc.hostUrl
}

func (uc *userClient) Close() {
	uc.conn.Close()
}

func GetDbClient(cfg *config.Config, host, port, database string) (*userClient, model.ServiceError) {
	var serr model.ServiceError

	util.LogIt("Cloudtacts", fmt.Sprintf("Getting client for host '%v', port '%v', database '%v'...", host, port, database))

	hostUrl := fmt.Sprintf("%v:%v", host, port)
	appDbClient := new(userClient)
	appDbClient.hostUrl = hostUrl
	serr = initClient(cfg, appDbClient, database)

	return appDbClient, serr
}

type userClient struct {
	hostUrl string
	conn    *sql.DB
}

func initClient(cfg *config.Config, uc *userClient, database string) model.ServiceError {
	var serr model.ServiceError
	var err error

	if uc.hostUrl != "" {
		uc.conn, err = sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v?tls=skip-verify&autocommit=true&parseTime=true", cfg.ValueOf(model.KEY_USERDB_LOGIN), cfg.ValueOf(model.KEY_USERDB_PASSWORD), uc.hostUrl, database))

		if err != nil {
			serr = model.DbOpenError.WithCause(err)
		} else {
			traceIt(cfg, fmt.Sprintf("DB client using user database on host %v.", uc.hostUrl))
		}
	} else {
		uc.conn, err = sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/cloudtacts", cfg.ValueOf(model.KEY_USERDB_LOGIN), cfg.ValueOf(model.KEY_USERDB_PASSWORD)))

		if err != nil {
			serr = model.DbOpenError.WithCause(err)
		} else {
			traceIt(cfg, "DB client using user database at localhost.")
		}
	}

	if serr == (model.ServiceError{}) {
		if ival, err := strconv.Atoi(cfg.ValueOfWithDefault(model.KEY_USERDB_MAX_POOL, "-1")); err == nil {
			uc.conn.SetMaxOpenConns(ival)
		} else {
			uc.conn.SetMaxOpenConns(0)
		}
		if ival, err := strconv.Atoi(cfg.ValueOfWithDefault(model.KEY_USERDB_MAX_IDLE, "2")); err == nil {
			uc.conn.SetMaxIdleConns(ival)
		} else {
			uc.conn.SetMaxIdleConns(2)
		}
		if ival, err := strconv.Atoi(cfg.ValueOfWithDefault(model.KEY_USERDB_MAX_IDTM, "300")); err == nil {
			uc.conn.SetConnMaxIdleTime(time.Second * time.Duration(ival))
		} else {
			uc.conn.SetConnMaxIdleTime(time.Second * 300)
		}
		if ival, err := strconv.Atoi(cfg.ValueOfWithDefault(model.KEY_USERDB_MAX_LFTM, "30")); err == nil {
			uc.conn.SetConnMaxLifetime(time.Minute * time.Duration(ival))
		} else {
			uc.conn.SetConnMaxLifetime(time.Minute * 30)
		}

		traceIt(cfg, fmt.Sprintf("Initial stats: %v", clientStats(uc)))
	}

	return (model.ServiceError{})
}

func updateDateTimeColumn(user *model.User, uc *userClient, colName, colVal string) model.ServiceError {
	var ferr = model.NoError

	var err error
	var stmtUpd *sql.Stmt
	if ferr == model.NoError {
		stmtUpd, err = uc.conn.Prepare(fmt.Sprintf(UPDATE_USER_STMT_TMPL, colName))
		if err != nil {
			ferr = model.DbPrepareError.WithCause(err)
		}
	}

	dtime := 0
	if ferr == model.NoError {
		defer stmtUpd.Close()

		if dtime, err = strconv.Atoi(colVal); err != nil {
			ferr = model.DatetimeError.WithCause(err)
		}
	}

	if ferr == model.NoError {
		_, err = stmtUpd.Exec(dtime, user.CtUser, user.CtProf, user.UEmail)
		if err != nil {
			ferr = model.DbExecuteError.WithCause(err)
		}
	}

	return ferr
}

func clientStats(uc *userClient) string {
	return fmt.Sprintf("%v, %v, %v, %v", uc.conn.Stats().MaxOpenConnections,
		uc.conn.Stats().MaxIdleClosed, uc.conn.Stats().MaxIdleTimeClosed,
		uc.conn.Stats().MaxLifetimeClosed)
}

func validateUserKey(user *model.User) (bool, model.UserError) {
	switch {
	case len(user.CtUser) == 0:
		return false, model.NoUserIdError
	case len(user.CtProf) == 0:
		return false, model.NoProfileIdError
	case len(user.UEmail) == 0:
		return false, model.NoEmailAddressError
	}

	return true, model.UserError{}
}

func traceIt(cfg *config.Config, message string) {
	if cfg.ValueOfWithDefault(model.KEY_USERDB_TEST_MODE, "false") == "true" {
		util.LogIt("Cloudtacts", message)
	}
}
