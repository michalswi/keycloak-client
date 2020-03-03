package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func NewServer(r *mux.Router, serverAddress string) *http.Server {
	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + serverAddress,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return srv
}
