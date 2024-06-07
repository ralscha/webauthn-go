package main

import (
	"net/http"
	"webauthn.rasc.ch/internal/response"
)

func (app *application) authenticateHandler(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.GetInt(r.Context(), "userID")
	if userID > 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.sessionManager.Destroy(r.Context()); err != nil {
		response.InternalServerError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
