/*
Package config consolidates and provides access to application configuration
data through the following methods with the given order of precedence (highest
to lowest):

	1. command line arguments, e.g.: --testMode=true
	2. environment variables, e.g.: APP_TEST_MODE=true
	3. configuration file properties, e.g.: app.testMode=true

The configuration parser is itself configured through parser configuration file
at: ./config/parameters_config.json, relative to the project root. The parser
configuration file has the following structure:

{
	"argSwitch": "-",
	"argSeparator": "SPACE",
	"parameters": [{
		"optionId": "configFileId"
		"cliArgument": "configFile",
		"environmentVar": "APP_CONFIG_FILE",
		"propertyName": "app.config.file",
		"defaultVal": "../../config/application.properties",
		"description": "Application configuration properties file."
	}]
}

Once parser is successfully configured, the package looks for its default
application configuration file at: ./config/application.properties, relative to
the project root. Once parsing is complete, the package presents a unified view
of the application configuration through its AssignedValue, ValueOf, and
ValueOfWithDefault functions.

The following shows a sample usage:

	import cfg "Cloudtacts/pkg/config"
	...
	if loaded, err := cfg.Parse(); !loaded {
		log.Fatal(fmt.Sprintf("Error parsing runtime parameters: %v", err))
	}
	optionId := "testModeKey"
	if cfg.AssignedValue(optionId) {
		testMode := cfg.ValueOf(optionId)
	} else {
		testMode := cfg.ValueOfWithDefault(optionId, "false")
	}
*/
package config

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/magiconair/properties"

	"Cloudtacts/pkg/model"
	"Cloudtacts/pkg/util"
)


// ValueOf returns the configuration value for the given parameter identifier.
// The parameter identifier can be any one of:
// 1. a CLI argument without its prefix switch (e.g. '--'; "userdbTestMode")
// 2. an environment variable name (e.g. "CT_USERDB_TEST_MODE"
// 3. a configuration file propert key (e.g. "user.auth.testMode")
func ValueOf(id string) string {
	if len(parameters[id]) > 0 {
		return parameters[id]
	}
	return ""
}

// ParameterWithDefault returns the configuration value for the given
// enumerated constant or the specified default value if not otherwise
// configured.
func ValueOfWithDefault(id string, defVal string) string {
	parmVal := ValueOf(id)
	if parmVal == model.USER_MUST_PROVIDE {
		parmVal = defVal
	}
	return parmVal
}

// AssignedValue returns true if the configuration parameter for the given
// enumerated constant has a user assigned value.
func AssignedValue(id string) bool {
	if len(parameters[id]) > 0 && parameters[id] != model.USER_MUST_PROVIDE {
		return true
	}
	return false
}

// Parse initializes configuration parameters with user assigned values or with
// their default values if not assigned a value by the user. User assigned
// values are given the following configuration precedence, from highest to
// lowest:
// 	1. CLI arguments
// 	2. Environment variables
// 	3. Configuration properties
func Parse() (bool, error) {
	// Load the parser configuration
	err := loadConfiguration()
	if err != nil {
		util.LogError("Failed to load parser configuration.", err)
	}

	// Set the initial default values with preference for CLI options, followed
	// by env overrides.
	options := util.ParseOptions(argSwitch, argSeparator, os.Args[1:])
	var parmVal string
	updateProps := []string{}
	for _, parm := range parserConfig.Parameters {
		switch {
		case len((*options)[parm.CliArgument]) > 0:
			parmVal = (*options)[parm.CliArgument]
		case len(os.Getenv(parm.EnvironmentVar)) > 0:
			parmVal = os.Getenv(parm.EnvironmentVar)
		case true:
			parmVal = parm.DefaultVal
			updateProps = append(updateProps, parm.PropertyName)
		}
		parameters[parm.OptionId] = parmVal
	}

	if AssignedValue(model.APP_CONFIG_FILE) {
		// ...after which, configuration properties can define any parameters
		// not already assigned.
		//
		success, err := loadProperties(ValueOf(model.APP_CONFIG_FILE), &updateProps)
		if err != nil {
			return success, err
		}
	}

	return true, nil
}


// CLI argument prefix switch (e.g. "--")
var argSwitch string = "--"

// CLI argument separator character (one of: SPACE [' '], COMMA [','], COLON [':'], SEMI-COLON [';'])
var argSeparator uint8 = ' '

// Table of parameters
var parameters map[string]string

// internal parser configuration instance
var parserConfig *model.ParserConfig


// configurationValue returns the assigned value for the specified key ('key')
// in the given properties ('props') interface unless an assigned argument
// value is given ('argVal'). Returns the specified default value ('defVal')
// if the argument value and property key are both unassigned.
func configurationValue(props *properties.Properties, key string, defVal string) string {
	propVal := props.GetString(key, model.USER_MUST_PROVIDE)
	if propVal == model.USER_MUST_PROVIDE {
		return defVal
	}
	return propVal
}

// loadConfiguration reads the parser configuration and pre-initializes the
// parser for loading the user's application configuration.
func loadConfiguration() error {
	parserConfig = new(model.ParserConfig)
	if err := util.LoadParserConfig(parserConfig); err != nil {
		util.LogError("Error loading parser configuration.", err)
	}

	parameters = make(map[string]string)

	if len(parserConfig.ArgSwitch) > 0 {
		argSwitch = parserConfig.ArgSwitch
	}

	if len(parserConfig.ArgSeparator) > 0 {
		switch parserConfig.ArgSeparator {
		case "EQUALS":
			argSeparator = '='
		case "COMMA":
			argSeparator = ','
		case "COLON":
			argSeparator = ':'
		case "SEMI-COLON":
			argSeparator = ';'
		case "SPACE":
			// default
		}
	}

	return nil
}

// loadProperties reads all key=value pair properties from the specified file
// path and assigns their values to any matching parameters not already
// assigned a value.
func loadProperties(filename string, propsList *[]string) (bool, error) {
	fmt.Println("Loading properties...")

	if filename == "" {
		return false, errors.New("File name not provided.")
	}

	props, err := properties.LoadFile(filename, properties.UTF8)
	if err != nil {
		return false, err
	}

	var parmVal string
	for _, parm := range parserConfig.Parameters {
		if slices.Contains(*propsList, parm.PropertyName) {
			parmVal = configurationValue(props, parm.PropertyName, parameters[parm.OptionId])
			parameters[parm.OptionId] = parmVal
			util.LogIt(fmt.Sprintf("Updated prop (%v): %v = %v", parm.OptionId, parm.PropertyName, parmVal))
		}
	}

	return true, nil
}
