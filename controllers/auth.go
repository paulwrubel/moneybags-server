package controllers

import (
	"net/http"

	"github.com/paulwrubel/moneybags-server/constants"
	"github.com/paulwrubel/moneybags-server/services"
	log "github.com/sirupsen/logrus"
)

type Auth struct {
	Service services.IAuth
}

type postLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type postLoginResponse struct {
	Token string `json:"access_token"`
}

func (a *Auth) PostToken() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var requestBody postLoginRequest
		err := unmarshalRequestBody(r.Body, &requestBody)
		if err != nil {
			writeResponse(rw, http.StatusBadRequest, errorsResponseFromErrors(err))
			return
		}

		authenticated, err := a.Service.Authenticate(requestBody.Username, requestBody.Password)
		switch err {
		case constants.ErrUserDoesNotExist:
			writeResponse(rw, http.StatusUnauthorized, errorsResponseFromErrors(err))
			return
		case nil:
			// noop, continue past switch
		default:
			log.WithError(err).Error("Error authenticating user")
			writeResponse(rw, http.StatusInternalServerError, errorsResponseFromErrors(err))
			return
		}
		if !authenticated {
			writeResponse(rw, http.StatusUnauthorized, errorsResponseFromMessages("Invalid username or password"))
			return
		}

		tokenString, err := a.Service.CreateAuthToken(requestBody.Username)
		if err != nil {
			log.WithError(err).Error("Error creating auth token")
			writeResponse(rw, http.StatusInternalServerError, errorsResponseFromErrors(err))
			return
		}

		writeResponse(rw, http.StatusOK, postLoginResponse{
			Token: tokenString,
		})
	}
}
