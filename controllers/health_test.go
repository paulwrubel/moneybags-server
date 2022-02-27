package controllers_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/paulwrubel/moneybags-server/controllers"
	"github.com/stretchr/testify/assert"
)

func TestHealthGet(t *testing.T) {
	tests := []struct {
		name                 string
		endpoint             string
		requestMethod        string
		requestBody          string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "get",
			endpoint:             "/health",
			requestMethod:        http.MethodGet,
			requestBody:          ``,
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			b := &controllers.Health{}

			rw := httptest.NewRecorder()
			r := httptest.NewRequest(tt.requestMethod, tt.endpoint, strings.NewReader(tt.requestBody))

			b.Get().ServeHTTP(rw, r)

			assert.Equal(t, tt.expectedStatusCode, rw.Result().StatusCode)
			resBody, _ := io.ReadAll(rw.Result().Body)
			if json.Valid(resBody) && json.Valid([]byte(tt.expectedResponseBody)) {
				assert.JSONEq(t, tt.expectedResponseBody, string(resBody))
			} else {
				assert.Equal(t, tt.expectedResponseBody, string(resBody))
			}
		})
	}
}
