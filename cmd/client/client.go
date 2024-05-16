package main

import (
	"os"

	"Cloudtacts/pkg/config"
	"Cloudtacts/pkg/model"
	"Cloudtacts/pkg/util"
)

func main() {
	var cfg *config.Config
	var err error

	if cfg, err = config.ContextConfig(); err != nil {
		util.LogError("CloudtactsClient", "Failed to parse configuration.", err)
	}

	port := cfg.ValueOfWithDefault(model.KEY_USERDB_FUNCTION_PORT, "8088")

	// By default, listen on all interfaces. If testing locally, run with
	// LOCAL_ONLY=true to avoid triggering firewall warnings and
	// exposing the server outside of your own machine.
	hostname := ""
	if localOnly := os.Getenv("LOCAL_ONLY"); localOnly == "true" {
		hostname = "127.0.0.1"
	}

	util.LogIt("CloudtactsClient", "Starting client...")

}
