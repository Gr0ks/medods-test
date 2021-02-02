package delivery

import (
	"github.com/gorilla/mux"
	"medods-test/pkg/auth"
)

func RegisterHTTPEndpoints(router *mux.Router, useCase auth.UseCase) {
	h := newHandler(useCase)

	router.HandleFunc("/api/auth", h.getNewPair)
	router.HandleFunc("/api/auth/refreshToken", h.refreshPair)
}