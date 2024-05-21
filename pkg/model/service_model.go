package model

import "fmt"

var (
	NoError = ServiceError{}

	ClientRequestError  = ServiceError{"A01", "Error sending function request.", nil}
	ClientReadError     = ServiceError{"A02", "Error reading function response.", nil}
	ClientInputError    = ServiceError{"A03", "Error reading input file.", nil}
	ClientOutputError   = ServiceError{"A04", "Error writing to output file.", nil}
	ClientProtocolError = ServiceError{"A05", "Error in service communication.", nil}
	ClientImageError    = ServiceError{"A06", "Error reading image file.", nil}
	ClientError         = ServiceError{"A07", "An internal client error has occurred.", nil}
	CloudStorageError   = ServiceError{"C01", "Error accessing cloud storage.", nil}
	DbQueryError        = ServiceError{"D01", "Error querying user info.", nil}
	DbScanError         = ServiceError{"D02", "Error scanning user info.", nil}
	DbResultsError      = ServiceError{"D03", "Got unknown results error.", nil}
	DbInsertError       = ServiceError{"D04", "Error inserting user info.", nil}
	DbPrepareError      = ServiceError{"D05", "Error preparing statement.", nil}
	DbExecuteError      = ServiceError{"D06", "Error executing statement.", nil}
	DbClientError       = ServiceError{"D07", "Error getting user info client.", nil}
	DbOpenError         = ServiceError{"D08", "Error opening user info.", nil}
	DbPKeyError         = ServiceError{"D09", "Primary key already exists.", nil}
	DbPKeyMissingError  = ServiceError{"D10", "Primary key not found.", nil}
	InvalidKeyError     = ServiceError{"I01", "Incomplete user info.", nil}
	InvalidMsgError     = ServiceError{"I02", "Invalid request message.", nil}
	InternalReadError   = ServiceError{"I03", "Error reading request message.", nil}
	InvalidLoginError   = ServiceError{"I04", "Invalid login credentials provided.", nil}
	InvalidTokenError   = ServiceError{"I05", "Invalid user access token provided.", nil}
	ExpiredTokenError   = ServiceError{"I06", "Expired user access token provided.", nil}
	ImageDecodingError  = ServiceError{"P01", "Error decoding image.", nil}
	SystemError         = ServiceError{"S00", "An internal error has occurred.", nil}
	DatetimeError       = ServiceError{"S01", "A datetime error has occurred.", nil}
	IOError             = ServiceError{"S02", "An input/output error has occurred.", nil}
	UserValidationError = ServiceError{"U01", "User validation period expired.", nil}

	HttpErrorStatus map[string]int
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

func init() {
	HttpErrorStatus = make(map[string]int)
	HttpErrorStatus[CloudStorageError.Code] = 502
	HttpErrorStatus[DbQueryError.Code] = 500
	HttpErrorStatus[DbScanError.Code] = 500
	HttpErrorStatus[DbResultsError.Code] = 500
	HttpErrorStatus[DbInsertError.Code] = 500
	HttpErrorStatus[DbPrepareError.Code] = 500
	HttpErrorStatus[DbExecuteError.Code] = 500
	HttpErrorStatus[DbClientError.Code] = 500
	HttpErrorStatus[DbOpenError.Code] = 500
	HttpErrorStatus[DbPKeyError.Code] = 409
	HttpErrorStatus[DbPKeyMissingError.Code] = 404
	HttpErrorStatus[InvalidKeyError.Code] = 400
	HttpErrorStatus[InvalidMsgError.Code] = 400
	HttpErrorStatus[InternalReadError.Code] = 500
	HttpErrorStatus[InvalidLoginError.Code] = 403
	HttpErrorStatus[InvalidTokenError.Code] = 400
	HttpErrorStatus[ExpiredTokenError.Code] = 403
	HttpErrorStatus[ImageDecodingError.Code] = 500
	HttpErrorStatus[SystemError.Code] = 500
	HttpErrorStatus[DatetimeError.Code] = 500
	HttpErrorStatus[IOError.Code] = 500
	HttpErrorStatus[UserValidationError.Code] = 403
}
