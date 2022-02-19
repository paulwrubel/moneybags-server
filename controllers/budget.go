package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/paulwrubel/moneybags-server/services"
	log "github.com/sirupsen/logrus"
)

type Budget struct {
	Service services.IBudget
}

type getAllResponse struct {
	Budgets []getAllResponseBudget `json:"budgets"`
}

type getAllResponseBudget struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (b *Budget) GetAll() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		budgets, err := b.Service.GetAll()
		if err != nil {
			log.WithError(err).Error("Error getting all budgets")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := getAllResponse{
			Budgets: []getAllResponseBudget{},
		}
		for _, budget := range budgets {
			response.Budgets = append(response.Budgets, getAllResponseBudget{
				ID:   budget.ID,
				Name: budget.Name,
			})
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			log.WithError(err).Error("Error marshalling response")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(responseBytes)
	}
}

type getResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (b *Budget) Get() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["budget_id"]

		exists, err := b.Service.Exists(id)
		if err != nil {
			log.WithError(err).Error("Error checking if budget exists")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !exists {
			log.Error("Budget does not exist")
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		budget, err := b.Service.Get(id)
		if err != nil {
			log.WithError(err).Error("Error getting budget")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := getResponse{
			ID:   budget.ID,
			Name: budget.Name,
		}
		responseBytes, err := json.Marshal(response)
		if err != nil {
			log.WithError(err).Error("Error marshalling response")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(responseBytes)
	}
}

type postRequest struct {
	Name string `json:"name"`
}

type postResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (b *Budget) Post() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).Error("Error reading request body")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		var requestBody postRequest
		err = json.Unmarshal(bodyBytes, &requestBody)
		if err != nil {
			log.WithError(err).Error("Error unmarshalling request body")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		createdBudget, err := b.Service.Create(requestBody.Name)
		if err != nil {
			log.WithError(err).Error("Error creating budget")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := postResponse{
			ID:   createdBudget.ID,
			Name: createdBudget.Name,
		}
		responseBytes, err := json.Marshal(response)
		if err != nil {
			log.WithError(err).Error("Error marshalling response")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(responseBytes)
	}
}
