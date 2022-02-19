package controllers

import (
	"net/http"
)

type Health struct{}

func (h *Health) Get() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}
}
