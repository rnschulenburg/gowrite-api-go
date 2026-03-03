package Http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type ReturnObject struct {
	Status        string      `json:"status"`
	StatusMessage string      `json:"statusMessage"`
	Data          interface{} `json:"data"`
	ResponseTime  time.Time   `json:"responseTime"`
}

func JsonResponse(w http.ResponseWriter, err error, statusCode string, statusMessage string, data interface{}) {
	if err != nil {
		statusNumber, _ := strconv.ParseInt(statusCode, 10, 16)
		JsonErrorResponse(w, int(statusNumber), err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err4 := json.NewEncoder(w).Encode(
		ReturnObject{
			Status:        statusCode,
			StatusMessage: statusMessage,
			Data:          data,
			ResponseTime:  time.Now()})
	if err4 != nil {
		return
	}

}

type ReturnErrorObject struct {
	Error        string      `json:"error"`
	ErrorMessage interface{} `json:"errorMessage"`
	ResponseTime time.Time   `json:"responseTime"`
}

func JsonErrorResponse(w http.ResponseWriter, code int, message interface{}) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err4 := json.NewEncoder(w).Encode(
		ReturnErrorObject{
			Error:        http.StatusText(code),
			ErrorMessage: message,
			ResponseTime: time.Now()})
	if err4 != nil {
		return
	}
}
