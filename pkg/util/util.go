package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/efficientgo/core/errors"

	"Cloudtacts/pkg/model"
)

func LoadUserListFile(path string, userList *model.UserList) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		LogIt("", fmt.Sprintf("Error opening users file %v.", path))
		return WrappedError(err, "open file")
	}
	defer jsonFile.Close()

	bbuff, err := io.ReadAll(jsonFile)
	if err != nil {
		LogIt("", fmt.Sprintf("Error reading users file %v.", path))
		return WrappedError(err, "read file")
	}

	err = ToUserList(bbuff, userList)
	if err != nil {
		LogIt("", fmt.Sprintf("Error unmarshalling users file %v to JSON.", path))
		return WrappedError(err, "unmarshal")
	}

	return nil
}

func LoadParserConfig(parms *model.ParserConfig) error {
	jsonFile, err := os.Open(model.ParserConfigPath)
	if err != nil {
		LogIt("", fmt.Sprintf("Error opening parser config file %v.", model.ParserConfigPath))
		return WrappedError(err, "open file")
	}
	defer jsonFile.Close()

	bbuff, err := io.ReadAll(jsonFile)
	if err != nil {
		LogIt("", fmt.Sprintf("Error reading parser config file %v.", model.ParserConfigPath))
		return WrappedError(err, "read file")
	}

	err = json.Unmarshal(bbuff, parms)
	if err != nil {
		LogIt("", fmt.Sprintf("Error unmarshalling parser config file %v to JSON.", model.ParserConfigPath))
		return WrappedError(err, "unmarshal")
	}

	return nil
}

func ImageFileType(fileName string) string {
	var itype string

	nameParts := strings.Split(strings.ToLower(fileName), ".")
	switch nameParts[len(nameParts)-1] {
	case "gif":
		itype = "gif"
	case "jpg", "jpeg":
		itype = "jpg"
	case "png":
		itype = "png"
	default:
		// TODO: try an determine from image mime.type
		itype = "unk"
	}

	return itype
}

func ToUserList(data []byte, userList *model.UserList) error {
	err := json.Unmarshal(data, userList)
	if err != nil {
		LogIt("", fmt.Sprintf("Error converting data to UserList:\n%v", data))
		return WrappedError(err, "unmarshal")
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
			LogIt("", fmt.Sprintf("Argument separator not found: '%v'", string(argSeparator)))
		}
	}

	return &opts
}

func StripDateStamp(datetime string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(datetime, "-", ""),
				":", ""),
			"T", ""),
		"Z", "")
}

func ToDatetime(tstamp string) (time.Time, model.ServiceError) {
	var datetime time.Time
	serr := model.NoError

	var year, month, day, hour, minute, second int
	var err error
	if year, err = strconv.Atoi(tstamp[0:4]); err != nil {
		serr = model.DatetimeError.WithCause(err)
	}
	if month, err = strconv.Atoi(tstamp[4:6]); err != nil && !serr.IsError() {
		serr = model.DatetimeError.WithCause(err)
	}
	if day, err = strconv.Atoi(tstamp[6:8]); err != nil && !serr.IsError() {
		serr = model.DatetimeError.WithCause(err)
	}
	if hour, err = strconv.Atoi(tstamp[8:10]); err != nil && !serr.IsError() {
		serr = model.DatetimeError.WithCause(err)
	}
	if minute, err = strconv.Atoi(tstamp[10:12]); err != nil && !serr.IsError() {
		serr = model.DatetimeError.WithCause(err)
	}
	if second, err = strconv.Atoi(tstamp[12:14]); err != nil && !serr.IsError() {
		serr = model.DatetimeError.WithCause(err)
	}
	if !serr.IsError() {
		tmonth, serr := ToMonth(month)
		if !serr.IsError() {
			datetime = time.Date(year, tmonth, day, hour, minute, second, 0, time.UTC)
		}
	}

	return datetime, serr
}

func ToMonth(monthNum int) (time.Month, model.ServiceError) {
	if monthNum < 1 || monthNum > 12 {
		LogIt("Cloudtacts", fmt.Sprintf("Received invalid month number for conversion: %d", monthNum))
		return -1, model.DatetimeError
	}

	month := []time.Month{time.January, time.February, time.March, time.April, time.May, time.June,
		time.July, time.August, time.September, time.October, time.November, time.December}[monthNum-1]

	return month, model.NoError
}

func WrappedError(err error, tag string) error {
	return errors.Wrap(err, tag)
}

func LogIt(tag string, message string) {
	if len(tag) > 0 {
		ct_log.Printf("%v - %v", tag, message)
	} else {
		ct_log.Printf("%v", message)
	}
}

func LogError(tag string, message string, cause error) {
	if len(tag) > 0 {
		ct_log.Fatalf("%v - %v\n%v", tag, message, cause)
	} else {
		ct_log.Fatalf("%v\n%v", message, cause)
	}
}

var ct_log *log.Logger

func init() {
	ct_log = log.Default()
}
