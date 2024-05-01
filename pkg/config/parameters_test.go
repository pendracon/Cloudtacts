package config

import (
	"fmt"
	"testing"

	//"Cloudtacts/pkg/model"
	"Cloudtacts/pkg/util"
)

func TestParse(t *testing.T) {
	t.Logf("Got user.auth.testMode = %v", ValueOf("user.auth.testMode"))
}

func init() {
	ok, err := Parse()
	if err != nil {
		util.LogError("TestConfig", err)
	}
	util.LogIt(fmt.Sprintf("Parsed configuration = %v", ok))
}
