package controllers_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/paulwrubel/moneybags-server/constants"
	"github.com/paulwrubel/moneybags-server/controllers"
	mockservices "github.com/paulwrubel/moneybags-server/mocks/services"
	"github.com/paulwrubel/moneybags-server/models"
	"github.com/stretchr/testify/assert"
)

func pointerify(s string) *string {
	p := s
	return &p
}

func TestUserAccountsGet(t *testing.T) {
	tests := []struct {
		name                 string
		endpoint             string
		requestSetupFunc     func(r *http.Request) *http.Request
		mockSetupFunc        func(m *mockservices.MockIUserAccounts)
		requestMethod        string
		requestBody          string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "get - success",
			endpoint: "/api/v1/user-accounts",
			requestSetupFunc: func(r *http.Request) *http.Request {
				r = r.WithContext(context.WithValue(r.Context(), constants.UsernameContextKey, "user_1"))
				return r
			},
			mockSetupFunc: func(m *mockservices.MockIUserAccounts) {
				existsCall := m.EXPECT().
					ExistsByUsername(gomock.Eq("user_1")).
					Times(1).
					Return(true, nil)

				m.EXPECT().
					GetInfo(gomock.Eq("user_1")).
					After(existsCall).
					Times(1).
					Return(&models.UserAccount{
						ID:           "__uaid_1__",
						Username:     "user_1",
						PasswordHash: "__hash__",
						Email:        pointerify("user1@testing.com"),
					}, nil)
			},
			requestMethod: http.MethodGet,
			requestBody: `{
				"username": "user_1"	
			}`,
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{	
				"id": "__uaid_1__",
				"username": "user_1",
				"email": "user1@testing.com"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserAccountService := mockservices.NewMockIUserAccounts(gomock.NewController(t))

			tt.mockSetupFunc(mockUserAccountService)

			ua := &controllers.UserAccounts{
				Service: mockUserAccountService,
			}

			rw := httptest.NewRecorder()
			r := httptest.NewRequest(tt.requestMethod, tt.endpoint, strings.NewReader(tt.requestBody))
			r = tt.requestSetupFunc(r)

			ua.Get().ServeHTTP(rw, r)

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

func TestUserAccountsPost(t *testing.T) {
	tests := []struct {
		name                 string
		endpoint             string
		requestSetupFunc     func(r *http.Request) *http.Request
		mockSetupFunc        func(m *mockservices.MockIUserAccounts)
		requestMethod        string
		requestBody          string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "post - success",
			endpoint: "/api/v1/user-accounts",
			requestSetupFunc: func(r *http.Request) *http.Request {
				return r
			},
			mockSetupFunc: func(m *mockservices.MockIUserAccounts) {
				m.EXPECT().
					Create(gomock.Eq("user_1"), gomock.Eq("password"), gomock.Eq(pointerify("user1@testing.com"))).
					Times(1).
					Return(&models.UserAccount{
						ID:           "__uaid_1__",
						Username:     "user_1",
						PasswordHash: "__hash__",
						Email:        pointerify("user1@testing.com"),
					}, nil)
			},
			requestMethod: http.MethodPost,
			requestBody: `{
				"username": "user_1",
				"password": "password",
				"email": "user1@testing.com"	
			}`,
			expectedStatusCode: http.StatusCreated,
			expectedResponseBody: `{	
				"id": "__uaid_1__",
				"username": "user_1",
				"email": "user1@testing.com"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserAccountService := mockservices.NewMockIUserAccounts(gomock.NewController(t))

			tt.mockSetupFunc(mockUserAccountService)

			ua := &controllers.UserAccounts{
				Service: mockUserAccountService,
			}

			rw := httptest.NewRecorder()
			r := httptest.NewRequest(tt.requestMethod, tt.endpoint, strings.NewReader(tt.requestBody))
			r = tt.requestSetupFunc(r)

			ua.Post().ServeHTTP(rw, r)

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
