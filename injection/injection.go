package injection

import (
	"github.com/paulwrubel/moneybags-server/config"
	"github.com/paulwrubel/moneybags-server/controllers"
	"github.com/paulwrubel/moneybags-server/repositories"
	"github.com/paulwrubel/moneybags-server/services"
)

type IInjector interface {
	InjectAuthService() *services.Auth

	InjectHealthController() *controllers.Health
	InjectAuthController(service services.IAuth) *controllers.Auth
	InjectUserAccountsController() *controllers.UserAccounts
	InjectBudgetsController() *controllers.Budgets
	InjectBankAccountsController() *controllers.BankAccounts
}

type Injector struct {
	AppInfo *config.AppInfo
}

func (i *Injector) InjectAuthService() *services.Auth {
	return &services.Auth{
		JWTIssuer:     i.AppInfo.AuthInfo.JWTIssuer,
		SigningMethod: i.AppInfo.AuthInfo.SigningMethod,
		PrivateKey:    i.AppInfo.AuthInfo.PrivateKey,
		UserAccounts: &repositories.UserAccounts{
			DB: i.AppInfo.DB,
		},
	}
}

func (i *Injector) InjectHealthController() *controllers.Health {
	return &controllers.Health{}
}

func (i *Injector) InjectAuthController(service services.IAuth) *controllers.Auth {
	return &controllers.Auth{
		Service: service,
	}
}

func (i *Injector) InjectUserAccountsController() *controllers.UserAccounts {
	return &controllers.UserAccounts{
		Service: &services.UserAccounts{
			Repository: &repositories.UserAccounts{
				DB: i.AppInfo.DB,
			},
		},
	}
}

func (i *Injector) InjectBudgetsController() *controllers.Budgets {
	return &controllers.Budgets{
		SBudgets: &services.Budgets{
			RBudgets: &repositories.Budgets{
				DB: i.AppInfo.DB,
			},
		},
		SUserAccounts: &services.UserAccounts{
			Repository: &repositories.UserAccounts{
				DB: i.AppInfo.DB,
			},
		},
	}
}

func (i *Injector) InjectBankAccountsController() *controllers.BankAccounts {
	return &controllers.BankAccounts{
		SBankAccounts: &services.BankAccounts{
			Repository: &repositories.BankAccounts{
				DB: i.AppInfo.DB,
			},
		},
		SBudgets: &services.Budgets{
			RBudgets: &repositories.Budgets{
				DB: i.AppInfo.DB,
			},
		},
		SUserAccounts: &services.UserAccounts{
			Repository: &repositories.UserAccounts{
				DB: i.AppInfo.DB,
			},
		},
	}
}
