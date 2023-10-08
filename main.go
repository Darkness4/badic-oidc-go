package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type OIDCClaims struct {
	jwt.RegisteredClaims
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	_ = godotenv.Load()
	ctx := context.Background()
	redirectURL := "http://localhost:3000/callback"
	listenAddress := ":3000"
	oidcIssuerURL := os.Getenv("OIDC_ISSUER")
	clientSecret := os.Getenv("CLIENT_SECRET")
	clientID := os.Getenv("CLIENT_ID")

	// Fetch and parse the discovery document
	provider, err := oidc.NewProvider(ctx, oidcIssuerURL)
	if err != nil {
		panic(err)
	}

	// Initialize OAuth2
	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	// Login page
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// Fetch csrf
		token := csrf.Token(r)
		cookie := &http.Cookie{
			Name:     "csrf_token",
			Value:    token,
			Expires:  time.Now().Add(1 * time.Minute), // Set expiration time as needed
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)

		// Redirect to authorization grant page
		http.Redirect(w, r, oauth2Config.AuthCodeURL(token), http.StatusFound)
	})

	// Callback: check code with OAuth2 authorization server
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// Fetch code
		val, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		code := val.Get("code")

		// Check CSRF
		csrfToken := val.Get("state")
		expectedCSRF, err := r.Cookie("csrf_token")
		if err == http.ErrNoCookie {
			http.Error(w, "no csrf cookie error", http.StatusUnauthorized)
			return
		}
		if csrfToken != expectedCSRF.Value {
			http.Error(w, "csrf error", http.StatusUnauthorized)
			return
		}

		// Fetch accessToken
		oauth2Token, err := oauth2Config.Exchange(r.Context(), code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Fetch id token and authenticate
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "missing ID token", http.StatusUnauthorized)
			return
		}

		idToken, err := provider.VerifierContext(
			r.Context(),
			&oidc.Config{
				ClientID: clientID,
			}).Verify(ctx, rawIDToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims := OIDCClaims{}
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Do something with access token and id token (store as session cookie, make API calls, redirect to user profile...). Your user is now authenticated.
		idTokenJSON, _ := json.MarshalIndent(claims, "", "  ")
		fmt.Fprintf(
			w,
			"Access Token: %s\nDecoded ID Token: %s\n",
			oauth2Token.AccessToken,
			idTokenJSON,
		)
	})

	// Callback page
	csrfKey := []byte("random-secret")
	slog.Info("listening", slog.String("address", listenAddress))
	if err := http.ListenAndServe(listenAddress, csrf.Protect(csrfKey)(http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
}
