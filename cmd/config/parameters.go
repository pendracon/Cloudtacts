package main

import (
	"fmt"
	"log"

	"Cloudtacts/pkg/config"
	//"Cloudtacts/pkg/model"
	"Cloudtacts/pkg/util"
)

func main() {
	logIt("Running VTIS 'parameters' demo...")

	logIt(fmt.Sprintf("Got userdbTestMode = %v", config.ValueOf("userdbTestModeId")))
	logIt(fmt.Sprintf("Got userdbHost = %v", config.ValueOf("userdbHostId")))
}

func logIt(message string) {
	log.Println(fmt.Sprintf("Parameters CLI - %v", message))
}

func init() {
	ok, err := config.Parse()
	if err != nil {
		util.LogError("Parameters CLI", err)
	}
	logIt(fmt.Sprintf("Parsed configuration = %v", ok))
}
