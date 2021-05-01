package models

import (
	"context"
	"github.com/EddieAlvarez01/sist-backend/storage"
	"log"
)

type ManageAccount struct {
	*storage.SistStorage
}

type Account struct {
	AccountHolderID string `json:"account_holder_id,omitempty"`
	InstitutionID string `json:"institution_id,omitempty"`
	AccountNumber string `json:"account_number,omitempty"`
	InstitutionName string `json:"institution_name,omitempty"`
	Type string `json:"type,omitempty"`
	Balance string `json:"balance,omitempty"`
	AssociatedAccounts []string `json:"associated_accounts,omitempty"`
}

//Get all accounts of account holder
func (m *ManageAccount) GetAllByAccountHolder(id string) ([]Account, error) {
	scanner := m.Session.Query(`SELECT cuentahabiente, institucion, no_cuenta, nombre_institucion, tipo, saldo, cuentas_asociadas
									FROM cuentas_de_cuentahabiente
									WHERE cuentahabiente = ?`, id).WithContext(context.TODO()).Iter().Scanner()
	accounts := make([]Account, 0)
	for scanner.Next() {
		var account Account
		if err := scanner.Scan(&account.AccountHolderID, &account.InstitutionID, &account.AccountNumber, &account.InstitutionName, &account.Type, &account.Balance, &account.AssociatedAccounts); err != nil {
			log.Println(err.Error())
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}
