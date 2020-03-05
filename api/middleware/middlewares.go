package middleware

import (
	"errors"
	"net/http"

	"github.com/nitomibimo/CRUD-GO/api/auth"
	"github.com/nitomibimo/CRUD-GO/api/response"
)

// SetMiddlewareJSON - set header response
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}

}

// SetMiddlewareAuthentication - set middleware auth
func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}

}
