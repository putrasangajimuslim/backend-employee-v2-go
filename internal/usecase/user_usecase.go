// internal/usecase/user_usecase.go
package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/yourname/go-clean-arch/internal/domain/entity"
	"github.com/yourname/go-clean-arch/internal/domain/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"omitempty,min=2,max=100"`
	Email string `json:"email" validate:"omitempty,email"`
}

type UserUseCase struct {
	userRepo   repository.UserRepository
	jwtSecret  string
	jwtExpiry  time.Duration
}

func NewUserUseCase(userRepo repository.UserRepository, jwtSecret string) *UserUseCase {
	return &UserUseCase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: 24 * time.Hour,
	}
}

func (u *UserUseCase) Register(ctx context.Context, req RegisterRequest) (*entity.User, error) {
	existingUser, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (u *UserUseCase) Login(ctx context.Context, req LoginRequest) (string, error) {
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := u.generateToken(user.ID)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (u *UserUseCase) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	user.Password = ""
	return user, nil
}

func (u *UserUseCase) GetAll(ctx context.Context, page, limit int) ([]*entity.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return u.userRepo.FindAll(ctx, page, limit)
}

func (u *UserUseCase) Update(ctx context.Context, id uint, req UpdateUserRequest) (*entity.User, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.Email != "" && req.Email != user.Email {
		existingUser, err := u.userRepo.FindByEmail(ctx, req.Email)
		if err == nil && existingUser != nil {
			return nil, errors.New("email already in use")
		}
		user.Email = req.Email
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (u *UserUseCase) Delete(ctx context.Context, id uint) error {
	_, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}
	return u.userRepo.Delete(ctx, id)
}

func (u *UserUseCase) generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(u.jwtExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.jwtSecret))
}

func (u *UserUseCase) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(u.jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, errors.New("invalid token")
}
