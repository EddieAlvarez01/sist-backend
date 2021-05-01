package models

import (
	"context"
	"github.com/EddieAlvarez01/sist-backend/storage"
	"github.com/gocql/gocql"
	"gopkg.in/inf.v0"
	"log"
)

type ManageInstitution struct {
	*storage.SistStorage
}

type Institution struct {
	ID string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Address string `json:"address,omitempty"`
	Phone string `json:"phone,omitempty"`
}

//Create a new institution
func (m *ManageInstitution) Create(institution *Institution) (*Institution, error) {
	institution.ID = gocql.TimeUUID().String()
	err := m.Session.Query(`INSERT INTO institucion_financiera_por_ID(id, nombre, tipo, direccion, telefono)
								VALUES(?, ?, ?, ?, ?)`, institution.ID, institution.Name, institution.Type, institution.Address, institution.Phone).WithContext(context.TODO()).Exec()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	initialTotal, _ := new(inf.Dec).SetString("0.00")
	err = m.Session.Query(`INSERT INTO totales_institucion_financiera(institucion, nombre_institucion, debitos, creditos)
								VALUES(?, ?, ?, ?)`, institution.ID, institution.Name, initialTotal, initialTotal).WithContext(context.TODO()).Exec()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return institution, nil
}

//Get all institutions
func (m *ManageInstitution) GetAll() ([]Institution, error) {
	scanner := m.Session.Query(`SELECT id, nombre, tipo, direccion, telefono FROM institucion_financiera_por_ID`).WithContext(context.TODO()).Iter().Scanner()
	institutions := make([]Institution, 0)
	for scanner.Next() {
		var institution Institution
		err := scanner.Scan(&institution.ID, &institution.Name, &institution.Type, &institution.Address, &institution.Phone)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		institutions = append(institutions, institution)
	}
	return institutions, nil
}
