package models

type Error struct {
	Code uint16 `json:"code"`
	Message string `json:"message"`
}
