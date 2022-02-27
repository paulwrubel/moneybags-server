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

func TestBankAccountsGetAll(t *testing.T) {
	tests := []struct {
		name                 string
		endpoint             string
		requestSetupFunc     func(r *http.Request) *http.Request
		mockSetupFunc        func(mba *mockservices.MockIBankAccounts, mb *mockservices.MockIBudgets, mua *mockservices.MockIUserAccounts)
		requestMethod        string
		requestBody          string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "get all - success",
			endpoint: "/api/v1/budgets/__bid_1__/bank-accounts",
			requestSetupFunc: func(r *http.Request) *http.Request {
				r = r.WithContext(context.WithValue(r.Context(), constants.UsernameContextKey, "user_1"))
				r = mux.SetURLVars(r, map[string]string{"budgetID": "__bid_1__"})
				return r
			},
			mockSetupFunc: func(mba *mockservices.MockIBankAccounts, mb *mockservices.MockIBudgets, mua *mockservices.MockIUserAccounts) {
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

				mba.EXPECT().
					GetAll(gomock.Eq("__bid_1__")).
					After(belongsToCall).
					Times(1).
					Return([]*models.BankAccount{
						{
							ID:       "__baid_1__",
							BudgetID: "__bid_1__",
							Name:     "bank_account_1",
						}, {
							ID:       "__baid_2__",
							BudgetID: "__bid_1__",
							Name:     "bank_account_2",
						},
					}, nil)

			},
			requestMethod:      http.MethodGet,
			requestBody:        ``,
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{
				"bank_accounts": [
					{	
						"id": "__baid_1__",
						"name": "bank_account_1"
					},
					{
						"id": "__baid_2__",
						"name": "bank_account_2"
					}
				]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBankAccountsService := mockservices.NewMockIBankAccounts(gomock.NewController(t))
			mockBudgetsService := mockservices.NewMockIBudgets(gomock.NewController(t))
			mockUserAccountsService := mockservices.NewMockIUserAccounts(gomock.NewController(t))

			tt.mockSetupFunc(mockBankAccountsService, mockBudgetsService, mockUserAccountsService)

			ba := &controllers.BankAccounts{
				SBankAccounts: mockBankAccountsService,
				SBudgets:      mockBudgetsService,
				SUserAccounts: mockUserAccountsService,
			}

			rw := httptest.NewRecorder()
			r := httptest.NewRequest(tt.requestMethod, tt.endpoint, strings.NewReader(tt.requestBody))
			r = tt.requestSetupFunc(r)

			ba.GetAll().ServeHTTP(rw, r)

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
