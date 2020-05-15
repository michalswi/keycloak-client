package callback

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	"github.com/michalswi/keycloak-client/auth"
	"github.com/michalswi/keycloak-client/store"
)

type handlers struct {
	logger        *log.Logger
	state         string
	oidcConfig    *oidc.Config
	authenticator *auth.Authenticator
}

func (c *handlers) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer c.logger.Printf("Callback request processed in %s\n", time.Now().Sub(startTime))
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

func (c *handlers) LinkRoutes(mux *mux.Router) {
	mux.HandleFunc("/demo/callback", c.Logger(c.CallbackHandler))
}

// OAuth2 Callback Handler
func (c *handlers) CallbackHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Query().Get("state") != c.state {
		c.logger.Printf("State did not match")
		http.Error(w, "State did not match", http.StatusBadRequest)
		return
	}

	// Verify state and errors
	oauth2Token, err := c.authenticator.Config.Exchange(context.TODO(), r.URL.Query().Get("code"))
	if err != nil {
		c.logger.Printf("Failed to exchange token, token not found?: " + err.Error())
		http.Error(w, "Failed to exchange token, token not found?: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Extract the ID Token from OAuth2 token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		c.logger.Printf("No id_token field in OAuth2 token.")
		http.Error(w, "No id_token field in OAuth2 token.", http.StatusInternalServerError)
		return
	}

	// 'rawIDToken' is displayed twice (two different, why?)
	// 'rawIDToken' is the proper Token to access App
	// log.Println(rawIDToken)

	// Parse and verify ID Token payload
	idToken, err := c.authenticator.Provider.Verifier(c.oidcConfig).Verify(context.TODO(), rawIDToken)
	if err != nil {
		c.logger.Printf("Failed to verify ID Token: " + err.Error())
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the userInfo and extract custom claims
	var resp map[string]interface{}

	if err := idToken.Claims(&resp); err != nil {
		c.logger.Printf("Claims failed >> " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp["access_token"] = rawIDToken
	resp["token_details"] = oauth2Token
	// resp["some_token"] = oauth2Token.AccessToken

	// 1 for '/demo' - return JSON with TOKEN + TOKEN info then using TOKEN:
	// curl -i -XGET -H "Authorization: Bearer $TOKEN" localhost:5050/demo

	// data, err := json.Marshal(resp)
	// if err != nil {
	// 	log.Printf("Marshalling failed >> " + err.Error())
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// w.Write(data)

	// 2 for '/home' - get token, keep session and redirect to '/home'
	session, err := store.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["access_token"] = rawIDToken
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.logger.Printf("Redirect to '/home'")
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
