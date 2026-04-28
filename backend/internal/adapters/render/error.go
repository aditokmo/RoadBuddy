package render

import "net/http"

type errorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func Error(w http.ResponseWriter, status int, message, code string) {
	JSON(w, status, errorResponse{
		Error: message,
		Code:  code,
	})
}
