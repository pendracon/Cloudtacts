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
	"context"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/magiconair/properties"

	"Cloudtacts/pkg/model"
	"Cloudtacts/pkg/util"
)

// A context configuration containing a loaded and parsed application
// configuration chain (see Parse).
type Config struct {
	// CLI argument prefix switch (e.g. "--")
	argSwitch string

	// CLI argument separator character (one of: SPACE [' '], COMMA [','], COLON [':'], SEMI-COLON [';'])
	argSeparator uint8

	// Table of parameters
	parameters map[string]string

	// internal parser configuration instance
	parserConfig *model.ParserConfig

	// flag configuration is loaded
	parserLoaded bool

	// go context
	ctx context.Context
}

// ContextConfig returns an application configuration instance with go context
// attached.
func ContextConfig() (*Config, error) {
	cfg := new(Config)
	cfg.parserLoaded = false

	ok, err := cfg.Parse()

	if !ok {
		util.LogIt("", fmt.Sprintf("Error parsing config: %v", err))
	}
	cfg.ctx = context.Background()

	return cfg, err
}

// ValueOf returns the configuration value for the given parameter identifier.
// The parameter identifier can be any one of:
// 1. a CLI argument without its prefix switch (e.g. '--'; "userdbTestMode")
// 2. an environment variable name (e.g. "CT_USERDB_TEST_MODE"
// 3. a configuration file propert key (e.g. "user.auth.testMode")
func (cfg *Config) ValueOf(id string) string {
	if len(cfg.parameters[id]) > 0 {
		return cfg.parameters[id]
	}
	return ""
}

// ParameterWithDefault returns the configuration value for the given
// enumerated constant or the specified default value if not otherwise
// configured.
func (cfg *Config) ValueOfWithDefault(id string, defVal string) string {
	parmVal := cfg.ValueOf(id)
	if parmVal == model.USER_MUST_PROVIDE {
		parmVal = defVal
	}
	return parmVal
}

// AssignedValue returns true if the configuration parameter for the given
// enumerated constant has a user assigned value.
func (cfg *Config) AssignedValue(id string) bool {
	return (len(cfg.parameters[id]) > 0) && (cfg.parameters[id] != model.USER_MUST_PROVIDE)
}

// IsParsed returns true if configuration has been loaded and parsed.
func (cfg *Config) IsParsed() bool {
	return cfg.parserLoaded
}

// Context returns the go context for the instance.
func (cfg *Config) Context() context.Context {
	return cfg.ctx
}

// Parse initializes configuration parameters with user assigned values or with
// their default values if not assigned a value by the user. User assigned
// values are given the following configuration precedence, from highest to
// lowest:
//  1. CLI arguments
//  2. Environment variables
//  3. Configuration properties
func (cfg *Config) Parse() (bool, error) {
	if !cfg.parserLoaded {
		// Load the parser configuration
		err := cfg.loadConfiguration()
		if err != nil {
			util.LogIt("", fmt.Sprintf("Failed to load parser configuration: %v", err))
			return false, util.WrappedError(err, "loadConfiguration")
		}

		// Set the initial default values with preference for CLI options, followed
		// by env overrides.
		options := util.ParseOptions(cfg.argSwitch, cfg.argSeparator, os.Args[1:])
		var parmVal string
		updateProps := []string{}
		for _, parm := range cfg.parserConfig.Parameters {
			switch {
			case len((*options)[parm.CliArgument]) > 0:
				parmVal = (*options)[parm.CliArgument]
			case len(os.Getenv(parm.EnvironmentVar)) > 0:
				parmVal = os.Getenv(parm.EnvironmentVar)
			case true:
				parmVal = parm.DefaultVal
				updateProps = append(updateProps, parm.PropertyName)
			}
			cfg.parameters[parm.OptionId] = parmVal
		}

		if cfg.AssignedValue(model.APP_CONFIG_ID) || len(model.ApplicationConfigPath) > 0 {
			// ...after which, configuration properties can define any parameters
			// not already assigned.
			//
			_, err := cfg.loadProperties(cfg.ValueOf(model.APP_CONFIG_ID), &updateProps)
			if err != nil {
				success, err := cfg.loadProperties(model.ApplicationConfigPath, &updateProps)
				if err != nil {
					return success, util.WrappedError(err, "loadProperties")
				}
				util.LogIt("", "Loaded application config from override.")
			}
		}

		cfg.parserLoaded = true
	}

	return true, nil
}

// loadConfiguration reads the parser configuration and pre-initializes the
// parser for loading the user's application configuration.
func (cfg *Config) loadConfiguration() error {
	cfg.parserConfig = new(model.ParserConfig)
	if err := util.LoadParserConfig(cfg.parserConfig); err != nil {
		util.LogIt("", fmt.Sprintf("Error loading parser configuration: %v", err))
		return util.WrappedError(err, "LoadParserConfig")
	}
	cfg.parameters = make(map[string]string)

	if len(cfg.parserConfig.ArgSwitch) > 0 {
		cfg.argSwitch = cfg.parserConfig.ArgSwitch
	}

	if len(cfg.parserConfig.ArgSeparator) > 0 {
		switch cfg.parserConfig.ArgSeparator {
		case "EQUALS":
			cfg.argSeparator = '='
		case "COMMA":
			cfg.argSeparator = ','
		case "COLON":
			cfg.argSeparator = ':'
		case "SEMI-COLON":
			cfg.argSeparator = ';'
		case "SPACE":
			// default
		}
	}

	return nil
}

// loadProperties reads all key=value pair properties from the specified file
// path and assigns their values to any matching parameters not already
// assigned a value.
func (cfg *Config) loadProperties(filename string, propsList *[]string) (bool, error) {
	if filename == "" {
		return false, util.WrappedError(errors.New("file name not provided"), "load properties")
	}

	props, err := properties.LoadFile(filename, properties.UTF8)
	if err != nil {
		return false, util.WrappedError(err, "load properties")
	}

	var parmVal string
	for _, parm := range cfg.parserConfig.Parameters {
		if slices.Contains(*propsList, parm.PropertyName) {
			parmVal = configurationValue(props, parm.PropertyName, cfg.parameters[parm.OptionId])
			cfg.parameters[parm.OptionId] = parmVal
			//util.LogIt("", fmt.Sprintf("Updated prop (%v): %v = %v", parm.OptionId, parm.PropertyName, parmVal))
		}
	}

	return true, nil
}

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
