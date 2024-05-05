package auth

import (
	"fmt"
	"testing"

	"Cloudtacts/pkg/config"
	"Cloudtacts/pkg/model"
	"Cloudtacts/pkg/util"
)

var cfg *config.Config
var testData *model.UserList

func TestNewUserDBClient(t *testing.T) {
	uc := connect(t)
	t.Logf("Got DB client: %v\n", uc.HostUrl())

	close(uc)
	t.Log("Closed DB client.\n")
}

func TestAddUser(t *testing.T) {
	uc := connect(t)
	defer close(uc)

	if serr := uc.AddUser(&(testData.Users[0])); serr.IsError() {
		t.Errorf("Error adding user: %v", serr)
	}
	t.Logf("Added user: %v", testData.Users[0].CtUser)

	if serr := uc.UpdateUser(&(testData.Users[0])); serr.IsError() {
		t.Errorf("Error updating user: %v", serr)
	}
	t.Logf("Updated user: %v", testData.Users[0].CtUser)
}

func TestUserInfo(t *testing.T) {
	uc := connect(t)
	defer close(uc)

	userData := model.User{
		CtUser: testData.Users[0].CtUser,
		CtProf: testData.Users[0].CtProf,
		UEmail: testData.Users[0].UEmail,
	}
	if serr := uc.UserInfo(&userData); serr.IsError() {
		t.Errorf("Error getting user info: %v", serr)
	}
	t.Logf("testData = %v", testData.Users[0])
	t.Logf("userData = %v", userData)
	if testData.Users[0] != userData {
		t.Error("Queried data doesn't match test data.")
	}
}

func TestDelUser(t *testing.T) {
	uc := connect(t)
	defer close(uc)

	if serr := uc.DeleteUser(&(testData.Users[0])); serr.IsError() {
		t.Errorf("Error deleting user info: %v", serr)
	}
	t.Logf("Deleted user: %v", testData.Users[0].CtUser)
}

func close(uc *userClient) {
	CloseUserDBClient(uc)
}

func connect(t *testing.T) *userClient {
	uc, serr := GetDbClient(cfg, cfg.ValueOf("userdbHostId"), cfg.ValueOf("userdbPortId"), cfg.ValueOf("userdbDatabaseId"))
	if serr.IsError() {
		t.Errorf("Error getting DB client: %v", serr)
	}

	return uc
}

func init() {
	model.ParserConfigPath = "../../config/parameters_config.json"
	model.ApplicationConfigPath = "../../config/application.properties"
	cfg, err := config.ContextConfig()
	if err != nil {
		util.LogError("", "parameters_test:TestConfig", err)
	}
	util.LogIt("", fmt.Sprintf("Parsed configuration = %v", cfg.IsParsed()))

	testData = new(model.UserList)
	err = util.LoadUserListFile("../../data/TestUsers.json", testData)
	if err != nil {
		util.LogError("Cloudtacts", "user_test:TestNewUser", err)
	}
}
