package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/util"
)


func (p *PointsHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("POST points")
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(data)
	util.CheckResponseError(err, w)
}

func (p *PointsHandler) HandlePut(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("PUT points")
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	util.CheckResponseError(err, w)
}

func (p *PointsHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("DELETE points")
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	util.CheckResponseError(err, w)
}
