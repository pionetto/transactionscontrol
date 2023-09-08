package transaction

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cajueiro/code/transactions/models"
	"cajueiro/pkg/app"

	"github.com/dgrijalva/jwt-go"
)

// ListTransactions lista as transaferencias da conta no banco de dados
func ListTransactions(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// criando a chave JWT usada para verificar a assinatura
		var jwtKey = []byte(app.Cfg.GetTokenKey())

		// capturando o token JWT no cabeçalho do request
		if r.Header["Token"] == nil {
			// caso o token seja nulo retorna 401
			http.Error(w, "Token nulo", http.StatusUnauthorized)
			return
		}

		// capturar a string do token JWT
		tknStr := r.Header.Get("Token")

		// inicializar um struct claims
		claims := &models.Claims{}

		// Parse da string JWT e armazena o resultado no struct claims
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Assinatura inválida", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			http.Error(w, "Token Expirou", http.StatusUnauthorized)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// capturando transactions no DB
		a := &models.Account{}
		if err := app.DB.Client.Preload("Transaction").First(&a, "cpf = ?", claims.CPF); err.Error != nil {
			// caso tenha erro ao procurar no banco retorna 500
			http.Error(w, "Erro na listagem das transferências", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(a.Transaction)

	}
}

// PostTransactions handler para criar transactions no DB
func PostTransactions(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// criando a chave JWT usada para verificar a assinatura
		var jwtKey = []byte(app.Cfg.GetTokenKey())

		// capturando o token JWT no cabeçalho do request
		if r.Header["Token"] == nil {
			// caso o token seja nulo retorna 401
			http.Error(w, "Token nulo", http.StatusUnauthorized)
			return
		}

		// capturar a string do token JWT
		tknStr := r.Header.Get("Token")

		// inicializar um struct claims
		claims := &models.Claims{}

		// Parse da string JWT e armazena o resultado no struct claims
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Assinatura inválida", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			http.Error(w, "Token Expirou", http.StatusUnauthorized)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// capturando account no DB
		a := &models.Account{}
		if err := app.DB.Client.First(&a, "cpf = ?", claims.CPF); err.Error != nil {
			// caso tenha erro ao procurar no banco retorna 500
			http.Error(w, "Erro na criação da transferência", http.StatusInternalServerError)
			return
		}

		// capturando transactions no request
		t := &models.Transaction{}
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			// caso tenha erro no decode do request retorna 400
			http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
			return
		}

		// adicionando ID da conta de origem
		// t.AccountOriginID = a.ID
		t.Account_id = a.ID

		// validando json do struct transaction
		if err := app.Vld.Struct(t); err != nil {
			// traduzindo os erros do JSON inválido
			errs := app.TranslateErrors(err)
			// caso o corpo do request seja inválido retorna 400
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, errs)
			return
		}

		// armazenando struct transaction no DB
		transaction, err := t.CreateTransaction(app)
		if err != nil {
			// caso tenha erro ao armazenar no banco retorna 500
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(transaction)

	}
}
