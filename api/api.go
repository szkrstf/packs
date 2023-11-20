package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/szkrstf/packs"
)

// CalculateHandler is a http.Handler for the calculate endpoint.
type CalculateHandler struct {
	calculator packs.Calculator
}

// NewCalculateHandler creates a new api.CalculateHandler.
func NewCalculateHangler(c packs.Calculator) *CalculateHandler {
	return &CalculateHandler{calculator: c}
}

type calculateReq struct {
	Items int `json:"items"`
}

type calculateResp struct {
	Data map[int]int `json:"data"`
}

// ServeHTTP implements the http.Handler interface.
func (h *CalculateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req calculateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	packs := h.calculator.Calculate(req.Items)

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(calculateResp{Data: packs}); err != nil {
		log.Println(err)
		return
	}
}
