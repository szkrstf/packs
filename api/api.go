package api

import (
	"encoding/json"
	"errors"
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

// SizeHandler is a http.Handler for the sizes endpoint.
type SizeHandler struct {
	store packs.SizeStore
}

// NewSizeHandler creates a new ui.SizeHandler
func NewSizeHandler(store packs.SizeStore) *SizeHandler {
	return &SizeHandler{store: store}
}

type sizeReq struct {
	Sizes []int `json:"sizes"`
}

type sizeResp struct {
	Data []int `json:"data"`
}

// ServeHTTP implements the http.Handler interface.
func (h *SizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		sizes := h.store.Get()
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(sizeResp{Data: sizes}); err != nil {
			log.Println(err)
		}
	case http.MethodPost:
		var req sizeReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if err := h.store.Set(req.Sizes); errors.Is(err, packs.ErrInvalidSizes) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
