package repository

import (
	"context"

	"main/internal/domain/entity"
	"main/internal/domain/repository"
	"main/internal/repository/conn"
	"main/internal/repository/query"

	"github.com/pkg/errors"
)

type userRepository struct {
	db *conn.Dao
}

func NewUserRepository(db *conn.Dao) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (repo *userRepository) Exist(ctx context.Context, email string) (bool, error) {
	count, err := repo.db.CountUser(ctx, email)
	if err != nil {
		return false, errors.Errorf("count email, err: %+v", err)
	}

	return count != 0, nil
}

func (repo *userRepository) Create(ctx context.Context, user *entity.User) error {
	id, err := repo.db.CreateUser(ctx, query.CreateUserParams{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		if conn.IsDuplicateKeyError(err) {
			return errors.Errorf("create user, err: %+v", repository.ErrDuplicateKey)
		}

		return errors.Errorf("create user, err: %+v", err)
	}

	user.ID = id

	return nil
}

func (repo *userRepository) Get(ctx context.Context, opt repository.GetUserOption) (*entity.User, error) {
	user, err := repo.db.GetUserByEmail(ctx, opt.Email)
	if err != nil {
		if conn.IsNotFoundError(err) {
			return nil, errors.Errorf("get user, err: %+v", repository.ErrNotFound)
		}

		return nil, errors.Errorf("get user, err: %+v", err)
	}

	return &entity.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (repo *userRepository) Delete(ctx context.Context, userID int64) error {
	err := repo.db.DeleteUser(ctx, userID)
	if err != nil {
		if conn.IsNotFoundError(err) {
			return nil
		}

		return errors.Errorf("delete user, err: %+v", err)
	}

	return nil
}
