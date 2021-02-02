package delivery

import (
	"medods-test/pkg/models"
	"medods-test/pkg/auth"
	"net/http"
	"encoding/json"
	"context"
	"time"
)

type requestForNewPair struct {
    UserId string `json:"userId"`
}

type requestForRefreshPair struct {
	UserId string `json:"userId"`
    AccessPair models.AccessPair `json:"accessPair"`
}

type handler struct {
	useCase auth.UseCase
}

func newHandler(useCase auth.UseCase) *handler {
	return &handler{
		useCase: useCase,
	}
}

func (h *handler) getNewPair(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req requestForNewPair
	err := json.NewDecoder(r.Body).Decode(&req);
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	inputSession := models.Session{
		UserId: req.UserId,
		UserIP: r.RemoteAddr,
	}
	accessPair, err := h.useCase.GetNewPair(ctx, &inputSession);
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		body, err := json.MarshalIndent(accessPair, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(body))
		}
	}
}

func (h *handler) refreshPair(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var req requestForRefreshPair
	err := json.NewDecoder(r.Body).Decode(&req);
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	inputSession := models.Session{
		UserId: req.UserId,
		UserIP: r.RemoteAddr,
	}
	accessPair, err := h.useCase.RefreshPair(ctx, &req.AccessPair, &inputSession);
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		body, err := json.MarshalIndent(accessPair, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(body))
		}
	}
}