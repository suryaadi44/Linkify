package service

import (
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/suryaadi44/linkify/internal/user/dto"
	entity "github.com/suryaadi44/linkify/internal/user/entitiy"
	"github.com/suryaadi44/linkify/internal/user/repository"
	"github.com/suryaadi44/linkify/pkg/utils"
)

type UserService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (u UserService) CreateUser(ctx context.Context, user dto.RegisterRequest) error {
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	userEntity := entity.User{
		UID:      uuid.NewString(),
		Username: user.Username,
		Email:    user.Email,
		Password: hash,
		Rank:     0,
		Created:  time.Now(),
	}

	err = u.repository.CreateUser(ctx, userEntity)
	if err != nil {
		log.Println("[User] Error creating user :", err)
		return err
	}

	return nil
}

func (u UserService) IsEmailExists(ctx context.Context, email string) bool {
	return u.repository.IsEmailExists(ctx, email)
}

func (u UserService) IsUsernameExists(ctx context.Context, username string) bool {
	return u.repository.IsUsernameExists(ctx, username)
}

func (u UserService) AuthenticateUser(ctx context.Context, user dto.LoginRequest) (*jwt.Token, error) {
	savedUser, err := u.repository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	if !utils.CheckPasswordHash(user.Password, savedUser.Password) {
		return nil, nil
	}

	//Create jwt claims
	claims := jwt.MapClaims{
		"uid":        savedUser.UID,
		"permission": savedUser.Rank,
		"exp":        time.Now().Add(time.Hour * 2).Unix(),
	}

	//Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token, nil
}
