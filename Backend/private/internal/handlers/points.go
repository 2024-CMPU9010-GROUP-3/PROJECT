package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/util"
)

type PointsHandler struct{}

func (p *PointsHandler) HandleGetByRadius(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("GET ByRadius")
	w.Header().Set(contentType, applicationJson)
	json.NewEncoder(w).Encode(data)
}

func (p *PointsHandler) HandleGetPointDetails(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("GET PointDetails")
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (p *PointsHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("POST points")
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

func (p *PointsHandler) HandlePut(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("PUT points")
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (p *PointsHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("DELETE points")
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
