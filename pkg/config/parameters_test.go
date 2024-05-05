package config

import (
	"fmt"
	"testing"

	"Cloudtacts/pkg/model"
	"Cloudtacts/pkg/util"
)

var cfg *Config

func TestParse(t *testing.T) {
	t.Logf("Got user.auth.testMode = %v", cfg.ValueOf("user.auth.testMode"))
}

func init() {
	model.ParserConfigPath = "../../config/parameters_config.json"
	model.ApplicationConfigPath = "../../config/application.properties"
	cfg, err := ContextConfig()
	if err != nil {
		util.LogError("", "parameters_test:TestConfig", err)
	}
	util.LogIt("", fmt.Sprintf("Parsed configuration = %v", cfg.IsParsed()))
}
