package util

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"Cloudtacts/pkg/model"
)

var ParserConfigPath string = model.PARSER_CONFIG_FILE

func LoadUserListFile(path string, userList *model.UserList) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error opening users file %v.", path))
		return err
	}
	defer jsonFile.Close()

	bbuff, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error reading users file %v.", path))
		return err
	}

	err = json.Unmarshal(bbuff, userList)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error unmarshalling users file %v to JSON.", path))
		return err
	}

	return nil
}

func LoadParserConfig(parms *model.ParserConfig) error {
	jsonFile, err := os.Open(ParserConfigPath)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error opening parser config file %v.", ParserConfigPath))
		return err
	}
	defer jsonFile.Close()

	bbuff, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error reading parser config file %v.", ParserConfigPath))
		return err
	}

	err = json.Unmarshal(bbuff, parms)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error unmarshalling parser config file %v to JSON.", ParserConfigPath))
		return err
	}

	return nil
}

func ParseOptions(argSwitch string, argSeparator uint8, args []string) *map[string]string {
	opts := make(map[string]string)

	var v []string
	for idx, arg := range args {
		if argSeparator == ' ' {
			if idx%2 == 1 {
				opts[v[0]] = arg
			} else {
				v[0] = arg
				opts[arg] = ""
			}
		} else if strings.Contains(arg, string(argSeparator)) {
			v = strings.Split(strings.TrimPrefix(arg, argSwitch), string(argSeparator))
			opts[v[0]] = v[1]
		} else {
			// TODO: log warning?
		}
	}

	return &opts
}

func DigestOf(data string) []byte {
	hashValue := sha256.Sum256([]byte(data))
	return hashValue[:]
}

func LogIt(message string) {
	log.Println(fmt.Sprintf("Cloudtacts - %v", message))
}

func LogError(message string, cause error) {
	log.Fatal(fmt.Sprintf("Cloudtacts - %v\n%v", message, cause))
}
