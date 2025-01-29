package response

import "net/http"

type Response struct {
	Status     bool        `json:"status"`
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Payload    interface{} `json:"data,omitempty"`
}

var (
	generalSuccess = Response{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    "SUCCESS",
	}
	createSuccess = Response{
		Status:     true,
		StatusCode: http.StatusCreated,
		Message:    "CREATED SUCCESS",
	}
)

func GeneralSuccess(message ...string) *Response {
	succ := generalSuccess
	if len(message) != 0 {
		succ.Message = message[0]
	}
	return &succ
}

func GeneralSuccessCustomMessageAndPayload(message string, payload interface{}) *Response {
	succ := generalSuccess
	succ.Message = message
	succ.Payload = payload
	return &succ
}

func CreatedSuccess(message ...string) *Response {
	succ := createSuccess
	if len(message) != 0 {
		succ.Message = message[0]
	}
	return &succ
}

func CreatedSuccessWithPayload(payload interface{}) *Response {
	succ := createSuccess
	succ.Payload = payload
	return &succ
}
