package model

import "fmt"

var (
	NoError = ServiceError{}

	DbQueryError      = ServiceError{"D01", "Error querying user info.", nil}
	DbScanError       = ServiceError{"D02", "Error scanning user info.", nil}
	DbResultsError    = ServiceError{"D03", "Got unknown results error.", nil}
	DbInsertError     = ServiceError{"D04", "Error inserting user info.", nil}
	DbPrepareError    = ServiceError{"D05", "Error preparing statement.", nil}
	DbExecuteError    = ServiceError{"D06", "Error executing statement.", nil}
	DbClientError     = ServiceError{"D07", "Error getting user info client.", nil}
	DbOpenError       = ServiceError{"D08", "Error opening user info.", nil}
	DbPKeyError       = ServiceError{"D09", "Primary key already exists.", nil}
	InvalidKeyError   = ServiceError{"I01", "Incomplete user info.", nil}
	InvalidMsgError   = ServiceError{"I02", "Invalid request message.", nil}
	InternalReadError = ServiceError{"I03", "Error reading request message.", nil}
	SystemError       = ServiceError{"S00", "An internal error has occurred.", nil}
	DatetimeError     = ServiceError{"S01", "A datetime error has occurred.", nil}
)

// ServiceError represents a base error result
type ServiceError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"causedBy"`
}

func (err ServiceError) Error() string {
	if err.Cause == nil {
		return fmt.Sprintf("%v: %v", err.Code, err.Message)
	} else {
		return fmt.Sprintf("%v: %v\n%v", err.Code, err.Message, err.Cause)
	}
}

func (err ServiceError) IsError() bool {
	return err.Code != ""
}

func (err ServiceError) WithCause(src error) ServiceError {
	return ServiceError{err.Code, err.Message, src}
}
