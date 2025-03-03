package logger

const (
	// RequestIDKey is unique identifier for request
	RequestIDKey = "request_id"
)

// Tag is key value pair with value in string
type Tag struct {
	Key   string
	Value interface{}
}

// Err to print tag error
func Err(err error) Tag {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	return Tag{
		Key:   "error",
		Value: errMsg,
	}
}
