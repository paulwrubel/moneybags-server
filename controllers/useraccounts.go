package controllers

import (
	"net/http"

	"github.com/paulwrubel/moneybags-server/constants"
	"github.com/paulwrubel/moneybags-server/services"
	log "github.com/sirupsen/logrus"
)

type UserAccounts struct {
	Service services.IUserAccounts
}

type getUserAccountResponse struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Email    *string `json:"email,omitempty"`
}

func (ua *UserAccounts) Get() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userAccount, ok := validateUser(rw, r, ua.Service)
		if !ok {
			return
		}

		writeResponse(rw, http.StatusOK, getUserAccountResponse{
			ID:       userAccount.ID,
			Username: userAccount.Username,
			Email:    userAccount.Email,
		})
	}
}

type postUserAccountRequest struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	Email    *string `json:"email"`
}

type postUserAccountResponse struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Email    *string `json:"email,omitempty"`
}

func (ua *UserAccounts) Post() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var requestBody postUserAccountRequest
		err := unmarshalRequestBody(r.Body, &requestBody)
		if err != nil {
			writeResponse(rw, http.StatusBadRequest, errorsResponseFromErrors(err))
			return
		}

		createdUserAccount, err := ua.Service.Create(requestBody.Username, requestBody.Password, requestBody.Email)
		switch err {
		case constants.ErrUserExists:
			writeResponse(rw, http.StatusConflict, errorsResponseFromErrors(err))
			return
		case constants.ErrInvalidUsername:
			fallthrough
		case constants.ErrInvalidPassword:
			fallthrough
		case constants.ErrInvalidEmail:
			writeResponse(rw, http.StatusBadRequest, errorsResponseFromErrors(err))
			return
		case nil:
			// noop, continue past switch
		default:
			log.WithError(err).Error("Error creating user account")
			writeResponse(rw, http.StatusInternalServerError, errorsResponseFromErrors(err))
			return
		}

		writeResponse(rw, http.StatusCreated, postUserAccountResponse{
			ID:       createdUserAccount.ID,
			Username: createdUserAccount.Username,
			Email:    createdUserAccount.Email,
		})
	}
}
