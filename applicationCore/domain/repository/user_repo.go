package repository

import "github.com/DJustProgrameR/internshipPVZ/applicationCore/domain/model"

type ProductRepo interface {
	Add(product *model.Product) error
	DeleteLastFromReception(receptionID string) error
	ListByReception(receptionID string) ([]*model.Product, error)
}
