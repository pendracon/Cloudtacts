package auth

import (
	"fmt"
	"testing"

	"Cloudtacts/pkg/config"
	"Cloudtacts/pkg/model"
	"Cloudtacts/pkg/util"
)

func TestNewUserDBClient(t *testing.T) {
	uc := connect()
	t.Logf("Got DB client: %v\n", uc.HostUrl())

	close(uc)
	t.Log("Closed DB client.\n")
}

func TestAddUser(t *testing.T) {
	uc := connect()
	defer close(uc)

	err := uc.AddUser(&(testData.Users[0]))
	if err != nil {
		util.LogError("Error adding user.", err)
	}
	err = uc.UpdateUser(&(testData.Users[0]))
	if err != nil {
		util.LogError("Error updating user.", err)
	}
}

func TestUserInfo(t *testing.T) {
	uc := connect()
	defer close(uc)

	userData := model.User{
		CtUser: testData.Users[0].CtUser,
		CtProf: testData.Users[0].CtProf,
		UEmail: testData.Users[0].UEmail,
	}
	uc.UserInfo(&userData)
	util.LogIt(fmt.Sprintf("testData = %v", testData.Users[0]))
	util.LogIt(fmt.Sprintf("userData = %v", userData))
}

func TestDelUser(t *testing.T) {
	uc := connect()
	defer close(uc)

	uc.DelUser(&(testData.Users[0]))
}

func close(uc *userClient) {
	CloseUserDBClient(uc)
}

func connect() *userClient {
	uc := NewUserDBClient(config.ValueOf("userdbHostId"), config.ValueOf("userdbPortId"), config.ValueOf("userdbDatabaseId"))
	return uc
}

func init() {
	util.ParserConfigPath = "../../config/parameters_config.json"
	config.Parse()
	testData = new(model.UserList)

	err := util.LoadUserListFile("../../data/TestUsers.json", testData)
	if err != nil {
		util.LogError("TestNewUser", err)
	}
}

var testData *model.UserList
