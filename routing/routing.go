package routing

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/paulwrubel/moneybags-server/injection"
	"github.com/paulwrubel/moneybags-server/middleware"
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

	authService := injector.InjectAuthService()

	router.Use(middleware.Logrus())
	router.Use(handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost}),
		handlers.AllowedOrigins([]string{"*"}),
	))
	auth := middleware.SessionValidation(authService)

	// healthcheck routes
	healthController := injector.InjectHealthController()
	router.HandleFunc("/ping", healthController.Get()).Methods(http.MethodGet)
	router.HandleFunc("/health", healthController.Get()).Methods(http.MethodGet)

	apiSubrouter := router.PathPrefix("/api/v1").Subrouter()

	// user account routes
	userAccountsController := injector.InjectUserAccountsController()
	userAccountsSubrouter := apiSubrouter.PathPrefix("/user-accounts").Subrouter()
	userAccountsSubrouter.HandleFunc("", userAccountsController.Post()).Methods(http.MethodPost)
	userAccountsSubrouter.Handle("", auth(userAccountsController.Get())).Methods(http.MethodGet)

	// auth routes
	authController := injector.InjectAuthController(authService)
	authSubrouter := apiSubrouter.PathPrefix("/auth").Subrouter()
	authSubrouter.HandleFunc("/token", authController.PostToken()).Methods(http.MethodPost)

	// budget routes
	budgetsController := injector.InjectBudgetsController()
	budgetsSubrouter := apiSubrouter.PathPrefix("/budgets").Subrouter()
	budgetsSubrouter.Use(auth)
	budgetsSubrouter.HandleFunc("", budgetsController.GetAll()).Methods(http.MethodGet)
	budgetsSubrouter.HandleFunc("/{budgetID}", budgetsController.Get()).Methods(http.MethodGet)
	budgetsSubrouter.HandleFunc("", budgetsController.Post()).Methods(http.MethodPost)

	// bank account routes
	bankAccountsController := injector.InjectBankAccountsController()
	bankAccountsSubrouter := apiSubrouter.PathPrefix("/budgets/{budgetID}/bank-accounts").Subrouter()
	bankAccountsSubrouter.Use(auth)
	bankAccountsSubrouter.HandleFunc("", bankAccountsController.GetAll()).Methods(http.MethodGet)

	return router
}
