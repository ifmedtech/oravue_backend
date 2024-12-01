package response

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOK    = "ok"
	StatusError = "error"
)

func WriteJson(w http.ResponseWriter, code int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	encodedData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.WriteHeader(code)
	_, writeErr := w.Write(encodedData)
	return writeErr
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  fmt.Sprintf("an unexpected error occurred: %s", err.Error()),
	}
}

func ValidationError(err validator.ValidationErrors) Response {
	var errMsgs []string
	for _, e := range err {
		switch e.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field '%s' is required", e.Field()))
		case "numeric":
			errMsgs = append(errMsgs, fmt.Sprintf("field '%s' must be numeric", e.Field()))
		case "len":
			errMsgs = append(errMsgs, fmt.Sprintf("field '%s' must be exactly %s characters", e.Field(), e.Param()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field '%s' is invalid", e.Field()))
		}
	}
	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
