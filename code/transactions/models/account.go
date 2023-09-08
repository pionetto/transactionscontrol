package models

import (
	"errors"
	"time"

	"cajueiro/pkg/app"
	"cajueiro/pkg/secret"

	"gorm.io/gorm"
)

// BeforeCreate hook do gorm para gerar uuid no create
func (a *Account) BeforeCreate(tx *gorm.DB) (err error) {
	a.Secret, err = secret.HashPassword(a.Secret)
	if err != nil {
		return errors.New("Erro ao criptografar senha")
	}
	return
}

// Account modelo para conta do usuário
type Account struct {
	gorm.Model  `json:"-"`
	ID          int            `json:"id" gorm:"not null"`
	CPF         string         `gorm:"unique" json:"cpf" validate:"required,len=11"`
	Secret      string         `json:"secret" validate:"required"`
	Amount_food float64        `json:"amount_food" validate:"required"`
	Amount_meal float64        `json:"amount_meal" validate:"required"`
	Amount_cash float64        `json:"amount_cash" validate:"required"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted"`
	Transaction []Transaction  `json:"-" gorm:"foreignKey:Account_id"`
}

// CreateAccount cria uma conta de usuário
func (a *Account) CreateAccount(app *app.App) (*Account, error) {

	account := &Account{
		ID:          a.ID,
		CPF:         a.CPF,
		Secret:      a.Secret,
		Amount_food: a.Amount_food,
		Amount_meal: a.Amount_meal,
		Amount_cash: a.Amount_cash,
		CreatedAt:   a.CreatedAt,
		Transaction: a.Transaction,
	}

	result := app.DB.Client.Create(account)

	if result.Error != nil {
		return nil, errors.New("Erro ao criar a conta")
	}

	return account, nil

}
