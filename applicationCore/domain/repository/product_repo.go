package repository

import "github.com/DJustProgrameR/internshipPVZ/applicationCore/domain/model"

type UserRepo interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id string) (*model.User, error)
}
