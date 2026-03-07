// internal/infrastructure/repository/user_repository.go
package repository

import (
	"context"
	"time"

	"github.com/yourname/go-clean-arch/internal/domain/entity"
	"github.com/yourname/go-clean-arch/internal/domain/repository"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:100;uniqueIndex;not null"`
	Password  string    `gorm:"size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return "users"
}

func (u *User) ToEntity() *entity.User {
	return &entity.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func ToEntityFromUsers(users []User) []*entity.User {
	result := make([]*entity.User, len(users))
	for i, u := range users {
		result[i] = u.ToEntity()
	}
	return result
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	u := &User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	var u User
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	return u.ToEntity(), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var u User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return u.ToEntity(), nil
}

func (r *userRepository) FindAll(ctx context.Context, page, limit int) ([]*entity.User, int64, error) {
	var users []User
	var total int64

	offset := (page - 1) * limit

	if err := r.db.WithContext(ctx).Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return ToEntityFromUsers(users), total, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"name":   user.Name,
		"email":  user.Email,
	}).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&User{}, id).Error
}
