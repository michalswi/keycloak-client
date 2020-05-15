package auth

import (
	"context"
	"log"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

func NewAuthenticator(
	logger *log.Logger,
	clientID string,
	clientSecret string,
	keycloakURL string,
	redirectURL string,
) (*Authenticator, error) {

	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, keycloakURL)
	if err != nil {
		logger.Printf("Failed to get provider, is it running?: %v", err)
		return nil, err
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

	return &Authenticator{
		Provider: provider,
		Config:   oauth2Config,
		Ctx:      ctx,
	}, nil
}
