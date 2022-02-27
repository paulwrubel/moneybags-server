package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/paulwrubel/moneybags-server/services"
	log "github.com/sirupsen/logrus"
)

type Budgets struct {
	SBudgets      services.IBudgets
	SUserAccounts services.IUserAccounts
}

type getAllBudgetsResponse struct {
	Budgets []getAllBudgetsResponseBudget `json:"budgets"`
}

type getAllBudgetsResponseBudget struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (b *Budgets) GetAll() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userAccount, ok := validateUser(rw, r, b.SUserAccounts)
		if !ok {
			return
		}

		budgets, err := b.SBudgets.GetAllByUserAccountID(userAccount.ID)
		if err != nil {
			log.WithError(err).Error("Error getting all budgets")
			writeResponse(rw, http.StatusInternalServerError, errorsResponseFromErrors(err))
			return
		}

		response := getAllBudgetsResponse{
			Budgets: []getAllBudgetsResponseBudget{},
		}
		for _, budget := range budgets {
			response.Budgets = append(response.Budgets, getAllBudgetsResponseBudget{
				ID:   budget.ID,
				Name: budget.Name,
			})
		}
		writeResponse(rw, http.StatusOK, response)
	}
}

type getBudgetResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (b *Budgets) Get() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userAccount, ok := validateUser(rw, r, b.SUserAccounts)
		if !ok {
			return
		}

		budgetID := mux.Vars(r)["budgetID"]

		exists, err := b.SBudgets.ExistsByID(budgetID)
		if err != nil {
			log.WithError(err).Error("Error checking if budget exists")
			writeResponse(rw, http.StatusInternalServerError, errorsResponseFromErrors(err))
			return
		}
		if !exists {
			writeResponse(rw, http.StatusNotFound, errorsResponseFromMessages("Budget does not exist"))
			return
		}

		belongsToRequestor, err := b.SBudgets.BelongsTo(userAccount.ID, budgetID)
		if err != nil {
			log.WithError(err).Error("Error checking if budget belongs to user")
			writeResponse(rw, http.StatusInternalServerError, errorsResponseFromErrors(err))
			return
		}
		if !belongsToRequestor {
			writeResponse(rw, http.StatusForbidden, errorsResponseFromMessages("Budget does not belong to user"))
			return
		}

		budget, err := b.SBudgets.GetByID(budgetID)
		if err != nil {
			log.WithError(err).Error("Error getting budget")
			writeResponse(rw, http.StatusInternalServerError, errorsResponseFromErrors(err))
			return
		}

		writeResponse(rw, http.StatusOK, getBudgetResponse{
			ID:   budget.ID,
			Name: budget.Name,
		})
	}
}

type postBudgetRequest struct {
	Name string `json:"name"`
}

type postBudgetResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (b *Budgets) Post() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userAccount, ok := validateUser(rw, r, b.SUserAccounts)
		if !ok {
			return
		}

		var requestBody postBudgetRequest
		err := unmarshalRequestBody(r.Body, &requestBody)
		if err != nil {
			writeResponse(rw, http.StatusBadRequest, errorsResponseFromErrors(err))
			return
		}

		exists, err := b.SBudgets.ExistsByUserIDAndName(userAccount.ID, requestBody.Name)
		if err != nil {
			log.WithError(err).Error("Error checking if budget exists")
			writeResponse(rw, http.StatusInternalServerError, errorsResponseFromErrors(err))
			return
		}
		if exists {
			writeResponse(rw, http.StatusConflict, errorsResponseFromMessages("Budget already exists"))
			return
		}

		createdBudget, err := b.SBudgets.Create(userAccount.ID, requestBody.Name)
		if err != nil {
			log.WithError(err).Error("Error creating budget")
			writeResponse(rw, http.StatusInternalServerError, errorsResponseFromErrors(err))
			return
		}

		writeResponse(rw, http.StatusCreated, postBudgetResponse{
			ID:   createdBudget.ID,
			Name: createdBudget.Name,
		})
	}
}
