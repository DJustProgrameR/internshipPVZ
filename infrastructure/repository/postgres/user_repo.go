package postgres

import (
	"database/sql"
	"github.com/DJustProgrameR/internshipPVZ/applicationCore/domain/model"
	repo "github.com/DJustProgrameR/internshipPVZ/applicationCore/domain/repository"
	sqr "github.com/Masterminds/squirrel"
)

type userRepo struct {
	db *sql.DB
	qb sqr.StatementBuilderType
}

func NewUserRepo(db *sql.DB) repo.UserRepo {
	return &userRepo{
		db: db,
		qb: sqr.StatementBuilder.PlaceholderFormat(sqr.Dollar),
	}
}

func (r *userRepo) Create(user *model.User) error {
	_, err := r.qb.Insert("users").
		Columns("id", "email", "password", "role").
		Values(user.ID, user.Email, user.Password, user.Role).
		RunWith(r.db).Exec()
	return err
}

func (r *userRepo) FindByEmail(email string) (*model.User, error) {
	row := r.qb.Select("id", "email", "password", "role").
		From("users").
		Where(sqr.Eq{"email": email}).
		RunWith(r.db).QueryRow()

	u := &model.User{}
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Role)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepo) FindByID(id string) (*model.User, error) {
	row := r.qb.Select("id", "email", "password", "role").
		From("users").
		Where(sqr.Eq{"id": id}).
		RunWith(r.db).QueryRow()

	u := &model.User{}
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Role)
	if err != nil {
		return nil, err
	}
	return u, nil
}
