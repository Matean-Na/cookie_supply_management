package services

import (
	"cookie_supply_management/core/config"
	"cookie_supply_management/internal/constants"
	"cookie_supply_management/internal/dto"
	"cookie_supply_management/internal/models"
	"cookie_supply_management/internal/repositories"
	"cookie_supply_management/pkg/security"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type UserServiceInterface interface {
	Registration(input dto.UserCreateDTO) (token string, httpCode int, err error)
	Login(input dto.UserLoginDTO) (token string, httpCode int, err error)
	Logout(username string) error
	GetToken(userName string) (token string, err error)
	GetByUsername(userName string) (*models.User, error)
	Update(username string, input dto.UserUpdateDTO) (httpCode int, err error)
}

type UserService struct {
	repo repositories.UserRepositoryInterface
	conf *config.Config
}

func NewUserService(repo repositories.UserRepositoryInterface, conf *config.Config) *UserService {
	return &UserService{
		repo: repo,
		conf: conf,
	}
}

func (s *UserService) Update(username string, input dto.UserUpdateDTO) (httpCode int, err error) {
	user, err := s.repo.SelectUserByUsername(username)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if input.NewPassword != "" && input.OldPassword != "" {
		if err = security.VerifyPassword(user.Password, input.OldPassword); err != nil {
			return http.StatusBadRequest, errors.New("неверный пароль")
		}

		if input.NewPassword != input.NewPasswordConfirm {
			return http.StatusBadRequest, errors.New("пароли не совпадают")
		}

		hashedPassword, err := security.Hash(input.NewPassword)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		user.Password = hashedPassword
	}

	if input.NewUserName != "" {
		user.Username = input.NewUserName
	}

	if input.NewRole != "" {
		user.Role = input.NewRole
	}

	if err = s.repo.UpdateUser(user); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *UserService) Logout(username string) error {
	if err := s.repo.DeleteToken(username); err != nil {
		return err
	}

	return nil
}

func (s *UserService) Login(input dto.UserLoginDTO) (token string, httpCode int, err error) {
	var user models.User

	user, err = s.repo.SelectUserByUsername(input.UserName)
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	if err = security.VerifyPassword(user.Password, input.Password); err != nil {
		return "", http.StatusBadRequest, errors.New("неверный пароль")
	}

	expiresIn := time.Duration(s.conf.Token.DurationHour) * time.Hour

	token, err = security.GenerateToken(user.Username, s.conf.Token.SecretKey, expiresIn)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	if err = s.repo.SetToken(user.Username, token); err != nil {
		return "", http.StatusInternalServerError, err
	}

	return token, http.StatusOK, nil
}

func (s *UserService) Registration(input dto.UserCreateDTO) (token string, httpCode int, err error) {
	if input.Role != constants.Admin && input.Role != constants.Accountant {
		return "", http.StatusBadRequest, errors.New(fmt.Sprintf("роль должна быть либо «%s», либо «%s»", constants.Admin, constants.Accountant))
	}

	if input.Password != input.PasswordConfirm {
		return "", http.StatusBadRequest, errors.New("пароли не совпадают")
	}

	if err = s.repo.SelectExistByUserName(input.UserName); err == nil {
		return "", http.StatusBadRequest, errors.New("пользователь существует")
	}

	input.Password, err = security.Hash(input.Password)

	if err = s.repo.InsertUser(input.UserName, input.Password, input.Role); err != nil {
		return "", http.StatusInternalServerError, err
	}

	expiresIn := time.Duration(s.conf.Token.DurationHour) * time.Hour

	token, err = security.GenerateToken(input.UserName, s.conf.SecretKey, expiresIn)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	if err = s.repo.SetToken(input.UserName, token); err != nil {
		return "", http.StatusInternalServerError, err
	}

	return token, http.StatusOK, nil
}

func (s *UserService) GetToken(userName string) (token string, err error) {
	token, err = s.repo.GetToken(userName)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) GetByUsername(userName string) (*models.User, error) {
	user, err := s.repo.SelectUserByUsername(userName)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
