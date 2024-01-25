package controllers

import (
	"encoding/json"
	"net/http"
)

type Result struct {
	Data interface{} `json:"results"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func HttpResponse(w http.ResponseWriter, err error, data interface{}, code int) {
	if err != nil {
		if code == 0 {
			code = http.StatusBadRequest
		}
		respondWithError(w, code, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, data)
}

func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	result := Result{
		Data: data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	result := ErrorResponse{
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
