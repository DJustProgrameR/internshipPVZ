package repository

import "github.com/DJustProgrameR/internshipPVZ/applicationCore/domain/model"

type ReceptionRepo interface {
	Create(reception *model.Reception) error
	GetLastByPVZ(pvzID string) (*model.Reception, error)
	Close(id string) error
	GetAllByPVZ(pvzID string) ([]*model.Reception, error)
}
