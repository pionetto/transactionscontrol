package merchant

import (
	"encoding/json"
	"net/http"

	"cajueiro/code/transactions/models"
	"cajueiro/pkg/app"

	"github.com/gorilla/mux"
)

// ListMerchants - handler para listar os estabelecimentos no DB
func ListMerchants(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()

		// capturando nome do estabelecimento na url
		merchant := mux.Vars(request)["merchant"]

		// capturando estabelecimentos no DB
		var t []models.Transaction
		if err := app.DB.Client.Where("merchant LIKE ?", "%"+merchant+"%").Find(&t); err.Error != nil {
			// caso tenha erro ao procurar no banco, retorna 500
			http.Error(w, "Erro na listagem dos estabelecimentos", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(t)

	}
}
