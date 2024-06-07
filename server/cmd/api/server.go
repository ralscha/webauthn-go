package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         app.config.HTTP.Addr,
		Handler:      app.routes(),
		ReadTimeout:  time.Duration(app.config.HTTP.ReadTimeoutInSeconds) * time.Second,
		WriteTimeout: time.Duration(app.config.HTTP.WriteTimeoutInSeconds) * time.Second,
		IdleTimeout:  time.Duration(app.config.HTTP.IdleTimeoutInSeconds) * time.Second,
	}

	shutdownErrorChan := make(chan error)

	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		shutdownErrorChan <- srv.Shutdown(ctx)
	}()

	slog.Info("starting server", "address", srv.Addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErrorChan
	if err != nil {
		return err
	}

	slog.Info("stopped server", "address", srv.Addr)

	shutdownChan := app.taskScheduler.Shutdown()
	<-shutdownChan

	app.wg.Wait()
	return nil
}

func (app *application) backgroundTask(fn func()) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()
		defer func() {
			if err := recover(); err != nil {
				if e, ok := err.(error); ok {
					slog.Error("background job failed", "error", e)
				} else {
					slog.Error("background job failed", "error", err)
				}
			}
		}()

		fn()
	}()
}
