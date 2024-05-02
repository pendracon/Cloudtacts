package function

import (
    "fmt"
    "io/ioutil"
    "net/http"

    "github.com/GoogleCloudPlatform/functions-framework-go/functions"

    "Cloudtacts/pkg/config"
    "Cloudtacts/pkg/util"
)

const (
	outputFile = "function_output.json"
)

// Function addNewUser is an HTTP handler
func addNewUser(w http.ResponseWriter, r *http.Request) {
    // Your code here


    // Send an HTTP response
    fmt.Fprintln(w, "OK")
}

// HTTP is a simple HTTP function that writes the request body to the response body.
func HTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := ioutil.WriteFile(outputFile, body, 0644); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func init() {
	if _, err := config.Parse(); err != nil {
        util.LogError("Failed to parse configuration.", err)
    }

    // Register an HTTP function with the Functions Framework
    target := config.ValueOfWithDefault("userdbAddUserId", "AddUser")
    util.LogIt(fmt.Sprintf("Binding function to target '%v'.", target))
    functions.HTTP(target, addNewUser)
    functions.HTTP("Echo", HTTP)
}
