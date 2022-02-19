package injection

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/paulwrubel/moneybags-server/controllers"
	"github.com/paulwrubel/moneybags-server/repositories"
	"github.com/paulwrubel/moneybags-server/services"
)

type IInjector interface {
	InjectHealthController() *controllers.Health
	InjectBudgetController() *controllers.Budget
}

type PostgresInjector struct {
	Database *pgxpool.Pool
}

func (i *PostgresInjector) InjectHealthController() *controllers.Health {
	return &controllers.Health{}
}

func (i *PostgresInjector) InjectBudgetController() *controllers.Budget {
	r := &repositories.Budget{
		DB: i.Database,
	}
	s := &services.Budget{
		Repository: r,
	}
	c := &controllers.Budget{
		Service: s,
	}
	return c
}
