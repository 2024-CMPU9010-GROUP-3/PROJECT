package handlers

import "net/http"

type PointsHandler struct{}

func (p *PointsHandler) HandleGetByRadius(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PointsHandler: GET ByRadius"))
}

func (p *PointsHandler) HandleGetPointDetails(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PointsHandler: GET PointDetails"))
}

func (p *PointsHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PointsHandler: POST"))
}

func (p *PointsHandler) HandlePut(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PointsHandler: PUT"))
}

func (p *PointsHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PointsHandler: Delete"))
}
