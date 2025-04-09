package model

type City string

const (
	CityMoscow City = "Москва"
	CitySPB    City = "Санкт-Петербург"
	CityKazan  City = "Казань"
)

type PVZ struct {
	ID               string
	RegistrationDate string
	City             City
}
