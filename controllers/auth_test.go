package controllers_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/paulwrubel/moneybags-server/constants"
	"github.com/paulwrubel/moneybags-server/controllers"
	mockservices "github.com/paulwrubel/moneybags-server/mocks/services"
	"github.com/stretchr/testify/assert"
)

func TestAuthLogin(t *testing.T) {
	tests := []struct {
		name                 string
		endpoint             string
		requestSetupFunc     func(r *http.Request) *http.Request
		mockSetupFunc        func(ma *mockservices.MockIAuth)
		requestMethod        string
		requestBody          string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "post - success",
			endpoint: "/api/v1/auth/token",
			requestSetupFunc: func(r *http.Request) *http.Request {
				return r
			},
			mockSetupFunc: func(m *mockservices.MockIAuth) {
				authCall := m.EXPECT().
					Authenticate(gomock.Eq("user_1"), gomock.Eq("pass_1")).
					Times(1).
					Return(true, nil)

				m.EXPECT().
					CreateAuthToken(gomock.Eq("user_1")).
					After(authCall).
					Times(1).
					Return("__token_1__", nil)
			},
			requestMethod: http.MethodPost,
			requestBody: `{	
				"username": "user_1",
				"password": "pass_1"
			}`,
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{
				"access_token": "__token_1__"
			}`,
		},
		{
			name:     "post - unauthorized - bad password",
			endpoint: "/api/v1/auth/token", mockSetupFunc: func(m *mockservices.MockIAuth) {
				m.EXPECT().
					Authenticate(gomock.Eq("user_1"), gomock.Eq("bad_pass_1")).
					Times(1).
					Return(false, nil)

				m.EXPECT().
					CreateAuthToken(gomock.Any()).
					Times(0)
			},
			requestMethod: http.MethodPost,
			requestBody: `{	
				"username": "user_1",
				"password": "bad_pass_1"
			}`,
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponseBody: `{
				"errors": [{
					"message": "Invalid username or password"
				}]
			}`,
		},
		{
			name:     "post - unauthorized - bad username",
			endpoint: "/api/v1/auth/login", mockSetupFunc: func(m *mockservices.MockIAuth) {
				m.EXPECT().
					Authenticate(gomock.Eq("bad_user_1"), gomock.Eq("pass_1")).
					Times(1).
					Return(false, constants.ErrUserDoesNotExist)

				m.EXPECT().
					CreateAuthToken(gomock.Any()).
					Times(0)
			},
			requestMethod: http.MethodPost,
			requestBody: `{	
				"username": "bad_user_1",
				"password": "pass_1"
			}`,
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponseBody: `{
				"errors": [{
					"message": "user does not exist"
				}]
			}`,
		},
		{
			name:     "post - failure - server error",
			endpoint: "/api/v1/auth/login", mockSetupFunc: func(m *mockservices.MockIAuth) {
				m.EXPECT().
					Authenticate(gomock.Eq("user_1"), gomock.Eq("pass_1")).
					Times(1).
					Return(false, errors.New("some internal problem occured"))

				m.EXPECT().
					CreateAuthToken(gomock.Eq("user_1")).
					Times(0)
			},
			requestMethod: http.MethodPost,
			requestBody: `{	
				"username": "user_1",
				"password": "pass_1"
			}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponseBody: `{
				"errors": [{
					"message": "some internal problem occured"
				}]
			}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := mockservices.NewMockIAuth(gomock.NewController(t))

			tt.mockSetupFunc(mockAuthService)

			a := &controllers.Auth{
				Service: mockAuthService,
			}

			rw := httptest.NewRecorder()
			r := httptest.NewRequest(tt.requestMethod, tt.endpoint, strings.NewReader(tt.requestBody))

			a.PostToken().ServeHTTP(rw, r)

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
