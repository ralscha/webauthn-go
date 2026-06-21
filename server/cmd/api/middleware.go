package main

import (
	"context"
	"fmt"
	"net/http"
	"webauthn.rasc.ch/internal/response"
)

func (app *application) authenticatedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := app.sessionManager.GetInt(r.Context(), "userID")
		if userID > 0 {
			next.ServeHTTP(w, r)
		} else {
			response.Forbidden(w)
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
			fmt.Println("rwTransaction: BeginTx failed")
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
		recorder := &statusRecorder{ResponseWriter: w}
		next.ServeHTTP(recorder, r.WithContext(ctx))

		if recorder.status >= http.StatusBadRequest {
			return
		}

		if err := tx.Commit(); err != nil {
			fmt.Println("Rolling back transaction")
			response.InternalServerError(w, err)
			return
		}
		committed = true
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	if r.status == 0 {
		r.status = status
		r.ResponseWriter.WriteHeader(status)
	}
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(b)
}
