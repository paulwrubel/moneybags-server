package routing

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/paulwrubel/moneybags-server/injection"
	log "github.com/sirupsen/logrus"
)

func RunServer(injector injection.IInjector) {
	router := getRouter(injector)

	go func() {
		err := http.ListenAndServe(":8080", router)
		if err != nil {
			log.WithError(err).Fatalln("error in RunServer()")
		}
	}()
}

func getRouter(injector injection.IInjector) *mux.Router {
	router := mux.NewRouter()

	router.Use(logrusMiddleware())

	// healthcheck routes
	healthController := injector.InjectHealthController()
	router.HandleFunc("/ping", healthController.Get()).Methods("GET")
	router.HandleFunc("/health", healthController.Get()).Methods("GET")

	apiSubrouter := router.PathPrefix("/api/v1").Subrouter()

	budgetController := injector.InjectBudgetController()
	budgetSubrouter := apiSubrouter.PathPrefix("/budgets").Subrouter()
	budgetSubrouter.HandleFunc("", budgetController.GetAll()).Methods("GET")
	budgetSubrouter.HandleFunc("/{budget_id}", budgetController.Get()).Methods("GET")
	budgetSubrouter.HandleFunc("", budgetController.Post()).Methods("POST")

	return router
}
