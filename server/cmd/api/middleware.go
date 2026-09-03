package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"webauthn.rasc.ch/internal/request"
	"webauthn.rasc.ch/internal/response"
)

func limitRequestBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, request.MaxBodyBytes)
		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticatedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := app.sessionManager.GetInt(r.Context(), authenticatedUserIDKey)
		if userID > 0 {
			next.ServeHTTP(w, r)
		} else {
			response.Unauthorized(w)
		}
	})
}

type contextKey string

const (
	transactionKey contextKey = "transaction"
)

func (app *application) rwTransaction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, err := app.database.BeginTx(r.Context(), nil)
		if err != nil {
			response.InternalServerError(w, err)
			return
		}

		committed := false
		defer func() {
			if !committed {
				_ = tx.Rollback()
			}
		}()

		ctx := context.WithValue(r.Context(), transactionKey, tx)
		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r.WithContext(ctx))

		if recorder.Code >= http.StatusBadRequest {
			writeRecordedResponse(w, recorder)
			return
		}

		if err := tx.Commit(); err != nil {
			response.InternalServerError(w, err)
			return
		}
		committed = true
		writeRecordedResponse(w, recorder)
	})
}

func writeRecordedResponse(w http.ResponseWriter, recorder *httptest.ResponseRecorder) {
	for key, values := range recorder.Header() {
		w.Header()[key] = append([]string(nil), values...)
	}
	w.WriteHeader(recorder.Code)
	if recorder.Body.Len() == 0 {
		return
	}
	if _, err := io.Copy(w, recorder.Body); err != nil {
		slog.Error("writing buffered response failed", "error", err)
	}
}
