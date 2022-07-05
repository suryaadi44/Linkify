package service

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/suryaadi44/linkify/internal/constant"
	linkRepositoryPkg "github.com/suryaadi44/linkify/internal/link/repository"
	"github.com/suryaadi44/linkify/internal/user/dto"
	entity "github.com/suryaadi44/linkify/internal/user/entitiy"
	userRepositoryPkg "github.com/suryaadi44/linkify/internal/user/repository"
	"github.com/suryaadi44/linkify/pkg/utils"
)

type UserService struct {
	userRepository userRepositoryPkg.UserRepository
	linkRepository linkRepositoryPkg.LinkRepository
}

func NewUserService(repository userRepositoryPkg.UserRepository, linkRepository linkRepositoryPkg.LinkRepository) *UserService {
	return &UserService{
		userRepository: repository,
		linkRepository: linkRepository,
	}
}

func (u UserService) CreateUser(ctx context.Context, user dto.RegisterRequest) error {
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Println("[User] Error hashing password :", err)
		return err
	}

	userEntity := entity.User{
		UID:      uuid.NewString(),
		Username: user.Username,
		Email:    user.Email,
		Password: hash,
		Picture:  constant.PICTURE_DEFAULT,
		Rank:     0,
		Created:  time.Now(),
	}

	err = u.userRepository.CreateUser(ctx, userEntity)
	if err != nil {
		log.Println("[User] Error creating user :", err)
		return err
	}

	err = u.linkRepository.CreateDefaultLink(ctx, user.Username)
	if err != nil {
		log.Println("[User] Error creating default link :", err)
		return err
	}

	return nil
}

func (u UserService) IsEmailExists(ctx context.Context, email string) bool {
	return u.userRepository.IsEmailExists(ctx, email)
}

func (u UserService) IsUsernameExists(ctx context.Context, username string) bool {
	return u.userRepository.IsUsernameExists(ctx, username)
}

func (u UserService) AuthenticateUser(ctx context.Context, user dto.LoginRequest) (*string, error) {
	savedUser, err := u.userRepository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		log.Println("[User] Error getting user :", err)
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

	//Sign token
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Println("[User] Error signing token :", err)
		return nil, err
	}

	return &tokenString, nil
}

func (u UserService) GetUserPictureByUsername(ctx context.Context, username string) (string, error) {
	picture, err := u.userRepository.GetUserPictureByUsername(ctx, username)
	if err != nil {
		log.Println("[User] Error getting user picture :", err)
		return "", err
	}

	return picture, nil
}
