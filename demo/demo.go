package demo

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	"github.com/michalswi/keycloak-client/auth"
)

type handlers struct {
	logger        *log.Logger
	state         string
	oidcConfig    *oidc.Config
	authenticator *auth.Authenticator
}

func (d *handlers) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer d.logger.Printf("Demo request processed in %s\n", time.Now().Sub(startTime))
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

func (d *handlers) LinkRoutes(mux *mux.Router) {
	mux.HandleFunc("/demo", d.Logger(d.AppHandler))
}

func (d *handlers) AppHandler(w http.ResponseWriter, r *http.Request) {

	rawAccessToken := r.Header.Get("Authorization")
	if rawAccessToken == "" {
		http.Redirect(w, r, d.authenticator.Config.AuthCodeURL(d.state), http.StatusFound)
		return
	}

	parts := strings.Split(rawAccessToken, " ")
	if len(parts) != 2 {
		w.WriteHeader(400)
		return
	}

	_, err := d.authenticator.Provider.Verifier(d.oidcConfig).Verify(context.TODO(), parts[1])
	if err != nil {
		http.Redirect(w, r, d.authenticator.Config.AuthCodeURL(d.state), http.StatusFound)
		return
	}

	w.Write([]byte("go oidc client works"))
}
