package usecase

import storage "github.com/practice-sem-2/user-service/internal/storages"

type UseCase struct {
	Users *UserUseCase
}

func NewUseCase(store *storage.Storage) *UseCase {
	return &UseCase{
		Users: NewUserUseCase(store),
	}
}
