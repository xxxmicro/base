package models

type Model interface {
	GetTable() string
	GenerateID()
}