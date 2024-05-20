package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"

	"Cloudtacts/pkg/auth"
	"Cloudtacts/pkg/config"
	"Cloudtacts/pkg/model"
	"Cloudtacts/pkg/storage"
	"Cloudtacts/pkg/util"
)

const (
	functionKeyHeader = "CT-Function-Name"
	errorCodeHeader   = "CT-Error-Code"

	getUserNameDef    = "GetUser"
	addUserNameDef    = "AddUser"
	deleteUserNameDef = "DeleteUser"
	updateUserNameDef = "UpdateUser"

	getUserResponse    = "{'username': '%v', 'profile': '%v', 'email': '%v', 'imageLoc': '%v', 'lastOn': '%v', 'validatedOn': '%v'}"
	addUserResponse    = "{'username': '%v', 'profile': '%v', 'result': 'added'}"
	deleteUserResponse = "{'username': '%v', 'profile': '%v', 'result': 'deleted'}"
	updateUserResponse = "{'username': '%v', 'profile': '%v', 'result': 'updated'}"
)

var testMode bool
var cfg *config.Config

// Function getUserInfo is an HTTP handler
func getUserInfo(w http.ResponseWriter, r *http.Request) {
	logIt("Executing 'getUserInfo'...")

	user, serr := findUser(w, r, cfg.ValueOfWithDefault(model.KEY_AUTH_FUNCTION_GET, getUserNameDef))
	if !serr.IsError() {
		body := fmt.Sprintf(getUserResponse, user.CtUser, user.CtProf, user.UEmail, user.CtPpic, user.LLogin, user.UValid)
		if cnt, err := fmt.Fprintln(w, body); err != nil {
			serr = model.IOError.WithCause(err)
			logIt(fmt.Sprintf("Wrote %d bytes\nError = %v", cnt, err))
		}
	}
	if serr.IsError() {
		writeErrorResponse(w, "Error reading user info: %v/%v.", user, serr)
	}
}

// Function addNewUser is an HTTP handler
func addNewUser(w http.ResponseWriter, r *http.Request) {
	logIt("Executing 'addNewUser'...")

	user, uc, serr := connect(w, r, cfg.ValueOfWithDefault(model.KEY_AUTH_FUNCTION_ADD, addUserNameDef))
	if !serr.IsError() {
		defer uc.Close()

		if user.HasTextPwd() {
			user.CtPass = user.PwdHash(true)
		}

		if len(user.CtPpic) > 0 && !user.HasProfilePicKey() {
			_, serr = storage.SaveProfilePic(cfg, user)
		}

		if !serr.IsError() {
			serr = uc.AddUser(user)
		}

		if !serr.IsError() {
			w.WriteHeader(http.StatusCreated)
			if cnt, err := fmt.Fprintln(w, fmt.Sprintf(addUserResponse, user.CtUser, user.CtProf)); err != nil {
				serr = model.IOError.WithCause(err)
				logIt(fmt.Sprintf("Wrote %d bytes\nError = %v", cnt, err))
			}
		}
	}
	if serr.IsError() {
		var msg, tmpl string
		if strings.Contains(serr.Cause.Error(), model.UserExistsError) {
			tmpl = "Error adding new user: %v/%v - user exists."
			msg = fmt.Sprintf(tmpl, user.CtUser, user.CtProf)
		} else {
			tmpl = "Error adding new user: %v/%v."
			msg = fmt.Sprintf(tmpl, user.CtUser, user.CtProf)
		}
		writeErrorResponse(w, tmpl, user, serr)
		util.LogIt("Cloudtacts", fmt.Sprintf("%v%v", msg, serr))
	}
}

// Function deleteUser is an HTTP handler
func deleteUser(w http.ResponseWriter, r *http.Request) {
	logIt("Executing 'deleteUser'...")

	user, uc, serr := connect(w, r, cfg.ValueOfWithDefault(model.KEY_AUTH_FUNCTION_DEL, deleteUserNameDef))
	if !serr.IsError() {
		defer uc.Close()

		quser := user.Clone()
		serr = queryUser(uc, quser)
		if len(quser.CtPpic) > 0 {
			quser.CtImgt = util.ImageFileType(quser.CtPpic)
		}

		if !serr.IsError() {
			if quser.HasProfilePicKey() {
				_, serr = storage.DeleteProfilePic(cfg, quser)
				if serr.IsError() {
					util.LogIt("Cloudtacts", fmt.Sprintf("Error attempting to delete profile image: %v", serr))
				}
			}
		}

		if serr = uc.DeleteUser(user); !serr.IsError() {
			if cnt, err := fmt.Fprintln(w, fmt.Sprintf(deleteUserResponse, user.CtUser, user.CtProf)); err != nil {
				serr = model.IOError.WithCause(err)
				logIt(fmt.Sprintf("Wrote %d bytes\nError = %v", cnt, err))
			}
		}
	}
	if serr.IsError() {
		writeErrorResponse(w, "Error deleting user: %v/%v.", user, serr)
	}
}

// Function updateUser is an HTTP handler
func updateUser(w http.ResponseWriter, r *http.Request) {
	logIt("Executing 'updateUser'...")

	user, uc, serr := connect(w, r, cfg.ValueOfWithDefault(model.KEY_AUTH_FUNCTION_UPD, updateUserNameDef))
	if !serr.IsError() {
		defer uc.Close()

		if serr = uc.UpdateUser(user); !serr.IsError() {
			if cnt, err := fmt.Fprintln(w, fmt.Sprintf(updateUserResponse, user.CtUser, user.CtProf)); err != nil {
				serr = model.IOError.WithCause(err)
				logIt(fmt.Sprintf("Wrote %d bytes\nError = %v", cnt, err))
			}
		}
	}
	if serr.IsError() {
		writeErrorResponse(w, "Error updating user: %v/%v.", user, serr)
	}
}

// Function connect verifies the request and returns a corresponding user
// instance and database connection handle. An instance of ServiceError is
// returned if an error occurs.
func connect(w http.ResponseWriter, r *http.Request, requestName string) (*model.User, auth.UserDBClient, model.ServiceError) {
	var user model.User
	var uc auth.UserDBClient
	serr := model.NoError

	if ok, _ := verifyRequestFunction(w, r, requestName); ok {

		if body, serr := readRequestBody(w, r); !serr.IsError() {

			if user, serr = getUser(w, body); !serr.IsError() {
				uc, serr = auth.GetDbClient(cfg, cfg.ValueOf(model.KEY_USERDB_HOST_IP), cfg.ValueOf(model.KEY_USERDB_PORT_NUM), cfg.ValueOf(model.KEY_USERDB_DATABASE))
				if serr.IsError() {
					util.LogIt("Cloudtacts", serr.Error())
				}
			}
		}
	}

	return &user, uc, serr
}

// Function queryUser is an HTTP handler
func findUser(w http.ResponseWriter, r *http.Request, funcName string) (*model.User, model.ServiceError) {
	user, uc, serr := connect(w, r, funcName)

	util.LogIt("Cloudtacts", fmt.Sprintf("Got client: %v", uc))

	if !serr.IsError() {
		defer uc.Close()

		serr = queryUser(uc, user)
		if !serr.IsError() && user.HasProfilePicKey() {
			user.CtPpic = user.CtPpic[2:]
		}
	}

	return user, serr
}

func queryUser(uc auth.UserDBClient, user *model.User) model.ServiceError {
	serr := uc.UserInfo(user)
	if serr.IsError() {
		logIt(fmt.Sprintf("Query error: %v", serr))
	}

	return serr
}

func getUser(w http.ResponseWriter, body []byte) (model.User, model.ServiceError) {
	var userList model.UserList

	if err := util.ToUserList(body, &userList); err == nil {
		return userList.Users[0], model.NoError
	} else {
		user := model.User{}
		serr := model.InvalidMsgError.WithCause(err)
		writeErrorResponse(w, "Error converting request body.", &user, serr)
		return user, serr
	}
}

func readRequestBody(w http.ResponseWriter, r *http.Request) ([]byte, model.ServiceError) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		serr := model.InternalReadError.WithCause(err)
		writeErrorResponse(w, "Error reading request body.", &model.User{}, serr)
		return nil, serr
	}
	return body, model.NoError
}

// Function writeErrorResponse writes an HTTP status and response message back
// to the calling client. Parameter tmpl should be either: 1. a complete message or
// a message template with fmt compatible placeholders for the user identifier
// and profile name in that order. If the given user is undefined then a
// complete
func writeErrorResponse(w http.ResponseWriter, tmpl string, user *model.User, serr model.ServiceError) {
	w.Header().Add(errorCodeHeader, serr.Code)
	w.WriteHeader(model.HttpErrorStatus[serr.Code])
	if serr.Cause == nil {
		if len(user.CtUser) > 0 {
			w.Write([]byte(fmt.Sprintf("%v\n%v.", fmt.Sprintf(tmpl, user.CtUser, user.CtProf), serr.Message)))
		} else {
			w.Write([]byte(fmt.Sprintf("%v\n%v.", tmpl, serr.Message)))
		}
	} else {
		if len(user.CtUser) > 0 {
			w.Write([]byte(fmt.Sprintf("%v\n%v\n%v", fmt.Sprintf(tmpl, user.CtUser, user.CtProf), serr.Message, serr.Cause)))
		} else {
			w.Write([]byte(fmt.Sprintf("%v\n%v\n%v", tmpl, serr.Message, serr.Cause)))
		}
	}
}

func verifyRequestFunction(w http.ResponseWriter, r *http.Request, name string) (bool, string) {
	ok, hval := headerValue(r, functionKeyHeader)
	if !ok {
		w.Header().Add(errorCodeHeader, model.InvalidMsgError.Code)
		w.WriteHeader(model.HttpErrorStatus[model.InvalidMsgError.Code]) // Bad Request: missing required header
		w.Write([]byte("Missing reqired header.\n"))
		return false, ""
	}
	if hval != name {
		w.Header().Add(errorCodeHeader, model.InvalidMsgError.Code)
		w.WriteHeader(model.HttpErrorStatus[model.InvalidMsgError.Code]) // Bad Request: wrong function request
		w.Write([]byte("Wrong function requested.\n"))
		return false, ""
	}
	return ok, hval
}

func headerValue(r *http.Request, key string) (bool, string) {
	val := r.Header.Get(key)
	if len(val) > 0 {
		return true, val
	} else {
		return false, val
	}
}

func logIt(message string) {
	if testMode {
		util.LogIt("Cloudtacts", message)
	}
}

func init() {
	cfgx, err := config.ContextConfig()
	if err != nil {
		util.LogError("Cloudtacts", "function - Failed to parse configuration.", err)
	}
	util.LogIt("Cloudtacts", fmt.Sprintf("Parsed configuration = %v", cfgx.IsParsed()))
	testMode = (cfgx.ValueOfWithDefault(model.KEY_USERDB_TEST_MODE, "false") == "true")

	// Register an HTTP function with the Functions Framework
	targetList := [4][2]string{
		{model.KEY_AUTH_FUNCTION_GET, getUserNameDef},
		{model.KEY_AUTH_FUNCTION_ADD, addUserNameDef},
		{model.KEY_AUTH_FUNCTION_DEL, deleteUserNameDef},
		{model.KEY_AUTH_FUNCTION_UPD, updateUserNameDef},
	}
	for _, targetName := range targetList {
		target := cfgx.ValueOfWithDefault(targetName[0], targetName[1])
		switch targetName[1] {
		case "GetUser":
			logIt(fmt.Sprintf("Binding function to target '%v'->'getUserInfo'.", target))
			functions.HTTP(target, getUserInfo)
		case "AddUser":
			logIt(fmt.Sprintf("Binding function to target '%v'->'addNewUser'.", target))
			functions.HTTP(target, addNewUser)
		case "DeleteUser":
			logIt(fmt.Sprintf("Binding function to target '%v'->'deleteUser'.", target))
			functions.HTTP(target, deleteUser)
		case "UpdateUser":
			logIt(fmt.Sprintf("Binding function to target '%v'->'updateUser'.", target))
			functions.HTTP(target, updateUser)
		}
	}
	cfg = cfgx
}
