package main

import (
	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/testdata/conformance/nondeclarative"

	"Cloudtacts/pkg/config"
	"Cloudtacts/pkg/util"
)

func main() {
	var cfg *config.Config
	var err error

	if cfg, err = config.ContextConfig(); err != nil {
		util.LogError("CloudtactsRunner", "Failed to parse configuration.", err)
	}

	port := cfg.ValueOfWithDefault("userdbFunctionPortId", "8088")

	if err := funcframework.RegisterHTTPFunctionContext(cfg.Context(), "/", nondeclarative.HTTP); err != nil {
		util.LogError("CloudtactsRunner", "Failed to register function context.", err)
	}

	// By default, listen on all interfaces. If testing locally, run with
	// LOCAL_ONLY=true to avoid triggering firewall warnings and
	// exposing the server outside of your own machine.
	hostname := ""
	if localOnly := os.Getenv("LOCAL_ONLY"); localOnly == "true" {
		hostname = "127.0.0.1"
	}

	util.LogIt("CloudtactsRunner", fmt.Sprintf("Starting function handler on port %v.", port))
	if err := funcframework.StartHostPort(hostname, port); err != nil {
		util.LogError("CloudtactsRunner", fmt.Sprintf("funcframework.StartHostPort: %v\n", port), err)
	}
}
