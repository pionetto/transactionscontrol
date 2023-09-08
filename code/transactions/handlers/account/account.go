package account

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cajueiro/code/transactions/models"
	"cajueiro/pkg/app"

	"github.com/gorilla/mux"
)

// ListAccounts lista as contas no banco de dados
func ListAccounts(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		//Pegando as contas no banco de dados
		var a []models.Account
		if err := app.DB.Client.Find(&a); err.Error != nil {
			// Se encontrar erro, retorna StatusInternalServerError (erro 500)
			http.Error(w, "Erro ao listar as contas", http.StatusInternalServerError)
			return
		}
		// Retorno do JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(a)

	}
}

// PostAccount cria uma conta no banco de dados
func PostAccount(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// Pegando account no request
		a := &models.Account{}
		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			// Se encontrar erro, retorna StatusBadRequest (erro 400)
			http.Error(w, "Formato(JSON) inválido", http.StatusBadRequest)
			return
		}

		// Validação do json de Account
		if err := app.Vld.Struct(a); err != nil {
			// tradução dos erros do JSON com formato inválido
			errs := app.TranslateErrors(err)
			// Se o body do request for inválido retorna StatusBadRequest (erro 400)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, errs)
			return
		}

		// armazenando struct account no DB
		account, err := a.CreateAccount(app)
		if err != nil {
			// caso tenha erro ao armazenar no banco retorna 500
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(account)

	}
}

// BalanceAccount retorna o saldo da conta no banco de dados
func BalanceAccount(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// Pegando id na url
		id := mux.Vars(r)["id"]

		// Pegando account no banco de dados
		a := &models.Account{}
		if err := app.DB.Client.First(&a, &id); err.Error != nil {
			// caso tenha erro ao procurar no banco retorna 404
			http.Error(w, "Conta não encontrada", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]float64{"amount_food": a.Amount_food})
		json.NewEncoder(w).Encode(map[string]float64{"amount_meal": a.Amount_meal})
		json.NewEncoder(w).Encode(map[string]float64{"amount_cash": a.Amount_cash})
	}
}
