package model

import (
	"context"
	"os"
)

const (
	// Constant APP_CONFIG_ID provides the default config file argument
	// identifier.
	APP_CONFIG_ID = "configFileId"

	// Constant PARSER_CONFIG_FILE provides the default location of the parser
	// configuration file, relative to the 'config' package.
	PARSER_CONFIG_FILE = "./config/parameters_config.json"

	// Constant APPLICATION_CONFIG_FILE provides the default location of the
	// application configuration file, relative to the 'config' package.
	APPLICATION_CONFIG_FILE = "./config/application.properties"

	// Constant USER_MUST_PROVIDE ("userMustProvide") serves as default value
	// for any parameter not assigned a value through configuration and which is
	// not otherwise assigned a default value.
	USER_MUST_PROVIDE = "userMustProvide"
)

// ApplicationConfig represents a loaded and parsed application configuration.
type ApplicationConfig interface {
	ValueOf(string) string
	ValueOfWithDefault(string, string) string
	AssignedValue(string) string
	IsParsed() bool
	Parse() (bool, error)
	Context() context.Context
}

// ParserConfig represents a loaded and parsed configuration chain. Field
// 'ArgSeparator' must be one of: "SPACE", "EQUALS", "COMMA", "COLON",
// "SEMI-COLON" (default is SPACE)
type ParserConfig struct {
	ArgSwitch    string      `json: "argSwitch"`
	ArgSeparator string      `json: "argSeparator"`
	Parameters   []Parameter `json: "parameters"`
}

// Parameter represents applicaiton options.
type Parameter struct {
	OptionId       string `json: "optionId"`
	CliArgument    string `json: "cliArgument"`
	EnvironmentVar string `json: "environmentVar"`
	PropertyName   string `json: "propertyName"`
	DefaultVal     string `json: "defaultVal"`
	Description    string `json: "description"`
}

var ParserConfigPath string
var ApplicationConfigPath string

func init() {
	if len(ParserConfigPath) < 1 {
		ParserConfigPath = os.Getenv("ParserConfigPath")
		if len(ParserConfigPath) < 1 {
			ParserConfigPath = PARSER_CONFIG_FILE
		}
	}
	if len(ApplicationConfigPath) < 1 {
		ApplicationConfigPath = os.Getenv("ApplicationConfigPath")
		if len(ApplicationConfigPath) < 1 {
			ApplicationConfigPath = APPLICATION_CONFIG_FILE
		}
	}
}
