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
	err := json.NewEncoder(w).Encode(data)
	util.CheckResponseError(err, w)
}

func (p *PointsHandler) HandleGetPointDetails(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("GET PointDetails")
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	util.CheckResponseError(err, w)
}
