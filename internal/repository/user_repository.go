package repository

import (
	"context"

	"main/internal/domain/entity"
	"main/internal/domain/repository"
	"main/internal/repository/model"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (repo *userRepository) Exist(ctx context.Context, email string) (bool, error) {
	var count int64

	err := repo.db.WithContext(ctx).Table(model.User{}.TableName()).
		Where("email = ?", email).
		Count(&count).Error
	if err != nil {
		return false, errors.Errorf("count email, err: %+v", err)
	}

	return count != 0, nil
}

func (repo *userRepository) Create(ctx context.Context, user *entity.User) error {
	u := model.NewUser(user)

	err := repo.db.WithContext(ctx).Create(u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.Errorf("create user, err: %+v", repository.ErrDuplicateKey)
		}

		return errors.Errorf("create user, err: %+v", err)
	}

	return nil
}

func (repo *userRepository) Get(ctx context.Context, opt repository.GetUserOption) (*entity.User, error) {
	var user model.User

	err := repo.db.WithContext(ctx).Table(user.TableName()).
		Scopes(repo.getUserOptionScope(opt)...).
		Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.Errorf("get user, err: %+v", repository.ErrDuplicateKey)
		}

		return nil, errors.Errorf("get user, err: %+v", err)
	}

	return user.ToEntity(), nil
}

func (userRepository) getUserOptionScope(opt repository.GetUserOption) []func(*gorm.DB) *gorm.DB {
	return []func(*gorm.DB) *gorm.DB{
		func(d *gorm.DB) *gorm.DB {
			if opt.ID != 0 {
				return d.Where("id = ?", opt.ID)
			}

			return d
		},
		func(d *gorm.DB) *gorm.DB {
			if len(opt.Email) != 0 {
				return d.Where("email = ?", opt.Email)
			}

			return d
		},
	}
}

func (repo *userRepository) Delete(ctx context.Context, opt repository.GetUserOption) error {
	err := repo.db.WithContext(ctx).Table(model.User{}.TableName()).
		Scopes(repo.getUserOptionScope(opt)...).
		Delete(&model.User{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return errors.Errorf("delete user, err: %+v", err)
	}

	return nil
}
