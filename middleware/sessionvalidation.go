package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/paulwrubel/moneybags-server/constants"
	"github.com/paulwrubel/moneybags-server/services"
	log "github.com/sirupsen/logrus"
)

func SessionValidation(authService services.IAuth) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Debug("No authorization header found")
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			authHeaderParts := strings.Split(authHeader, " ")
			if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
				log.Debug("Invalid authorization header")
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			tokenString := authHeaderParts[1]
			token, err := authService.ValidateSession(tokenString)
			if err != nil {
				log.WithError(err).Error("Error validating session")
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}
			log.Debug("Session validated")

			// set context
			ctx := r.Context()
			if usernameClaim, ok := token.Claims.(jwt.MapClaims)["sub"].(string); ok {
				ctx = context.WithValue(ctx, constants.UsernameContextKey, usernameClaim)
			}

			log.WithField("username", ctx.Value(constants.UsernameContextKey)).Debug("Context set")

			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
