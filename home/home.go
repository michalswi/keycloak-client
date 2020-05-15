package home

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	"github.com/michalswi/keycloak_client/auth"
	"github.com/michalswi/keycloak_client/store"
)

type handlers struct {
	logger        *log.Logger
	state         string
	oidcConfig    *oidc.Config
	authenticator *auth.Authenticator
}

func (h *handlers) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer h.logger.Printf("Home request processed in %s\n", time.Now().Sub(startTime))
		next(w, r)
	}
}

func NewHandlers(logger *log.Logger, state string, oidcConfig *oidc.Config, authenticator *auth.Authenticator) *handlers {
	return &handlers{
		logger:        logger,
		state:         state,
		oidcConfig:    oidcConfig,
		authenticator: authenticator,
	}
}

func (h *handlers) LinkRoutes(mux *mux.Router) {
	mux.HandleFunc("/home", h.Logger(h.AppHandler))
}

func (h *handlers) AppHandler(w http.ResponseWriter, r *http.Request) {

	session, err := store.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rawToken := fmt.Sprintf("%s", session.Values["access_token"])

	_, err = h.authenticator.Provider.Verifier(h.oidcConfig).Verify(context.TODO(), rawToken)
	if err != nil {
		http.Redirect(w, r, h.authenticator.Config.AuthCodeURL(h.state), http.StatusFound)
		return
	}

	message := "keycloak_client"
	version := "0.0.1"
	hostname, err := os.Hostname()
	if err != nil {
		h.logger.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)

	var html = `
	<html>
	<h1>%s</h1>
	<p><b>Hostname</b>: %s; <b>Version</b>: %s</p>
	</html>
	`
	fmt.Fprintf(w, html, message, hostname, version)
}
