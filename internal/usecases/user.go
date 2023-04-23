package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/practice-sem-2/user-service/internal/models"
	storage "github.com/practice-sem-2/user-service/internal/storages"
)

type UserCRUD interface {
	CreateUser(ctx context.Context, create *models.UserCreate) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, username string, fields models.UpdateFields) (*models.User, error)
	DeleteUser(ctx context.Context, username string) error
	GetManyUsers(ctx context.Context, usernames []string) ([]models.User, error)
}

type UserUseCase struct {
	store UserCRUD
}

func NewUserUseCase(store UserCRUD) *UserUseCase {
	return &UserUseCase{store: store}
}

func (u *UserUseCase) Create(ctx context.Context, user *models.UserCreate) (*models.User, error) {
	user.Password = hashPassword(user.Password)
	createdUser, err := u.store.CreateUser(ctx, user)
	return createdUser, err
}

func hashPassword(password string) string {
	hasher := sha256.New()
	return hex.EncodeToString(hasher.Sum([]byte(password)))
}

func verifyPassword(password string, hash string) bool {
	return hashPassword(password) == hash
}

func (u *UserUseCase) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return u.store.GetUserByUsername(ctx, username)
}

func (u *UserUseCase) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return u.store.GetUserByEmail(ctx, email)
}

func (u *UserUseCase) GetMany(ctx context.Context, usernames []string) ([]models.User, error) {
	return u.store.GetManyUsers(ctx, usernames)
}

func (u *UserUseCase) GetUserByCredentials(ctx context.Context, username string, password string) (*models.User, error) {
	user, err := u.GetByUsername(ctx, username)

	if err != nil {
		return nil, err
	}
	if verifyPassword(password, user.PasswordHash) {
		return user, nil
	} else {
		return nil, storage.ErrUserNotFound
	}
}

func (u *UserUseCase) Update(ctx context.Context, username string, fields models.UpdateFields) (*models.User, error) {
	return u.store.UpdateUser(ctx, username, fields)
}

func (u *UserUseCase) Delete(ctx context.Context, username string) error {
	return u.store.DeleteUser(ctx, username)
}
