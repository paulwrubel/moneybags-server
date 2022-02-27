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
	"github.com/gorilla/mux"
	"github.com/paulwrubel/moneybags-server/constants"
	"github.com/paulwrubel/moneybags-server/controllers"
	mockservices "github.com/paulwrubel/moneybags-server/mocks/services"
	"github.com/paulwrubel/moneybags-server/models"
	"github.com/stretchr/testify/assert"
)

func TestBudgetsGetAll(t *testing.T) {
	tests := []struct {
		name                 string
		endpoint             string
		requestSetupFunc     func(r *http.Request) *http.Request
		mockSetupFunc        func(mb *mockservices.MockIBudgets, mua *mockservices.MockIUserAccounts)
		requestMethod        string
		requestBody          string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "get all - success",
			endpoint: "/api/v1/budgets",
			requestSetupFunc: func(r *http.Request) *http.Request {
				r = r.WithContext(context.WithValue(r.Context(), constants.UsernameContextKey, "user_1"))
				return r
			},
			mockSetupFunc: func(mb *mockservices.MockIBudgets, mua *mockservices.MockIUserAccounts) {
				userExistsCall := mua.EXPECT().
					ExistsByUsername(gomock.Eq("user_1")).
					Times(1).
					Return(true, nil)

				getUserCall := mua.EXPECT().
					GetInfo(gomock.Eq("user_1")).
					After(userExistsCall).
					Times(1).
					Return(&models.UserAccount{
						ID:           "__uaid_1__",
						Username:     "user_1",
						PasswordHash: "__hash__",
						Email:        nil,
					}, nil)

				mb.EXPECT().
					GetAllByUserAccountID(gomock.Eq("__uaid_1__")).
					After(getUserCall).
					Times(1).
					Return([]*models.Budget{
						{
							ID:            "__bid_1__",
							UserAccountID: "__uaid_1__",
							Name:          "budget_1",
						},
						{
							ID:            "__bid_2__",
							UserAccountID: "__uaid_1__",
							Name:          "budget_2",
						},
					}, nil)
			},
			requestMethod:      http.MethodGet,
			requestBody:        ``,
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{
				"budgets": [
					{
						"id": "__bid_1__",
						"name": "budget_1"
					},
					{
						"id": "__bid_2__",
						"name": "budget_2"
					}
				]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBudgetsService := mockservices.NewMockIBudgets(gomock.NewController(t))
			mockUserAccountsService := mockservices.NewMockIUserAccounts(gomock.NewController(t))

			tt.mockSetupFunc(mockBudgetsService, mockUserAccountsService)

			b := &controllers.Budgets{
				SBudgets:      mockBudgetsService,
				SUserAccounts: mockUserAccountsService,
			}

			rw := httptest.NewRecorder()
			r := httptest.NewRequest(tt.requestMethod, tt.endpoint, strings.NewReader(tt.requestBody))
			r = tt.requestSetupFunc(r)

			b.GetAll().ServeHTTP(rw, r)

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

func TestBudgetsGet(t *testing.T) {
	tests := []struct {
		name                 string
		endpoint             string
		requestSetupFunc     func(r *http.Request) *http.Request
		urlVars              map[string]string
		mockSetupFunc        func(mb *mockservices.MockIBudgets, mua *mockservices.MockIUserAccounts)
		requestMethod        string
		requestBody          string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "get - success",
			endpoint: "/api/v1/budgets/__bid_1__",
			requestSetupFunc: func(r *http.Request) *http.Request {
				r = mux.SetURLVars(r, map[string]string{
					"budgetID": "__bid_1__",
				})
				r = r.WithContext(context.WithValue(r.Context(), constants.UsernameContextKey, "user_1"))
				return r
			},
			mockSetupFunc: func(mb *mockservices.MockIBudgets, mua *mockservices.MockIUserAccounts) {
				userExistsCall := mua.EXPECT().
					ExistsByUsername(gomock.Eq("user_1")).
					Times(1).
					Return(true, nil)

				getUserCall := mua.EXPECT().
					GetInfo(gomock.Eq("user_1")).
					After(userExistsCall).
					Times(1).
					Return(&models.UserAccount{
						ID:           "__uaid_1__",
						Username:     "user_1",
						PasswordHash: "__hash__",
						Email:        nil,
					}, nil)

				existsCall := mb.EXPECT().
					ExistsByID(gomock.Eq("__bid_1__")).
					After(getUserCall).
					Times(1).
					Return(true, nil)

				belongsToCall := mb.EXPECT().
					BelongsTo(gomock.Eq("__uaid_1__"), gomock.Eq("__bid_1__")).
					After(existsCall).
					Times(1).
					Return(true, nil)

				mb.EXPECT().
					GetByID(gomock.Eq("__bid_1__")).
					After(belongsToCall).
					Times(1).
					Return(&models.Budget{
						ID:   "__bid_1__",
						Name: "budget_1",
					}, nil)
			},
			requestMethod:      http.MethodGet,
			requestBody:        ``,
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{
				"id": "__bid_1__",
				"name": "budget_1"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBudgetsService := mockservices.NewMockIBudgets(gomock.NewController(t))
			mockUserAccountsService := mockservices.NewMockIUserAccounts(gomock.NewController(t))

			tt.mockSetupFunc(mockBudgetsService, mockUserAccountsService)

			b := &controllers.Budgets{
				SBudgets:      mockBudgetsService,
				SUserAccounts: mockUserAccountsService,
			}

			rw := httptest.NewRecorder()
			r := httptest.NewRequest(tt.requestMethod, tt.endpoint, strings.NewReader(tt.requestBody))
			r = tt.requestSetupFunc(r)

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

func TestBudgetsPost(t *testing.T) {
	tests := []struct {
		name                 string
		endpoint             string
		requestSetupFunc     func(r *http.Request) *http.Request
		mockSetupFunc        func(mb *mockservices.MockIBudgets, mua *mockservices.MockIUserAccounts)
		requestMethod        string
		requestBody          string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "post - success",
			endpoint: "/api/v1/budgets",
			requestSetupFunc: func(r *http.Request) *http.Request {
				r = r.WithContext(context.WithValue(r.Context(), constants.UsernameContextKey, "user_1"))
				return r
			},
			mockSetupFunc: func(mb *mockservices.MockIBudgets, mua *mockservices.MockIUserAccounts) {
				userExistsCall := mua.EXPECT().
					ExistsByUsername(gomock.Eq("user_1")).
					Times(1).
					Return(true, nil)

				getUserCall := mua.EXPECT().
					GetInfo(gomock.Eq("user_1")).
					After(userExistsCall).
					Times(1).
					Return(&models.UserAccount{
						ID:           "__uaid_1__",
						Username:     "user_1",
						PasswordHash: "__hash__",
						Email:        nil,
					}, nil)

				existsCall := mb.EXPECT().
					ExistsByUserIDAndName(gomock.Eq("__uaid_1__"), gomock.Eq("budget_1")).
					After(getUserCall).
					Times(1).
					Return(false, nil)

				mb.EXPECT().
					Create(gomock.Eq("__uaid_1__"), gomock.Eq("budget_1")).
					After(existsCall).
					Times(1).
					Return(&models.Budget{
						ID:            "__bid_1__",
						UserAccountID: "__uaid_1__",
						Name:          "budget_1",
					}, nil)
			},
			requestMethod: http.MethodPost,
			requestBody: `{
				"name": "budget_1"
			}`,
			expectedStatusCode: http.StatusCreated,
			expectedResponseBody: `{
				"id": "__bid_1__",
				"name": "budget_1"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBudgetsService := mockservices.NewMockIBudgets(gomock.NewController(t))
			mockUserAccountsService := mockservices.NewMockIUserAccounts(gomock.NewController(t))

			tt.mockSetupFunc(mockBudgetsService, mockUserAccountsService)

			b := &controllers.Budgets{
				SBudgets:      mockBudgetsService,
				SUserAccounts: mockUserAccountsService,
			}

			rw := httptest.NewRecorder()
			r := httptest.NewRequest(tt.requestMethod, tt.endpoint, strings.NewReader(tt.requestBody))
			r = tt.requestSetupFunc(r)

			b.Post().ServeHTTP(rw, r)

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
