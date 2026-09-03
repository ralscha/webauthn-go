package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexedwards/scs/v2"
)

const transactionTestDriverName = "webauthn-transaction-test"

func init() {
	sql.Register(transactionTestDriverName, transactionTestDriver{})
}

type transactionTestDriver struct{}

func (transactionTestDriver) Open(name string) (driver.Conn, error) {
	return &transactionTestConnection{commitFails: name == "commit-fails"}, nil
}

type transactionTestConnection struct {
	commitFails bool
}

func (c *transactionTestConnection) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("not implemented")
}

func (c *transactionTestConnection) Close() error { return nil }

func (c *transactionTestConnection) Begin() (driver.Tx, error) {
	return &transactionTestTransaction{commitFails: c.commitFails}, nil
}

type transactionTestTransaction struct {
	commitFails bool
}

func (tx *transactionTestTransaction) Commit() error {
	if tx.commitFails {
		return errors.New("commit failed")
	}
	return nil
}

func (*transactionTestTransaction) Rollback() error { return nil }

func TestRWTransactionDoesNotSendSuccessBeforeCommit(t *testing.T) {
	db, err := sql.Open(transactionTestDriverName, "commit-fails")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	app := &application{database: db}
	handler := app.rwTransaction(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "success")
	}))

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/", nil))

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	if recorder.Body.String() == "success" {
		t.Fatal("success response was sent before the transaction committed")
	}
}

func TestRWTransactionForwardsResponseAfterCommit(t *testing.T) {
	db, err := sql.Open(transactionTestDriverName, "commit-succeeds")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	app := &application{database: db}
	handler := app.rwTransaction(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Context().Value(transactionKey).(*sql.Tx); !ok {
			t.Error("transaction is missing from request context")
		}
		w.Header().Set("X-Test", "value")
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, "created")
	}))

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/", nil))

	if recorder.Code != http.StatusCreated || recorder.Body.String() != "created" {
		t.Fatalf("response = %d %q, want %d %q", recorder.Code, recorder.Body.String(), http.StatusCreated, "created")
	}
	if recorder.Header().Get("X-Test") != "value" {
		t.Fatal("response headers were not forwarded")
	}
}

func TestAuthenticatedOnlyReturnsUnauthorized(t *testing.T) {
	sm := scs.New()
	app := &application{sessionManager: sm}
	handler := sm.LoadAndSave(app.authenticatedOnly(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("protected handler was called")
	})))

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil))

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}
