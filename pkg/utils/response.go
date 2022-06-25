package utils

import "strings"

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

type EmptyObj struct{}

func BuildResponse(success bool, message string, data interface{}) Response {
	return Response{
		Success: success,
		Message: message,
		Errors:  EmptyObj{},
		Data:    data,
	}
}

func BuildErrorResponse(message string, errors string, data interface{}) Response {
	splittedError := strings.Split(errors, "\n")
	return Response{
		Success: false,
		Message: message,
		Errors:  splittedError,
		Data:    data,
	}
}
