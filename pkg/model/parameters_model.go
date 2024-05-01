package model

// Constant APP_CONFIG_FILE provides the default config file argument
// identifier.
const APP_CONFIG_FILE = "configFile"

// Constant PARSER_CONFIG_FILE provides the default location of the parser
// configuration file, relative to the 'config' package.
const PARSER_CONFIG_FILE = "./config/parameters_config.json"

// Constant USER_MUST_PROVIDE ("userMustProvide") serves as default value
// for any parameter not assigned a value through configuration and which is
// not otherwise assigned a default value.
const USER_MUST_PROVIDE = "userMustProvide"


// ArgSeparator must be one of: "SPACE", "EQUALS", "COMMA", "COLON",
// "SEMI-COLON"
type ParserConfig struct {
	ArgSwitch		string		`json: "argSwitch"`
	ArgSeparator	string		`json: "argSeparator"`
	Parameters		[]Parameter	`json: "parameters"`
}

type Parameter struct {
	OptionId		string	`json: "optionId"`
	CliArgument 	string	`json: "cliArgument"`
	EnvironmentVar	string	`json: "environmentVar"`
	PropertyName	string	`json: "propertyName"`
	DefaultVal		string	`json: "defaultVal"`
	Description		string	`json: "description"`
}
