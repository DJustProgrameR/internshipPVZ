package model

type Role string

const (
	RoleEmployee  Role = "employee"
	RoleModerator Role = "moderator"
)

type User struct {
	ID       string
	Email    string
	Password string
	Role     Role
}
