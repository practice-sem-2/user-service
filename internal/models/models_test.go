package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserCreate_CanValidateMinimalData(t *testing.T) {
	u := UserCreate{
		Username: "joe",
		Password: "qwerty1",
		Email:    "joe@example.com",
	}

	err := Validate.Struct(u)
	assert.Nil(t, err, "Should not return any errors")
}

func TestUserCreate_CanValidateRichData(t *testing.T) {
	avatarID := "09f28b68-697b-48fd-809d-dcb0c460d992"
	u := UserCreate{
		Username:  "joe",
		Password:  "qwerty1",
		Email:     "joe@example.com",
		FirstName: "John",
		LastName:  "Doe",
		AvatarID:  &avatarID,
	}

	err := Validate.Struct(u)
	assert.Nil(t, err, "Should not return any errors")
}

func TestUserCreate_ReturnsErrorIfDataIncorrect(t *testing.T) {
	avatarID := "asd"
	u := UserCreate{
		Username:  "j",
		Email:     "not_a_correct_email",
		FirstName: "John",
		// Last is too long
		LastName: "123456789_123456789_123456789_123456789",
		AvatarID: &avatarID,
	}

	err := Validate.Struct(u)
	hasError := map[string]bool{
		"Username": false,
		"Password": false,
		"LastName": false,
		"AvatarID": false,
	}

	for _, e := range err.(validator.ValidationErrors) {
		switch e.Field() {
		case "Username":
			hasError["Username"] = true
			assert.True(t, e.Tag() == "min", "Should fail, because username must contain at least 3 characters")
		case "Password":
			hasError["Password"] = true
			assert.True(t, e.Tag() == "required", "Should fail, because username must contain at least 3 characters")
		case "LastName":
			hasError["LastName"] = true
			assert.True(t, e.Tag() == "max", "Should fail, because username must contain at least 3 characters")
		case "AvatarID":
			hasError["AvatarID"] = true
			assert.True(t, e.Tag() == "uuid", "Should fail, because username must contain at least 3 characters")
		default:
			assert.Failf(t, "Field %s considered to be correct, but got: %s", e.Field(), e.Error())
		}
	}

	for key, hasErr := range hasError {
		assert.True(t, hasErr, "Should return error in %s", key)
	}
}

func TestUser_CanValidateModel(t *testing.T) {
	avatarID := "09f28b68-697b-48fd-809d-dcb0c460d992"
	u := User{
		Username:     "joe",
		PasswordHash: "b6ad34b0b6b7e38f878a513b3f7927ebeb4cffb01aeb6d9fd9f9ad67fbc76517",
		Email:        "john@example.com",
		FirstName:    "John",
		LastName:     "Doe",
		AvatarID:     &avatarID,
		IsActive:     false,
	}
	err := Validate.Struct(u)
	if err != nil {
		assert.Failf(t, "Should not return any validation errors but got %s", err.Error())
	}
}

func TestUpdateFields_CanValidateCorrectData(t *testing.T) {
	password := "qwerty1"
	email := "john@example.com"
	avatarId := "09f28b68-697b-48fd-809d-dcb0c460d992"

	u := UpdateFields{
		Password:  &password,
		Email:     &email,
		FirstName: nil,
		LastName:  nil,
		AvatarID:  &avatarId,
	}

	err := Validate.Struct(u)
	assert.Nilf(t, err, "Should return no errors, if data is correct")
}
