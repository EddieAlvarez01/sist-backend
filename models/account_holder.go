package models

import (
	"context"
	"github.com/EddieAlvarez01/sist-backend/storage"
	"github.com/gocql/gocql"
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
