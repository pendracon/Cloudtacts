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

	uc.Close()
	t.Log("Closed DB client.\n")
}

func TestUserInfo(t *testing.T) {
	uc := connect(t)
	defer uc.Close()

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
	if !testData.Users[0].Equals(&userData) {
		t.Error("Queried data doesn't match test data.")
	}
}

func TestAddUser(t *testing.T) {
	uc := connect(t)
	defer uc.Close()

	if serr := uc.AddUser(&(testData.Users[1])); serr.IsError() {
		t.Errorf("Error adding user: %v", serr)
	}
	t.Logf("Added user: %v", testData.Users[1].CtUser)

	userData := model.User{
		CtUser: testData.Users[1].CtUser,
		CtProf: testData.Users[1].CtProf,
		UEmail: testData.Users[1].UEmail,
	}
	if serr := uc.UserInfo(&userData); serr.IsError() {
		t.Errorf("Error querying added user: %v", serr)
	}
	if !testData.Users[1].Equals(&userData) {
		t.Error("Queried data doesn't match test data.")
	}
}

func TestUpdateUser(t *testing.T) {
	uc := connect(t)
	defer uc.Close()

	userData := model.User{
		CtUser: testData.Users[1].CtUser,
		CtProf: testData.Users[1].CtProf,
		UEmail: testData.Users[1].UEmail,
		CtPpic: "org1/pendracon1/image.png",
	}

	if serr := uc.UpdateUser(&userData); serr.IsError() {
		t.Errorf("Error updating user: %v", serr)
	}
	t.Logf("Updated user: %v", testData.Users[1].CtUser)

	userData.CtPpic = ""
	if serr := uc.UserInfo(&userData); serr.IsError() {
		t.Errorf("Error querying updated user: %v", serr)
	}
	if userData.CtPpic != "org1/pendracon1/image.png" {
		t.Error("Updated data not returned in query.")
	}
}

func TestDelUser(t *testing.T) {
	uc := connect(t)
	defer uc.Close()

	userData := model.User{
		CtUser: testData.Users[1].CtUser,
		CtProf: testData.Users[1].CtProf,
		UEmail: testData.Users[1].UEmail,
	}

	if serr := uc.DeleteUser(&userData); serr.IsError() {
		t.Errorf("Error deleting user info: %v", serr)
	}
	t.Logf("Deleted user: %v", testData.Users[1].CtUser)

	if serr := uc.UserInfo(&userData); serr.IsError() {
		t.Errorf("Error querying updated user: %v", serr)
	}
	fmt.Printf("Queried back data: %v\n", userData)
	if userData.CtPpic != "" {
		t.Error("Deleted data returned in query.")
	}
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
	var err error
	cfg, err = config.ContextConfig()
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
