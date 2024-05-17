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
var statusMap map[string]int

// Function getUserInfo is an HTTP handler
func getUserInfo(w http.ResponseWriter, r *http.Request) {
	logIt("Executing 'getUserInfo'...")

	if user, serr := findUser(w, r, cfg.ValueOfWithDefault(model.KEY_AUTH_FUNCTION_GET, getUserNameDef)); !serr.IsError() {
		fmt.Fprintln(w, fmt.Sprintf(getUserResponse, user.CtUser, user.CtProf, user.UEmail, user.CtPpic, user.LLogin, user.UValid))
	} else {
		w.Header().Add(errorCodeHeader, model.SystemError.Code)
		w.WriteHeader(statusMap[model.SystemError.Code])
		w.Write([]byte(fmt.Sprintf("Error reading user info: %v/%v.\n%v", user.CtUser, user.CtProf, serr.Cause)))
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
			fmt.Fprintln(w, fmt.Sprintf(addUserResponse, user.CtUser, user.CtProf))
		} else {
			var msg string
			status := statusMap[serr.Code]
			if strings.Contains(serr.Cause.Error(), model.UserExistsError) {
				msg = fmt.Sprintf("Error adding new user: %v/%v - user exists.\n", user.CtUser, user.CtProf)
			} else {
				msg = fmt.Sprintf("Error adding new user: %v/%v.\n", user.CtUser, user.CtProf)
			}
			util.LogIt("Cloudtacts", fmt.Sprintf("%v%v", msg, serr))
			w.Header().Add(errorCodeHeader, serr.Code)
			w.WriteHeader(status)
			w.Write([]byte(msg))
		}
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
			fmt.Fprintln(w, fmt.Sprintf(deleteUserResponse, user.CtUser, user.CtProf))
		} else {
			w.Header().Add(errorCodeHeader, serr.Code)
			w.WriteHeader(statusMap[serr.Code])
			w.Write([]byte(fmt.Sprintf("Error deleting user: %v/%v.\n", user.CtUser, user.CtProf)))
		}
	}
}

// Function updateUser is an HTTP handler
func updateUser(w http.ResponseWriter, r *http.Request) {
	logIt("Executing 'updateUser'...")

	user, uc, serr := connect(w, r, cfg.ValueOfWithDefault(model.KEY_AUTH_FUNCTION_UPD, updateUserNameDef))
	if !serr.IsError() {
		defer uc.Close()

		if serr = uc.UpdateUser(user); !serr.IsError() {
			fmt.Fprintln(w, fmt.Sprintf(updateUserResponse, user.CtUser, user.CtProf))
		} else {
			w.Header().Add(errorCodeHeader, serr.Code)
			w.WriteHeader(statusMap[serr.Code])
			w.Write([]byte(fmt.Sprintf("Error updating user: %v/%v.\n", user.CtUser, user.CtProf)))
		}
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
		w.Header().Add(errorCodeHeader, model.SystemError.Code)
		w.WriteHeader(statusMap[model.SystemError.Code])
		w.Write([]byte("Error converting request body.\n"))
		return model.User{}, model.InvalidMsgError.WithCause(err)
	}
}

func readRequestBody(w http.ResponseWriter, r *http.Request) ([]byte, model.ServiceError) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Add(errorCodeHeader, model.InternalReadError.Code)
		w.WriteHeader(statusMap[model.InternalReadError.Code])
		w.Write([]byte("Error reading request body.\n"))
		return nil, model.InternalReadError.WithCause(err)
	}
	return body, model.NoError
}

func verifyRequestFunction(w http.ResponseWriter, r *http.Request, name string) (bool, string) {
	ok, hval := headerValue(r, functionKeyHeader)
	if !ok {
		w.Header().Add(errorCodeHeader, model.InvalidMsgError.Code)
		w.WriteHeader(statusMap[model.InvalidMsgError.Code]) // Bad Request: missing required header
		w.Write([]byte("Missing reqired header.\n"))
		return false, ""
	}
	if hval != name {
		w.Header().Add(errorCodeHeader, model.InvalidMsgError.Code)
		w.WriteHeader(statusMap[model.InvalidMsgError.Code]) // Bad Request: wrong function request
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

	statusMap = make(map[string]int)
	statusMap[model.DbQueryError.Code] = 500
	statusMap[model.DbScanError.Code] = 500
	statusMap[model.DbResultsError.Code] = 500
	statusMap[model.DbInsertError.Code] = 500
	statusMap[model.DbPrepareError.Code] = 500
	statusMap[model.DbExecuteError.Code] = 500
	statusMap[model.DbClientError.Code] = 500
	statusMap[model.DbOpenError.Code] = 500
	statusMap[model.InvalidKeyError.Code] = 400
	statusMap[model.InvalidMsgError.Code] = 400
	statusMap[model.InternalReadError.Code] = 500
	statusMap[model.SystemError.Code] = 500
	statusMap[model.DatetimeError.Code] = 500
}
