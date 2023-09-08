package models

import (
	"errors"
	"time"

	"cajueiro/pkg/app"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BeforeCreate hook do gorm para gerar uuid no create
func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return
}

// Transaction modelo para transação do usuário
type Transaction struct {
	gorm.Model         `json:"-"`
	ID                 uuid.UUID      `json:"id" gorm:"type:uuid"` // IDENTIFICADOR UNICO DA TRANSAÇÃO
	Accounttocredit_id int            `json:"accounttocredit_id"`
	Account_id         int            `json:"account_id"` // IDENTIFICADOR DA CONTA DA QUAL FOI DEBITADO
	AccountID          int            // ID DE REFERÊNCIA NA TABELA DE CONTA
	Amount             float64        `json:"amount" gorm:"type:numeric"`
	Merchant           string         `json:"merchant"`
	Mcc                string         `json:"mcc"`
	Message            string         `json:"message"`
	Code               string         `json:"code"`
	CreatedAt          time.Time      `json:"created"`
	UpdatedAt          time.Time      `json:"updated"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted"`
}

// CreateTransaction realiza uma transação entre contas
func (t *Transaction) CreateTransaction(app *app.App) (*Transaction, error) {

	// inicia o modo de transaction
	tx := app.DB.Client.Begin()

	// verifica se a conta de destino existe
	if err := t.checkDestinationAccount(app); err != nil {

		// caso não encontre faz rollback
		tx.Rollback()
		return nil, err
	}

	// verifica se a conta de origem tem saldo suficiente
	if err := t.checkOriginBalance(app); err != nil {

		// caso não tenha saldo faz rollback
		tx.Rollback()

		return nil, err
	}

	// cria o struct transaction no DB
	transaction := &Transaction{
		ID:                 t.ID,
		Accounttocredit_id: t.Accounttocredit_id,
		Account_id:         t.Account_id,
		AccountID:          t.AccountID,
		Amount:             t.Amount,
		Merchant:           t.Merchant,
		Mcc:                t.Mcc,
		Message:            t.Message,
		Code:               t.Code,
		CreatedAt:          t.CreatedAt,
		UpdatedAt:          t.UpdatedAt,
		DeletedAt:          t.DeletedAt,
	}
	if err := tx.Create(&transaction); err.Error != nil {

		// caso ocorra erro faz rollback
		tx.Rollback()
		return nil, errors.New("Erro na criação da transação")
	}

	switch t.Code {
	case "500":
		// transação sem erros é comitada
		tx.Commit()
	case "200":

		// atualiza o saldo da conta de origem
		if err := t.balanceOriginAccount(app); err != nil {

			// caso ocorra erro faz rollback
			tx.Rollback()
			return nil, err
		}

		// atualiza o saldo da conta de destino
		if err := t.balanceDestinationAccount(app); err != nil {

			// caso ocorra erro faz rollback
			tx.Rollback()
			return nil, err
		}

		// transação sem erros é comitada
		tx.Commit()

		// caso sucesso retorna erro nulo
		return transaction, nil

	}
	return transaction, nil

}

// checkDestinationAccount verifica se a conta de destino existe
func (t *Transaction) checkDestinationAccount(app *app.App) error {

	if t.Accounttocredit_id == t.Account_id {
		return errors.New("Contas de transação devem ser diferentes")
	}

	// captura a conta de destino no banco
	a := &Account{}

	if result := app.DB.Client.First(&a, &t.Accounttocredit_id); result.Error != nil {
		return errors.New("Conta de destino não encontrada")
	}

	// retorna exista conta de destino retorna erro nulo
	return nil

}

// checkOriginBalance verifica se a conta de origem tem saldo suficiente
func (t *Transaction) checkOriginBalance(app *app.App) error {

	// captura a conta de origem no banco
	a := &Account{}
	if result := app.DB.Client.First(&a, &t.Account_id); result.Error != nil {
		return errors.New("Conta de origem não encontrada")
	}

	/*
		Se o `mcc` for `"5411" ou "5412"`, deve-se utilizar o saldo de `FOOD` - Amount_food
		Se o `mcc` for `"5811" ou "5812"`, deve-se utilizar o saldo de `MEAL`.- Amount_meal
		Para quaisquer outros valores do `mcc`, deve-se utilizar o saldo de `CASH` - Amount_cash
	*/

	// caso não tenha saldo suficiente retorna erro adequado
	switch t.Mcc {

	case "5411":

		if (a.Amount_food - t.Amount) < 0 {
			t.Message = "Transação não autorizada - Saldo na conta insuficiente - food"
			t.Code = "500"
		} else {
			t.Code = "200"
			t.Message = "Transação autorizada"
		}

	case "5412":
		if (a.Amount_food - t.Amount) < 0 {
			t.Message = "Transação não autorizada - Saldo na conta insuficiente - food"
			t.Code = "500"
		} else {
			t.Code = "200"
			t.Message = "Transação autorizada"
		}
	case "5811":
		if (a.Amount_meal - t.Amount) < 0 {
			t.Message = "Transação não autorizada - Saldo na conta insuficiente - food"
			t.Code = "500"
		} else {
			t.Code = "200"
			t.Message = "Transação autorizada"
		}
	case "5812":
		if (a.Amount_meal - t.Amount) < 0 {
			t.Message = "Transação não autorizada - Saldo na conta insuficiente - food"
			t.Code = "500"
		} else {
			t.Code = "200"
			t.Message = "Transação autorizada"
		}
	default:
		if (a.Amount_cash - t.Amount) < 0 {
			t.Message = "Transação não autorizada - Saldo na conta insuficiente - food"
			t.Code = "500"
		} else {
			t.Code = "200"
			t.Message = "Transação autorizada"
		}
	}

	// caso tenha saldo suficiente retorna erro nulo
	return nil

}

// balanceOriginAccount atualiza o saldo da conta de origem
func (t *Transaction) balanceOriginAccount(app *app.App) error {

	// inicia o modo de transaction
	tx := app.DB.Client.Begin()

	// captura a conta de origem no DB
	origem := &Account{}
	if result := tx.First(&origem, &t.Accounttocredit_id); result.Error != nil {
		tx.Rollback()
		return errors.New("Conta de origem não encontrada")
	}

	// atualiza o saldo da conta de origem

	switch t.Mcc {
	case "5411":
		origem.Amount_food = origem.Amount_food - t.Amount
		if result := tx.Save(&origem); result.Error != nil {
			tx.Rollback()
			return errors.New("Erro ao atualizar saldo da conta de origem-food")
		}
		tx.Commit()
	case "5412":
		origem.Amount_food = origem.Amount_food - t.Amount
		if result := tx.Save(&origem); result.Error != nil {
			tx.Rollback()
			return errors.New("Erro ao atualizar saldo da conta de origem-food")
		}
		tx.Commit()
	case "5811":
		origem.Amount_meal = origem.Amount_meal - t.Amount
		if result := tx.Save(&origem); result.Error != nil {
			tx.Rollback()
			return errors.New("Erro ao atualizar saldo da conta de origem-meal")
		}
		tx.Commit()
	case "5812":
		origem.Amount_meal = origem.Amount_meal - t.Amount
		if result := tx.Save(&origem); result.Error != nil {
			tx.Rollback()
			return errors.New("Erro ao atualizar saldo da conta de origem-meal")
		}
		tx.Commit()
	default:
		origem.Amount_cash = origem.Amount_cash - t.Amount
		if result := tx.Save(&origem); result.Error != nil {
			tx.Rollback()
			return errors.New("Erro ao atualizar saldo da conta de origem-cash")
		}
		tx.Commit()
	}

	switch t.Code {
	case "500":
		origem.Amount_food = origem.Amount_food + t.Amount
		tx.Commit()
	default:
		origem.Amount_food = origem.Amount_food - t.Amount
		tx.Commit()
	}

	// atualização sem erros é comitada
	tx.Commit()

	// caso sucesso retorna erro nulo
	return nil

}

// balanceDestinationAccount atualiza o saldo da conta de destino
func (t *Transaction) balanceDestinationAccount(app *app.App) error {

	// inicia o modo de transaction
	tx := app.DB.Client.Begin()

	// captura a conta de destino no DB
	destino := &Account{}
	if result := tx.First(&destino, &t.Account_id); result.Error != nil {
		tx.Rollback()
		return errors.New("Conta de destino não encontrada")
	}

	// atualiza o saldo da conta de destino
	switch t.Mcc {
	case "5411":
		destino.Amount_food = destino.Amount_food - t.Amount
		if result := tx.Save(&destino); result.Error != nil {
			tx.Rollback()
			return errors.New("Erro ao atualizar saldo da conta de destino")
		}
		tx.Commit()
	case "5412":
		destino.Amount_food = destino.Amount_food - t.Amount
		if result := tx.Save(&destino); result.Error != nil {
			tx.Rollback()
			return errors.New("Erro ao atualizar saldo da conta de destino")
		}
		tx.Commit()
	case "5811":
		destino.Amount_meal = destino.Amount_meal - t.Amount
		if result := tx.Save(&destino); result.Error != nil {
			tx.Rollback()
			return errors.New("Erro ao atualizar saldo da conta de destino")
		}
		tx.Commit()
	case "5812":
		destino.Amount_meal = destino.Amount_meal - t.Amount
		if result := tx.Save(&destino); result.Error != nil {
			tx.Rollback()
			return errors.New("Erro ao atualizar saldo da conta de destino")
		}

		tx.Commit()
	default:
		destino.Amount_cash = destino.Amount_cash - t.Amount
		if result := tx.Save(&destino); result.Error != nil {
			tx.Rollback()
			return errors.New("Erro ao atualizar saldo da conta de destino")
		}
		tx.Commit()
	}

	// atualização sem erros é comitada
	tx.Commit()

	// caso sucesso retorna erro nulo
	return nil

}
