package main

import (
	"fmt"

	"Cloudtacts/pkg/config"
	//"Cloudtacts/pkg/model"
	"Cloudtacts/pkg/util"
)

func main() {
	cfg, err := config.ContextConfig()
	if err != nil {
		util.LogError("Parameters", "TestConfig", err)
	}
	logIt(fmt.Sprintf("Parsed configuration = %v", cfg.IsParsed()))

	logIt("Running VTIS 'parameters' demo...")
	logIt(fmt.Sprintf("Got userdbTestMode = %v", cfg.ValueOf("userdbTestModeId")))
	logIt(fmt.Sprintf("Got userdbHost = %v", cfg.ValueOf("userdbHostId")))
}

func logIt(message string) {
	util.LogIt("ParametersCLI", message)
}
