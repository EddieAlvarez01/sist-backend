package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/EddieAlvarez01/sist-backend/storage"
	"github.com/gocql/gocql"
	"gopkg.in/inf.v0"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	descriptionDebit = "Traslado de saldo hacia otra cuenta"
	descriptionCredit = "Deposito de saldo en su cuenta"
)

type ManageOperation struct {
	*storage.SistStorage
}

type Operation struct {
	AccountHolderID string `json:"account_holder_id,omitempty"`
	Month string `json:"month,omitempty"`
	AccountNumber string `json:"account_number,omitempty"`
	OperationID string `json:"operation_id,omitempty"`
	Type string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	Value string `json:"value,omitempty"`
	Date time.Time `json:"date,omitempty"`
	State string `json:"state,omitempty"`
	ReasonFailure string `json:"reason_failure,omitempty"`
}

type DTOOperation struct {
	AccountHolderID string `json:"account_holder_id,omitempty"`
	SourceAccount string `json:"source_account,omitempty"`
	DestinationAccount string `json:"destination_account,omitempty"`
	Value string `json:"value,omitempty"`
}

type TotalInstitution struct {
	InstitutionID string `json:"institution_id"`
	InstitutionName string `json:"institution_name"`
	Debits string `json:"debits"`
	Credits string `json:"credits"`
}

//Create new operation
func (m *ManageOperation) Create(operation *DTOOperation) (*Operation, error) {
	modelAccount := ManageAccount{SistStorage: m.SistStorage}
	account, err := modelAccount.GetByID(operation.SourceAccount)
	if err != nil {
		return nil, err
	}
	fromAccount := m.ConvertStringFloatToCurrency(account.Balance)
	fromOperation := m.ConvertStringFloatToCurrency(operation.Value)
	if fromAccount < fromOperation {
		if len(account.AssociatedAccounts) > 0 {
			var associatedAccounts []Account
			var associatedBalances []int
			accumulatedBalance := fromAccount
			for index, accountString := range account.AssociatedAccounts {
				associatedAccount, err := modelAccount.GetByID(accountString)
				if err != nil {
					return nil, err
				}
				associatedBalance := m.ConvertStringFloatToCurrency(associatedAccount.Balance)
				associatedAccounts = append(associatedAccounts, *associatedAccount)
				associatedBalances = append(associatedBalances, associatedBalance)
				accumulatedBalance += associatedBalance
				if accumulatedBalance >= fromOperation {
					break
				}
				if index == len(account.AssociatedAccounts) - 1 {
					return nil, errors.New(fmt.Sprintf("La cuenta no. %s no tiene fondos para suplir la operación", account.AccountNumber))
				}
			}
			for i, associatedAccount := range associatedAccounts {
				fromOperation, err = m.balanceUpdateLogic(&associatedAccount, modelAccount, associatedBalances[i], fromOperation, "DEBITO")
				if err != nil {
					return nil, err
				}
				if fromOperation <= 0 {
					break
				}
			}
		}else{
			return nil, errors.New(fmt.Sprintf("La cuenta no. %s no tiene fondos para suplir la operación", account.AccountNumber))
		}
	}
	if fromOperation > 0 {
		_, err = m.balanceUpdateLogic(account, modelAccount, fromAccount, fromOperation, "DEBITO")
		if err != nil {
			return nil, err
		}
	}
	op, err := m.sendOperationToDB(account, "DEBITO", operation.Value)
	if err != nil {
		return nil, err
	}
	err = m.updateInstitutionTotal(*account, "DEBITO", m.ConvertStringFloatToCurrency(operation.Value))
	if err != nil {
		return nil, err
	}
	account2, err := modelAccount.GetByID(operation.DestinationAccount)
	if err != nil {
		return nil, err
	}
	_, err = m.balanceUpdateLogic(account2, modelAccount, m.ConvertStringFloatToCurrency(account2.Balance), m.ConvertStringFloatToCurrency(operation.Value), "CREDITO")
	if err != nil {
		return nil, err
	}
	_, err = m.sendOperationToDB(account2, "CREDITO", operation.Value)
	if err != nil {
		return nil, err
	}
	err = m.updateInstitutionTotal(*account2, "CREDITO", m.ConvertStringFloatToCurrency(operation.Value))
	if err != nil {
		return nil, err
	}
	return op, nil
}

//Balance update logic
func (m *ManageOperation) balanceUpdateLogic(account *Account, modelAccount ManageAccount, accountBalance int, operationCost int, typeOperation string) (int, error) {
	var subtractOperation int
	if typeOperation == "DEBITO"{
		subtractOperation = operationCost - accountBalance
		if subtractOperation < 0 {
			accountBalance -= operationCost
		}else{
			accountBalance -= operationCost - subtractOperation
		}
	}else{
		accountBalance += operationCost
	}
	account.Balance = m.ConvertCurrencyToStringFloat(accountBalance)
	err := modelAccount.UpdateBalanceAccount(*account)
	if err != nil {
		return 0, err
	}
	return subtractOperation, nil
}

//Send operation to DB
func (m *ManageOperation) sendOperationToDB(account *Account, typeOperation string, value string) (*Operation, error) {
	var description string
	if typeOperation == "DEBITO" {
		description = descriptionDebit
	}else{
		description = descriptionCredit
	}
	operation := Operation{
		AccountHolderID: account.AccountHolderID,
		AccountNumber:   account.AccountNumber,
		OperationID:     gocql.TimeUUID().String(),
		Type:            typeOperation,
		Description:     description,
		Value:           value,
		Date:            time.Now(),
		State:           "OK",
		ReasonFailure:   "",
	}
	valueforDB, _ := new(inf.Dec).SetString(operation.Value)
	err := m.Session.Query(`INSERT INTO operaciones_por_cuentahabiente(cuentahabiente, no_cuenta, operacionID, tipo, descripcion, valor, fecha, estado, razon_falla)
								VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`, operation.AccountHolderID, operation.AccountNumber, operation.OperationID, operation.Type, operation.Description, valueforDB, operation.Date, operation.State, operation.ReasonFailure).WithContext(context.TODO()).Exec()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	err = m.Session.Query(`INSERT INTO operaciones_por_cuentahabiente_por_mes(cuentahabiente, mes, no_cuenta, operacionID, tipo, descripcion, valor, fecha, estado, razon_falla)
								VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, operation.AccountHolderID, strings.ToUpper(operation.Date.Month().String()), operation.AccountNumber, operation.OperationID, operation.Type, operation.Description, valueforDB, operation.Date, operation.State, operation.ReasonFailure).WithContext(context.TODO()).Exec()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &operation, nil
}

//Update a institution total
func (m *ManageOperation) updateInstitutionTotal(account Account, typeMovement string, value int) error {
	var totalInstitution TotalInstitution
	var debitInstitution inf.Dec
	var creditInstitution inf.Dec
	err := m.Session.Query(`SELECT institucion, nombre_institucion, debitos, creditos
								FROM totales_institucion_financiera
								WHERE institucion = ?`, account.InstitutionID).WithContext(context.TODO()).Consistency(gocql.One).Scan(&totalInstitution.InstitutionID, &totalInstitution.InstitutionName, &debitInstitution, &creditInstitution)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	totalInstitution.Debits = debitInstitution.String()
	totalInstitution.Credits = creditInstitution.String()
	var query *gocql.Query
	var amountExchange int
	if typeMovement == "DEBITO" {
		amountExchange = m.ConvertStringFloatToCurrency(totalInstitution.Debits)
		amountExchange += value
		totalInstitution.Debits = m.ConvertCurrencyToStringFloat(amountExchange)
		debitDec, _ := new(inf.Dec).SetString(totalInstitution.Debits)
		query = m.Session.Query(`UPDATE totales_institucion_financiera
				SET debitos = ?
				WHERE institucion = ?`, debitDec, totalInstitution.InstitutionID)
	}else{
		amountExchange = m.ConvertStringFloatToCurrency(totalInstitution.Credits)
		amountExchange += value
		totalInstitution.Credits = m.ConvertCurrencyToStringFloat(amountExchange)
		creditDec, _ := new(inf.Dec).SetString(totalInstitution.Credits)
		query = m.Session.Query(`UPDATE totales_institucion_financiera
				SET creditos = ?
				WHERE institucion = ?`, creditDec, totalInstitution.InstitutionID)
	}
	err = query.WithContext(context.TODO()).Exec()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

//Convert string float to int
func (m *ManageOperation) ConvertStringFloatToCurrency(value string) int {
	if value == "0" {
		return 0
	}
	dividedNumber := strings.Split(value, ".")
	entirePart, _ := strconv.Atoi(dividedNumber[0])
	entirePart *= 100
	decimalPart, _ := strconv.Atoi(dividedNumber[1])
	return entirePart + decimalPart
}

//Convert int to string float
func (m *ManageOperation) ConvertCurrencyToStringFloat(value int) string {
	if value == 0 {
		return "0.00"
	}
	numberString := strconv.Itoa(value)
	if value < 100 {
		if value >= 10 {
			return fmt.Sprintf("0.%s", numberString)
		}
		return fmt.Sprintf("0.0%s", numberString)
	}
	values := []string{numberString[:len(numberString) - 2], numberString[len(numberString) - 2:]}
	return strings.Join(values, ".")
}

//Get all operations by account holder ID
func (m *ManageOperation) GetAllByAccountHolderID(id string, month string, op int) ([]Operation, error) {
	operations := make([]Operation, 0)
	var query *gocql.Query
	if op == 1 {
		query = m.Session.Query(`SELECT cuentahabiente, no_cuenta, operacionID, tipo, descripcion, valor, fecha, estado, razon_falla
								FROM operaciones_por_cuentahabiente
								WHERE cuentahabiente = ?`, id).WithContext(context.TODO())
	}else{
		query = m.Session.Query(`SELECT cuentahabiente, mes, no_cuenta, operacionID, tipo, descripcion, valor, fecha, estado, razon_falla
								FROM operaciones_por_cuentahabiente_por_mes
								WHERE cuentahabiente = ? AND mes = ?`, id, month).WithContext(context.TODO())
	}
	scanner := query.Iter().Scanner()
	for scanner.Next() {
		var operation Operation
		var valueDec inf.Dec
		var err error
		if op == 1 {
			err = scanner.Scan(&operation.AccountHolderID, &operation.AccountNumber, &operation.OperationID, &operation.Type, &operation.Description, &valueDec, &operation.Date, &operation.State, &operation.ReasonFailure)
		}else{
			err = scanner.Scan(&operation.AccountHolderID, &operation.Month, &operation.AccountNumber, &operation.OperationID, &operation.Type, &operation.Description, &valueDec, &operation.Date, &operation.State, &operation.ReasonFailure)
		}
		if  err != nil {
			log.Println(err.Error())
			return nil, err
		}
		operation.Value = valueDec.String()
		operations = append(operations, operation)
	}
	return operations, nil
}

//Get institution totals
func (m *ManageOperation) GetInstitutionTotals(id string) (*TotalInstitution, error) {
	var total TotalInstitution
	var debitsDec inf.Dec
	var creditsDec inf.Dec
	err := m.Session.Query(`SELECT institucion, nombre_institucion, debitos, creditos
								FROM totales_institucion_financiera
								WHERE institucion = ?`, id).WithContext(context.TODO()).Consistency(gocql.One).Scan(&total.InstitutionID, &total.InstitutionName, &debitsDec, &creditsDec)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	total.Debits = debitsDec.String()
	total.Credits = creditsDec.String()
	return &total, nil
}


