package handlers

import "net/http"

type AuthHandler struct{}

func (p *AuthHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("AuthHandler: GET"))
}

func (p *AuthHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("AuthHandler: POST"))
}

func (p *AuthHandler) HandlePut(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("AuthHandler: PUT"))
}

func (p *AuthHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("AuthHandler: DELETE"))
}

func (p *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("AuthHandler: POST login"))
}
