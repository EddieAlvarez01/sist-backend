package models

import (
	"context"
	"github.com/EddieAlvarez01/sist-backend/storage"
	"github.com/gocql/gocql"
	"gopkg.in/inf.v0"
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

//Create a new account for an account holder
func (m *ManageAccount) Create(account *Account) (*Account, error) {
	account.AccountNumber = gocql.TimeUUID().String()
	balance, _ := new(inf.Dec).SetString(account.Balance)
	err := m.Session.Query(`INSERT INTO cuentas_de_cuentahabiente(cuentahabiente, institucion, no_cuenta, nombre_institucion, tipo, saldo, cuentas_asociadas)
								VALUES(?, ?, ?, ?, ?, ?, ?)`, account.AccountHolderID, account.InstitutionID, account.AccountNumber, account.InstitutionName, account.Type, balance, account.AssociatedAccounts).WithContext(context.TODO()).Exec()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return account, nil
}

//Get account by ID
func (m *ManageAccount) GetByID(id string) (*Account, error) {
	var accountNumber string
	var accountHolderID string
	var institutionName string
	err := m.Session.Query(`SELECT no_cuenta, cuentahabiente, institucion
								FROM cuenta_por_ID
								WHERE no_cuenta = ?`, id).WithContext(context.TODO()).Consistency(gocql.One).Scan(&accountNumber, &accountHolderID, &institutionName)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	var account Account
	var balance inf.Dec
	err = m.Session.Query(`SELECT cuentahabiente, institucion, no_cuenta, nombre_institucion, tipo, saldo, cuentas_asociadas
								FROM cuentas_de_cuentahabiente
								WHERE cuentahabiente = ? AND institucion = ? AND no_cuenta = ?`, accountHolderID, institutionName, accountNumber).WithContext(context.TODO()).Consistency(gocql.One).Scan(&account.AccountHolderID, &account.InstitutionID, &account.AccountNumber, &account.InstitutionName, &account.Type, &balance, &account.AssociatedAccounts)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	account.Balance = balance.String()
	return &account, nil
}

//Update Balance account
func (m *ManageAccount) UpdateBalanceAccount(account Account) error {
	balanceDec, _ := new(inf.Dec).SetString(account.Balance)
	err := m.Session.Query(`UPDATE cuentas_de_cuentahabiente
								SET saldo = ?
								WHERE cuentahabiente = ? AND institucion = ? AND no_cuenta = ?`, balanceDec, account.AccountHolderID, account.InstitutionID, account.AccountNumber).WithContext(context.TODO()).Exec()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
