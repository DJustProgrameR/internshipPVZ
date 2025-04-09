package repository

import "github.com/DJustProgrameR/internshipPVZ/applicationCore/domain/model"

type PVZRepo interface {
	Create(pvz *model.PVZ) error
	GetAllWithFilter(start, end string, page, limit int) ([]*model.PVZ, error)
	FindByID(id string) (*model.PVZ, error)
}
