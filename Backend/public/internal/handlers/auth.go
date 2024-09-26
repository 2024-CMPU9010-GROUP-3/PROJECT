package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/util"
)

type AuthHandler struct{}

func (p *AuthHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("GET User")
	w.Header().Set(contentType, applicationJson)
	err := json.NewEncoder(w).Encode(data)
	util.CheckResponseError(err, w)
}

func (p *AuthHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("POST User")
	w.Header().Set(contentType, applicationJson)
	err := json.NewEncoder(w).Encode(data)
	util.CheckResponseError(err, w)
}

func (p *AuthHandler) HandlePut(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("PUT User")
	w.Header().Set(contentType, applicationJson)
	err := json.NewEncoder(w).Encode(data)
	util.CheckResponseError(err, w)
}

func (p *AuthHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("DELETE User")
	w.Header().Set(contentType, applicationJson)
	err := json.NewEncoder(w).Encode(data)
	util.CheckResponseError(err, w)
}

func (p *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("POST User/Login")
	w.Header().Set(contentType, applicationJson)
	err := json.NewEncoder(w).Encode(data)
	util.CheckResponseError(err, w)
}
