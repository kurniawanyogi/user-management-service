package user

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"sync"
	"time"
	"user-management-service/common"
	"user-management-service/common/logger"
	"user-management-service/config"
	"user-management-service/model"
	"user-management-service/repository"
)

type IUserService interface {
	Register(ctx context.Context, payload model.RegistrationUserRequest) error
	Detail(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, payload model.UpdateUserRequest) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]model.User, error)
	Login(ctx context.Context, payload model.LoginRequest) (*string, *model.User, error)
	Authenticate(ctx context.Context, token string) error
}

type userService struct {
	common       common.IRegistry
	repositories repository.IRegistry
}

func NewUserService(common common.IRegistry, repositories repository.IRegistry) *userService {
	return &userService{
		common:       common,
		repositories: repositories,
	}
}

func (s *userService) Register(ctx context.Context, payload model.RegistrationUserRequest) error {
	// Implementasi goroutine untuk menjalankan 2 method ke repository
	var wg sync.WaitGroup

	// Buffer untuk menghindari goroutine leak
	errChan := make(chan error, 2)
	wg.Add(2)

	go func() {
		defer wg.Done()
		existingUsername, err := s.repositories.GetUserRepository().FindByUsername(ctx, payload.Username)
		if err != nil && err.Error() != common.ErrDataNotFound.Error() {
			errChan <- err
			return
		}
		if existingUsername != nil && existingUsername.Status == common.StatusUserActive {
			errChan <- common.ErrUsernameAlreadyTaken
			return
		}
	}()

	go func() {
		defer wg.Done()
		existingEmail, err := s.repositories.GetUserRepository().FindByEmail(ctx, payload.Email)
		if err != nil && err.Error() != common.ErrDataNotFound.Error() {
			errChan <- err
			return
		}
		if existingEmail != nil && existingEmail.Status == common.StatusUserActive {
			errChan <- common.ErrEmailAlreadyTaken
			return
		}
	}()

	// Close channel ketika goroutine selesai
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// mengecek error dari kedua goroutine di atas
	for err := range errChan {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return common.ErrFailHashPassword
	}

	_, err = s.repositories.GetUserRepository().Insert(
		ctx,
		model.User{
			Username:  payload.Username,
			Email:     payload.Email,
			Password:  string(hashedPassword),
			FirstName: payload.FirstName,
			LastName:  *payload.LastName,
			Status:    common.StatusUserActive,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) Detail(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.repositories.GetUserRepository().FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Update(ctx context.Context, payload model.UpdateUserRequest) error {
	existingUser, err := s.repositories.GetUserRepository().FindById(ctx, payload.Id)
	if err != nil {
		return err
	}

	existingUser.FirstName = payload.FirstName
	existingUser.LastName = *payload.LastName
	err = s.repositories.GetUserRepository().Update(ctx, *existingUser)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) Delete(ctx context.Context, id int64) error {
	existingUser, err := s.repositories.GetUserRepository().FindById(ctx, id)
	if err != nil {
		return err
	}

	if existingUser != nil && existingUser.Status == common.StatusUserInactive {
		return common.ErrUserAlreadyDeleted
	}

	existingUser.Status = common.StatusUserInactive
	err = s.repositories.GetUserRepository().Update(ctx, *existingUser)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) List(ctx context.Context) ([]model.User, error) {
	users, err := s.repositories.GetUserRepository().GetAll(ctx)
	if err != nil {
		logger.Error(ctx, "error getting all users", err, logger.Tag{Key: "logCtx", Value: ctx})
		return nil, err
	}

	// return empty slice instead nil
	if users == nil {
		return []model.User{}, nil
	}
	return users, nil
}

func (s *userService) Login(context context.Context, payload model.LoginRequest) (*string, *model.User, error) {
	existingUser, err := s.repositories.GetUserRepository().FindByUsername(context, payload.Username)
	if err != nil {
		return nil, nil, err
	}

	if existingUser != nil && existingUser.Status != common.StatusUserActive {
		return nil, nil, common.ErrUserAlreadyDeleted
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(payload.Password))
	if err != nil {
		return nil, nil, common.ErrInvalidPassword
	}

	token, err := generateJWT(existingUser)
	if err != nil {
		logger.Error(context, "error generating token", err)
		return nil, nil, err
	}
	return token, existingUser, nil
}

func (s *userService) Authenticate(ctx context.Context, token string) error {
	parsedToken, err := jwt.ParseWithClaims(token, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the token's signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Cold.SecretKey), nil
	})

	if err != nil {
		logger.Error(ctx, "Error parsing token", err)
		return common.ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(*model.Claims)
	if !ok || !parsedToken.Valid {
		return common.ErrInvalidToken
	}

	// validate expired token
	if claims.RegisteredClaims.ExpiresAt.Before(time.Now()) {
		return common.ErrInvalidToken
	}

	// validate username
	username := claims.Username
	if username == "" {
		return common.ErrInvalidToken
	}

	return nil
}

func generateJWT(user *model.User) (*string, error) {
	expirationTime := time.Now().Add(time.Hour * 24)
	claims := &model.Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: expirationTime,
			},
		},
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Cold.SecretKey))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}
