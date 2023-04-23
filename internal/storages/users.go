package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"github.com/practice-sem-2/user-service/internal/models"
	"reflect"
	"strings"
)

type UserStorage struct {
	db         Scope
	selectUser sq.SelectBuilder
	insertUser sq.InsertBuilder
	updateUser sq.UpdateBuilder
	deleteUser sq.DeleteBuilder
}

func NewUserStorage(db Scope) UserStorage {
	return UserStorage{
		db:         db,
		selectUser: sq.Select("*").From("users").PlaceholderFormat(sq.Dollar),
		insertUser: sq.Insert("users").PlaceholderFormat(sq.Dollar),
		updateUser: sq.Update("users").PlaceholderFormat(sq.Dollar),
		deleteUser: sq.Delete("users").PlaceholderFormat(sq.Dollar),
	}
}

type MissingUsersError struct {
	Usernames []string
}

func (e *MissingUsersError) Error() string {
	users := strings.Join(e.Usernames, ", ")
	return fmt.Sprintf("several users are not exists: %s", users)
}

var (
	ErrEmailAlreadyExists = errors.New("user with provided email already exists")
	ErrUserAlreadyExists  = errors.New("user with provided username already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCode        = errors.New("activation code is incorrect")
)

func (s *UserStorage) CreateUser(ctx context.Context, user *models.UserCreate) (*models.User, error) {
	builder := s.insertUser.
		Columns("username", "email", "is_active", "password_hash", "first_name", "last_name", "avatar_id").
		Values(user.Username, user.Email, false, user.Password, user.FirstName, user.LastName, user.AvatarID).
		Suffix("RETURNING *")

	query, args, err := builder.ToSql()

	if err != nil {
		return nil, err
	}

	row := s.db.QueryRowxContext(ctx, query, args...)
	var createdUser models.User

	err = row.StructScan(&createdUser)

	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok {
			if pgErr.ConstraintName == "users_pkey" {
				return nil, ErrUserAlreadyExists
			} else if pgErr.ConstraintName == "users_email_key" {
				return nil, ErrEmailAlreadyExists
			}
		}
		return nil, err
	}
	return &createdUser, err
}

func (s *UserStorage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	builder := s.selectUser.Where(sq.Eq{"username": username})
	query, args, err := builder.ToSql()

	if err != nil {
		return nil, err
	}
	var user models.User
	err = s.db.GetContext(ctx, &user, query, args...)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	} else {
		return &user, err
	}
}

func (s *UserStorage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	builder := s.selectUser.Where(sq.Eq{"email": email})
	query, args, err := builder.ToSql()

	if err != nil {
		return nil, err
	}
	var user models.User
	err = s.db.GetContext(ctx, &user, query, args...)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	} else {
		return &user, err
	}
}

func (s *UserStorage) GetManyUsers(ctx context.Context, usernames []string) ([]models.User, error) {
	q := make(sq.Or, len(usernames))
	for i, name := range usernames {
		q[i] = sq.Eq{"username": name}
	}
	query, args, err := s.selectUser.Where(q).ToSql()
	if err != nil {
		return nil, err
	}
	users := make([]models.User, 0, len(usernames))
	err = s.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, err
	}
	if len(users) < len(usernames) {
		missingCount := len(usernames) - len(users)
		hashmap := make(map[string]struct{}, missingCount)
		for _, user := range users {
			hashmap[user.Username] = struct{}{}
		}
		missing := make([]string, 0, missingCount)
		for _, username := range usernames {
			if _, ok := hashmap[username]; !ok {
				missing = append(missing, username)
			}
		}
		return users, &MissingUsersError{Usernames: missing}
	}
	return users, err
}

func (s *UserStorage) UpdateUser(ctx context.Context, username string, fields models.UpdateFields) (*models.User, error) {
	patchList := filterNil(fields)
	q := s.updateUser.Where(sq.Eq{"username": username}).Suffix("RETURNING *")

	for field, value := range patchList {
		q = q.Set(field, value)
	}

	user := &models.User{}
	query, args, err := q.ToSql()

	if err != nil {
		return nil, err
	}

	err = s.db.GetContext(ctx, user, query, args...)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		if pgErr, ok := err.(pgx.PgError); ok {
			if pgErr.ConstraintName == "users_pkey" {
				return nil, ErrUserAlreadyExists
			} else if pgErr.ConstraintName == "users_email_key" {
				return nil, ErrEmailAlreadyExists
			}
		}
		return nil, err
	}

	return user, nil
}

func (s *UserStorage) ActivateUser(ctx context.Context, username string, code string) error {
	query, args, err := sq.
		Select("username, code").From("users").
		LeftJoin("users_activation_codes c ON u.username = c.code").
		Where(sq.Eq{"username": username, "code": code}).
		Suffix("FOR UPDATE").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	var data []struct {
		Username string  `db:"username"`
		Code     *string `db:"code"`
	}

	err = s.db.SelectContext(ctx, &data, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrUserNotFound
	} else if err != nil {
		return err
	}

	if data[0].Code == nil {
		return ErrInvalidCode
	}

	query, args, err = sq.Update("users").
		Set("is_active", true).
		Where(sq.Eq{"username": username}).
		ToSql()

	if err != nil {
		return err
	}

	if _, err = s.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	// No need to reactivation user, so delete all activation codes
	query, args, err = sq.Delete("users_activation_code").
		Where(sq.Eq{"username": username}).
		ToSql()

	return nil
}

func (s *UserStorage) DeleteUser(ctx context.Context, username string) error {
	query, args, err := s.deleteUser.Where(sq.Eq{"username": username}).ToSql()

	res, err := s.db.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}

	if affected, err := res.RowsAffected(); err != nil && affected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func filterNil(arg interface{}) map[string]interface{} {
	av := reflect.ValueOf(arg)
	at := reflect.TypeOf(arg)
	result := make(map[string]interface{}, av.NumField())
	for i := 0; i < av.NumField(); i++ {
		if !av.Field(i).IsNil() {
			name := at.Field(i).Tag.Get("db")
			result[name] = av.Field(i).Interface()
		}
	}
	return result
}
