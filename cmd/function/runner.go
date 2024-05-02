package main

import (
	"context"
	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/testdata/conformance/nondeclarative"

	"Cloudtacts/pkg/config"
	"Cloudtacts/pkg/util"
)

func main() {
	if _, err := config.Parse(); err != nil {
		util.LogError("Failed to parse configuration.", err)
	}

	port := config.ValueOfWithDefault("userdbFunctionPortId", "8088")
	target := config.ValueOf("userdbAddUserId")

	ctx := context.Background()
	if err := funcframework.RegisterHTTPFunctionContext(ctx, "/", nondeclarative.HTTP); err != nil {
		util.LogError(fmt.Sprint("Failed to register function: %v", target), err)
	}

	// By default, listen on all interfaces. If testing locally, run with
	// LOCAL_ONLY=true to avoid triggering firewall warnings and
	// exposing the server outside of your own machine.
	hostname := ""
	if localOnly := os.Getenv("LOCAL_ONLY"); localOnly == "true" {
		hostname = "127.0.0.1"
	}
	if err := funcframework.StartHostPort(hostname, port); err != nil {
		util.LogError(fmt.Sprintf("funcframework.StartHostPort/Function: %v/%v\n", port, target), err)
	}
}
