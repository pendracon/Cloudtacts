package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"Cloudtacts/pkg/config"
	"Cloudtacts/pkg/model"
	"Cloudtacts/pkg/util"
)

const (
	FUNCTION_URL = "http://%v:%v/%v"
)

func postRequest(cfg *config.Config, token *string, function, body string) (int, string, model.ServiceError) {
	serr := model.NoError

	client := &http.Client{}
	url := fmt.Sprintf(FUNCTION_URL,
		cfg.ValueOf(model.KEY_AUTH_FUNCTION_HOST),
		cfg.ValueOf(model.KEY_AUTH_FUNCTION_PORT),
		function)

	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		serr = model.ClientRequestError.WithCause(err)
		logIt(fmt.Sprintf("Error creating new request instance: %v.", serr))
		return 0, "", serr
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("CT-Function-Name", function)
	if len(*token) > 0 {
		req.Header.Add("CT-User-Token", *token)
	}

	resp, err := client.Do(req)
	if err != nil {
		serr = model.ClientRequestError.WithCause(err)
		logIt(fmt.Sprintf("Error executing request: %v.", serr))
		return 0, "", serr
	}
	defer resp.Body.Close()

	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	logIt(fmt.Sprintf("Got %d length response", size))

	var buff []byte
	if err != nil {
		serr = model.ClientProtocolError.WithCause(err)
	} else {
		buff = make([]byte, size-1)
		if size, err = resp.Body.Read(buff); err != nil {
			serr = model.ClientReadError.WithCause(err)
			logIt(fmt.Sprintf("Error reading response of size %d: %v.", size, serr))
			return 0, "", serr
		}
	}

	if len(resp.Header.Get("CT-User-Token")) > 0 {
		*token = resp.Header.Get("CT-User-Token")
	}

	return resp.StatusCode, string(buff[:]), serr
}

func readInput(cfg *config.Config) (string, model.ServiceError) {
	serr := model.NoError
	var data string

	jsonFile, err := os.Open(cfg.ValueOf(model.KEY_CLIENT_INPUT_FILE))
	if err != nil {
		serr = model.ClientInputError.WithCause(err)
	} else {
		defer jsonFile.Close()

		bbuff, err := io.ReadAll(jsonFile)
		if err != nil {
			serr = model.ClientInputError.WithCause(err)
		} else {
			data = string(bbuff[:])
		}
	}

	return data, serr
}

func loadImage(cfg *config.Config) ([]byte, string, model.ServiceError) {
	serr := model.NoError
	ifileName := cfg.ValueOf(model.KEY_CLIENT_IMAGE_FILE)
	itype := strings.ToLower(cfg.ValueOf(model.KEY_CLIENT_IMAGE_TYPE))

	var data []byte

	imgFile, err := os.Open(ifileName)
	if err != nil {
		serr = model.ClientImageError.WithCause(err)
	} else {
		defer imgFile.Close()

		data, err = io.ReadAll(imgFile)
		if err != nil {
			serr = model.ClientImageError.WithCause(err)
		}
		if !cfg.AssignedValue(model.KEY_CLIENT_IMAGE_TYPE) {
			itype = util.ImageFileType(ifileName)
		}
	}

	return data, itype, serr
}

func attachImage(data string, img []byte, itype string) (string, model.ServiceError) {
	var userList model.UserList
	var rval string
	serr := model.NoError

	err := util.ToUserList([]byte(data), &userList)
	if err == nil {
		userList.Users[0].CtPpic = base64.StdEncoding.EncodeToString(img)
		userList.Users[0].CtImgt = itype
		if bbuff, err := json.Marshal(userList); err == nil {
			rval = string(bbuff[:])
		} else {
			serr = model.ClientError.WithCause(err)
		}
	} else {
		serr = model.ClientError.WithCause(err)
	}

	return rval, serr
}

func attachCreds(data, creds string) (string, model.ServiceError) {
	var userList model.UserList
	var rval string
	serr := model.NoError

	err := util.ToUserList([]byte(data), &userList)
	if err == nil {
		userList.Users[0].CtPass = creds
		if bbuff, err := json.Marshal(userList); err == nil {
			rval = string(bbuff[:])
		} else {
			serr = model.ClientError.WithCause(err)
		}
	} else {
		serr = model.ClientError.WithCause(err)
	}

	return rval, serr
}

func main() {
	var cfg *config.Config
	var err error
	var data string
	serr := model.NoError

	if cfg, err = config.ContextConfig(); err != nil {
		util.LogError("Client", "Failed to parse configuration.", err)
	}

	switch cfg.ValueOf(model.KEY_CLIENT_COMMAND) {
	case "GetUser":
	case "AddUser":
	case "DeleteUser":
	case "UpdateUser":
	case "ValidateUser":
	case "LoginUser":
	default:
		if !cfg.AssignedValue(model.KEY_CLIENT_COMMAND) {
			util.LogError("Client", "Command not specified", nil)
		} else {
			util.LogError("Client", "Unknown command specified", nil)
		}
	}

	if !cfg.AssignedValue(model.KEY_CLIENT_INPUT_FILE) {
		util.LogError("Client", "Input not specified.", nil)
	}

	data, serr = readInput(cfg)

	if !serr.IsError() && cfg.AssignedValue(model.KEY_CLIENT_USER_CREDS) {
		data, serr = attachCreds(data, cfg.ValueOf(model.KEY_CLIENT_USER_CREDS))
	}

	if !serr.IsError() && cfg.AssignedValue(model.KEY_CLIENT_IMAGE_FILE) {
		var img []byte
		var itype string

		img, itype, serr = loadImage(cfg)
		if !serr.IsError() {
			data, serr = attachImage(data, img, itype)
		}
	}

	if !serr.IsError() {
		var status int
		var resp string
		var token string

		status, resp, serr = postRequest(cfg, &token, cfg.ValueOf(model.KEY_CLIENT_COMMAND), data)

		if !serr.IsError() {
			logIt(fmt.Sprintf("Command executed with status %v", status))
			logIt(fmt.Sprintf("User token: %v", token))
			logIt(fmt.Sprintf("Response:\n%v", resp))
		}
	}

	if serr.IsError() {
		util.LogError("Client", serr.Message, serr.Cause)
	}
}

func logIt(message string) {
	util.LogIt("Client", message)
}
