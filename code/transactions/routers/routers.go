package routers

import (
	"cajueiro/code/transactions/handlers/account"
	"cajueiro/code/transactions/handlers/login"
	"cajueiro/code/transactions/handlers/merchant"
	"cajueiro/code/transactions/handlers/transaction"
	"cajueiro/pkg/app"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// GetRouter retorna o roteador mux da API
func GetRouter(app *app.App) *mux.Router {

	// middleware compartilhado em todas as rotas da API
	common := negroni.New(
		negroni.NewLogger(),
	)

	// criando roteador base
	router := mux.NewRouter()

	// rota de login
	loginRoutes := mux.NewRouter()
	router.Path("/login").Handler(common.With(
		negroni.Wrap(loginRoutes),
	))
	logins := loginRoutes.Path("/login").Subrouter()
	logins.Methods("POST").HandlerFunc(login.HandlerLogin(app))

	// rota de accounts
	accountsRoutes := mux.NewRouter()
	router.Path("/accounts").Handler(common.With(
		negroni.Wrap(accountsRoutes),
	))
	accounts := accountsRoutes.Path("/accounts").Subrouter()
	accounts.Methods("GET").HandlerFunc(account.ListAccounts(app))
	accounts.Methods("POST").HandlerFunc(account.PostAccount(app))

	// rota de balance
	balanceRoutes := mux.NewRouter()
	router.Path("/accounts/{id}/balance").Handler(common.With(
		negroni.Wrap(balanceRoutes),
	))
	balances := balanceRoutes.Path("/accounts/{id}/balance").Subrouter()
	balances.Methods("GET").HandlerFunc(account.BalanceAccount(app))

	// rota de transações (transactions)
	transactionsRoutes := mux.NewRouter()
	router.Path("/transactions").Handler(common.With(
		negroni.Wrap(transactionsRoutes),
	))
	transactions := transactionsRoutes.Path("/transactions").Subrouter()
	transactions.Methods("GET").HandlerFunc(transaction.ListTransactions(app))
	transactions.Methods("POST").HandlerFunc(transaction.PostTransactions(app))

	// rota de estabelecimentos
	merchantsRoutes := mux.NewRouter()
	router.Path("/merchants/{merchant}").Handler(common.With(
		negroni.Wrap(merchantsRoutes),
	))
	merchants := merchantsRoutes.Path("/merchants/{merchant}").Subrouter()
	merchants.Methods("GET").HandlerFunc(merchant.ListMerchants(app))

	return router
}
