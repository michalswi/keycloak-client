package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/michalswi/keycloak_client/home"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	"github.com/michalswi/keycloak_client/auth"
	"github.com/michalswi/keycloak_client/callback"
	"github.com/michalswi/keycloak_client/demo"
	"github.com/michalswi/keycloak_client/server"
	"github.com/michalswi/keycloak_client/store"
)

type accessToken struct {
	IDToken string `json:"accessToken"`
}

func main() {

	logger := log.New(os.Stdout, "oidc client ", log.LstdFlags|log.Lshortfile)

	serverAddress := "5050"
	// serverAddress := os.Getenv("SERVICE_ADDR")
	clientID := "demo-client"
	clientSecret := "085f48aa-525c-44ba-a7ba-f388630167cf"
	keycloakURL := "http://localhost:8080/auth/realms/demo"
	redirectURL := "http://localhost:5050/demo/callback"

	// Generate random state
	// non-empty string to protect from CSRF attacks - oauth2Config.AuthCodeURL(state)
	state := randToken(8)

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}

	// auth/auth.go
	authenticator, err := auth.NewAuthenticator(clientID, clientSecret, keycloakURL, redirectURL)
	if err != nil {
		logger.Printf("Authenticator failed: %v", err)
	}

	// store.InitStore()
	r := mux.NewRouter()

	// redirectURL - keycloak related
	c := callback.NewHandlers(logger, state, oidcConfig, authenticator)
	// '/demo' handler
	d := demo.NewHandlers(logger, state, oidcConfig, authenticator)
	// '/home' handler
	h := home.NewHandlers(logger, state, oidcConfig, authenticator)

	c.LinkRoutes(r)
	d.LinkRoutes(r)
	h.LinkRoutes(r)

	// initialize session
	store.InitStore()

	// start server
	srv := server.NewServer(r, serverAddress)
	go func() {
		logger.Printf("Starting server on port: %s.\n", serverAddress)
		err := srv.ListenAndServe()
		if err != nil {
			logger.Fatalf("Server failed to start: %v\n", err)
		}
	}()
	gracefulShutdown(srv, logger)
}

// Generate random token
func randToken(len int) string {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("Generate random state error: %v", err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

// Graceful shutdown
func gracefulShutdown(srv *http.Server, logger *log.Logger) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	logger.Printf("Shutting down the server...\n")
	os.Exit(0)
}
