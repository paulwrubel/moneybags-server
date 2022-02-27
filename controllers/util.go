package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/paulwrubel/moneybags-server/constants"
	"github.com/paulwrubel/moneybags-server/models"
	"github.com/paulwrubel/moneybags-server/services"
	log "github.com/sirupsen/logrus"
)

type errorsResponse struct {
	Errors []errorResponse `json:"errors"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func validateUser(rw http.ResponseWriter, r *http.Request, uaService services.IUserAccounts) (*models.UserAccount, bool) {
	username, ok := r.Context().Value(constants.UsernameContextKey).(string)
	if !ok {
		log.Error("Could not retrieve username from context")
		writeResponse(rw, http.StatusInternalServerError, errorsResponseFromMessages("User not found in provided token. Please contact the site administrator"))
		return nil, false
	}

	userExists, err := uaService.ExistsByUsername(username)
	if err != nil {
		log.Errorf("Error checking if user exists: %s", err.Error())
		writeResponse(rw, http.StatusInternalServerError, errorsResponseFromErrors(err))
		return nil, false
	}
	if !userExists {
		writeResponse(rw, http.StatusNotFound, errorsResponseFromMessages("User does not exist"))
		return nil, false
	}

	userAccount, err := uaService.GetInfo(username)
	if err != nil {
		log.Errorf("Error getting user account: %s", err.Error())
		writeResponse(rw, http.StatusInternalServerError, errorsResponseFromErrors(err))
		return nil, false
	}
	return userAccount, true
}

func unmarshalRequestBody(body io.Reader, dst interface{}) error {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("error reading request body: %w", err)
	}
	err = json.Unmarshal(bodyBytes, dst)
	if err != nil {
		return fmt.Errorf("error unmarshalling request body: %w", err)
	}
	return nil
}

func writeResponse(rw http.ResponseWriter, status int, response interface{}) {
	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.WithError(err).Error("Error marshalling response")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)
	_, err = rw.Write(responseBytes)
	if err != nil {
		log.WithError(err).Error("Error writing response")
	}
}

func errorsResponseFromErrors(errs ...error) errorsResponse {
	errorsResponse := errorsResponse{}
	for _, err := range errs {
		errorsResponse.Errors = append(errorsResponse.Errors, errorResponse{
			Message: err.Error(),
		})
	}
	return errorsResponse
}

func errorsResponseFromMessages(msgs ...string) errorsResponse {
	errorsResponse := errorsResponse{}
	for _, msg := range msgs {
		errorsResponse.Errors = append(errorsResponse.Errors, errorResponse{
			Message: msg,
		})
	}
	return errorsResponse
}
