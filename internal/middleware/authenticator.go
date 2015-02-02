package middleware

import (
	"encoding/base64"
	"net/http"
	"regexp"
	"strings"
)

type Credentialer interface {
	Credentials() (string, string)
}

type Authenticator struct {
	Handler      http.Handler
	credentialer Credentialer
}

func NewAuthenticator(handler http.Handler, credentialer Credentialer) http.Handler {
	return Authenticator{
		Handler:      handler,
		credentialer: credentialer,
	}
}

func (a Authenticator) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Authorization")
	expression := regexp.MustCompile(`(?i)basic (.*)`)
	regexMatches := expression.FindStringSubmatch(header)
	if len(regexMatches) != 2 {
		a.Fail(w)
		return
	}

	encodedAuth := regexMatches[1]
	decodedAuth, err := base64.StdEncoding.DecodeString(encodedAuth)
	if err != nil {
		a.Fail(w)
		return
	}

	auth := strings.Split(string(decodedAuth), ":")
	if len(auth) != 2 {
		a.Fail(w)
		return
	}

	username, password := a.credentialer.Credentials()
	if username != auth[0] || password != auth[1] {
		a.Fail(w)
		return
	}

	a.Handler.ServeHTTP(w, req)
}

func (a Authenticator) Fail(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}
