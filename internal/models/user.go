package models

import "github.com/go-playground/validator/v10"

type UserCreate struct {
	Username  string  `db:"username" validate:"required,min=2,max=40"`
	Password  string  `db:"password" validate:"required,min=5,max=64"`
	Email     string  `db:"email" validate:"required,min=3,max=64"`
	FirstName string  `db:"first_name" validate:"omitempty,max=32"`
	LastName  string  `db:"last_name" validate:"omitempty,max=32"`
	AvatarID  *string `db:"avatar_id" validate:"omitempty,uuid"`
}

type User struct {
	Username     string  `db:"username" validate:"required,min=3,max=40"`
	PasswordHash string  `db:"password_hash" validate:"required"`
	Email        string  `db:"email" validate:"required,email,max=64"`
	FirstName    string  `db:"first_name" validate:"omitempty,max=32"`
	LastName     string  `db:"last_name" validate:"omitempty,max=32"`
	AvatarID     *string `db:"avatar_id" validate:"omitempty,uuid"`
	IsActive     bool    `db:"is_active" validate:""`
}

type UpdateFields struct {
	Password  *string `db:"password_hash" validate:"omitempty,min=5,max=64"`
	Email     *string `db:"email" validate:"omitempty,email,min=3,max=64"`
	FirstName *string `db:"first_name" validate:"omitempty,max=32"`
	LastName  *string `db:"last_name" validate:"omitempty,max=32"`
	AvatarID  *string `db:"avatar_id" validate:"omitempty,uuid"`
}

var Validate = validator.New()
