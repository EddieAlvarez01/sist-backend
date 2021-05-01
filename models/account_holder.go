package models

import (
	"context"
	"errors"
	"github.com/EddieAlvarez01/sist-backend/storage"
	"github.com/gocql/gocql"
	"log"
)

type ManageAccountHolder struct {
	*storage.SistStorage
}

type AccountHolder struct {
	ID string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	LastName string `json:"last_name,omitempty"`
	Email string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Role string `json:"role,omitempty"`
}

//Create a new user
func (m *ManageAccountHolder) Create(accountHolder *AccountHolder) (*AccountHolder, error) {
	idToCrate := gocql.TimeUUID()
	err := m.Session.Query(`INSERT INTO cuentahabiente_por_ID (id, nombre, apellido, correo, contrasena, rol)
									VALUES(?, ?, ?, ?, ?, ?)`,
									idToCrate, accountHolder.Name, accountHolder.LastName, accountHolder.Email, accountHolder.Password, accountHolder.Role).WithContext(context.TODO()).Exec()
	if err != nil {
		return nil, err
	}
	accountHolder.ID = idToCrate.String()
	return accountHolder, nil
}

//Get user by Email and password
func (m *ManageAccountHolder) GetUserByEmailAndPassword(email, password string)(*AccountHolder, error) {
	var id string
	var passwordDB string
	err := m.Session.Query(`SELECT id, contrasena
								FROM cuentahabiente_por_correo
								WHERE correo = ?`, email).WithContext(context.TODO()).Consistency(gocql.One).Scan(&id, &passwordDB)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if id == "" || passwordDB == "" {
		return nil, errors.New("User not found")
	}
	if passwordDB != password {
		return nil, errors.New("Incorrect credentials")
	}
	var accountHolder AccountHolder
	err = m.Session.Query(`SELECT id, nombre, apellido, correo, rol
								FROM cuentahabiente_por_ID
								WHERE id = ?`, id).WithContext(context.TODO()).Consistency(gocql.One).Scan(&accountHolder.ID, &accountHolder.Name, &accountHolder.LastName, &accountHolder.Email, &accountHolder.Role)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &accountHolder, nil
}

//Get all accounts holder
func (m *ManageAccountHolder) GetAll() ([]AccountHolder, error) {
	scanner := m.Session.Query(`SELECT id, nombre, apellido, correo, rol
								FROM cuentahabiente_por_ID`).WithContext(context.TODO()).Iter().Scanner()
	accountsHolder := make([]AccountHolder, 0)
	for scanner.Next() {
		var accountHolder AccountHolder
		err := scanner.Scan(&accountHolder.ID, &accountHolder.Name, &accountHolder.LastName, &accountHolder.Email, &accountHolder.Role)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		accountsHolder = append(accountsHolder, accountHolder)
	}
	return accountsHolder, nil
}
