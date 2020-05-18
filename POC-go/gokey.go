package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type accessToken struct {
	IDToken string `json:"accessToken"`
}

func main() {

	keycloakURL := "http://localhost:8080/auth/realms/demo"

	var clientSecret string
	if len(os.Args) < 2 {
		log.Fatalf("Missing 'clientSecret' for demo-client")
	} else {
		clientSecret = os.Args[1]
	}

	clientID := "demo-client"
	redirectURL := "http://localhost:5050/demo/callback"

	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, keycloakURL)
	if err != nil {
		log.Fatalf("Failed to get provider: %v", err)
	}

	// Configure an OpenID Connect aware OAuth2 client
	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		// Discovery returns the OAuth2 endpoints
		Endpoint: provider.Endpoint(),
		// "openid" is a required scope for OpenID Connect flows
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	// Generate random state
	// non-empty string to protect from CSRF attacks - oauth2Config.AuthCodeURL(state)
	state := randToken(8)

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier := provider.Verifier(oidcConfig)

	// App Route Handler
	http.HandleFunc("/demo", func(w http.ResponseWriter, r *http.Request) {

		rawAccessToken := r.Header.Get("Authorization")
		if rawAccessToken == "" {
			http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
			return
		}

		parts := strings.Split(rawAccessToken, " ")
		if len(parts) != 2 {
			w.WriteHeader(400)
			return
		}

		_, err := verifier.Verify(ctx, parts[1])
		if err != nil {
			http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
			return
		}

		w.Write([]byte("go oidc client works"))
	})

	// OAuth2 Callback Handler
	http.HandleFunc("/demo/callback", func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Query().Get("state") != state {
			log.Printf("State did not match")
			http.Error(w, "State did not match", http.StatusBadRequest)
			return
		}

		// Verify state and errors
		oauth2Token, err := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
		if err != nil {
			log.Printf("Failed to exchange token, token not found?: " + err.Error())
			http.Error(w, "Failed to exchange token, token not found?: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Extract the ID Token from OAuth2 token
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			log.Printf("No id_token field in OAuth2 token.")
			http.Error(w, "No id_token field in OAuth2 token.", http.StatusInternalServerError)
			return
		}

		// 'rawIDToken' is displayed twice (two different, why?)
		// 'rawIDToken' is the proper Token to access App
		// log.Println(rawIDToken)

		// Parse and verify ID Token payload
		idToken, err := verifier.Verify(ctx, rawIDToken)
		if err != nil {
			log.Printf("Failed to verify ID Token: " + err.Error())
			http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Get the userInfo and extract custom claims
		// 1 - return json with token + token info
		var resp map[string]interface{}

		if err := idToken.Claims(&resp); err != nil {
			log.Printf("Claims failed >> " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp["access_token"] = rawIDToken
		resp["token_details"] = oauth2Token
		// resp["some_token"] = oauth2Token.AccessToken

		// 2 - return only token
		// resp := accessToken{
		// 	IDToken: rawIDToken,
		// }

		// if err := idToken.Claims(&resp); err != nil {
		// 	log.Printf("Claims failed >> " + err.Error())
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// common for 1 and 2
		data, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Marshalling failed >> " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)

	})
	log.Printf("Starting server on port 5050")
	log.Fatal(http.ListenAndServe(":5050", nil))
}

func randToken(len int) string {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("Generate random state error: %v", err)
	}
	return base64.StdEncoding.EncodeToString(b)
}
