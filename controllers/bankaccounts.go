package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/paulwrubel/moneybags-server/services"
	log "github.com/sirupsen/logrus"
)

type BankAccounts struct {
	SBankAccounts services.IBankAccounts
	SBudgets      services.IBudgets
	SUserAccounts services.IUserAccounts
}

type getAllBankAccountsResponse struct {
	BankAccounts []getAllBankAccountsResponseBankAccount `json:"bank_accounts"`
}

type getAllBankAccountsResponseBankAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (ba *BankAccounts) GetAll() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userAccount, ok := validateUser(rw, r, ba.SUserAccounts)
		if !ok {
			return
		}

		budgetID := mux.Vars(r)["budgetID"]

		exists, err := ba.SBudgets.ExistsByID(budgetID)
		if err != nil {
			log.WithError(err).Error("Error checking if budget exists")
			writeResponse(rw, http.StatusInternalServerError, errorsResponseFromErrors(err))
			return
		}
		if !exists {
			writeResponse(rw, http.StatusNotFound, errorsResponseFromMessages("Budget does not exist"))
			return
		}

		belongsToRequestor, err := ba.SBudgets.BelongsTo(userAccount.ID, budgetID)
		if err != nil {
			log.WithError(err).Error("Error checking if budget belongs to user")
			writeResponse(rw, http.StatusInternalServerError, errorsResponseFromErrors(err))
			return
		}
		if !belongsToRequestor {
			writeResponse(rw, http.StatusForbidden, errorsResponseFromMessages("Budget does not belong to user"))
			return
		}

		bankAccounts, err := ba.SBankAccounts.GetAll(budgetID)
		if err != nil {
			log.WithError(err).Error("Error getting all accounts")
			writeResponse(rw, http.StatusInternalServerError, errorsResponseFromErrors(err))
			return
		}

		response := getAllBankAccountsResponse{
			BankAccounts: []getAllBankAccountsResponseBankAccount{},
		}
		for _, account := range bankAccounts {
			response.BankAccounts = append(response.BankAccounts, getAllBankAccountsResponseBankAccount{
				ID:   account.ID,
				Name: account.Name,
			})
		}

		writeResponse(rw, http.StatusOK, response)
	}
}
