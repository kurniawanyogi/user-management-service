package user

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"sync"
	"testing"
	"time"
	"user-management-service/common"
	"user-management-service/config"
	mocksUserRepo "user-management-service/mocks/repository/user"
	"user-management-service/model"
	"user-management-service/repository"
)

type listMock struct {
	database sqlmock.Sqlmock
	sqlmock  sqlmock.Sqlmock
	userRepo mocksUserRepo.IUserRepository
}
type UserServiceSuite struct {
	suite.Suite
	mocks        listMock
	repositories repository.IRegistry
	UserService  IUserService
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}

func (s *UserServiceSuite) SetupSuite() {
	fmt.Println("SetUpSuite: UserTestSuite")
}

func (s *UserServiceSuite) TearDown() {
	fmt.Println("TearDownSuite: UserServiceSuite")
}

func (s *UserServiceSuite) SetupTest() {
	fmt.Println("SetUpTest: UserServiceSuite")
	db, mockDB, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	s.NoError(err)
	master := sqlx.NewDb(db, "sqlmock")

	s.mocks = listMock{
		database: mockDB,
		sqlmock:  mockDB,
		userRepo: *mocksUserRepo.NewIUserRepository(s.T()),
	}
	s.repositories = repository.NewRegistry(
		master,
		&s.mocks.userRepo,
	)
	s.UserService = NewUserService(common.NewRegistry(), s.repositories)
}

func (s *UserServiceSuite) TearDownTest() {
	fmt.Println("TearDownTest: UserServiceSuite")
}

func (s *UserServiceSuite) TestRegister() {
	var wg sync.WaitGroup
	type (
		args struct {
			ctx     context.Context
			payload model.RegistrationUserRequest
		}
		want struct {
			err error
		}

		testCase struct {
			name     string
			args     args
			mockFunc func(listMock *listMock, args args)
			want     want
		}
	)
	testCases := []testCase{
		{
			name: "register failed, username duplicate",
			args: args{
				ctx: context.Background(),
				payload: model.RegistrationUserRequest{
					Username:  "test123",
					Password:  "test123",
					FirstName: "test",
					LastName:  new(string),
					Email:     "test123@gmail.com",
				},
			},
			mockFunc: func(m *listMock, args args) {
				// membutuhkan mock ke kedua method
				// karena implementasi goroutine menyebabkan kedua method berjalan parallel

				wg.Add(2)
				go func() {
					defer wg.Done()
					s.mocks.userRepo.On("FindByUsername", mock.Anything, args.payload.Username).
						Return(&model.User{Username: "test123", Status: common.StatusUserActive}, nil)
				}()

				go func() {
					defer wg.Done()
					s.mocks.userRepo.On("FindByEmail", mock.Anything, args.payload.Email).
						Return(nil, nil)
				}()
			},
			want: want{
				err: common.ErrUsernameAlreadyTaken,
			},
		},
		{
			name: "register failed, error when FindByUsername",
			args: args{
				ctx: context.Background(),
				payload: model.RegistrationUserRequest{
					Username:  "test123",
					Password:  "test123",
					FirstName: "test",
					LastName:  new(string),
					Email:     "test123@gmail.com",
				},
			},
			mockFunc: func(m *listMock, args args) {
				// membutuhkan mock ke kedua method
				// karena implementasi goroutine menyebabkan kedua method berjalan parallel
				wg.Add(2)
				go func() {
					defer wg.Done()
					s.mocks.userRepo.On("FindByUsername", mock.Anything, args.payload.Username).
						Return(nil, common.ErrSQLExec)
				}()

				go func() {
					defer wg.Done()
					s.mocks.userRepo.On("FindByEmail", mock.Anything, args.payload.Email).
						Return(nil, nil)
				}()
			},
			want: want{
				err: common.ErrSQLExec,
			},
		},
		{
			name: "register failed, email duplicate",
			args: args{
				ctx: context.Background(),
				payload: model.RegistrationUserRequest{
					Username:  "test123",
					Password:  "test123",
					FirstName: "test",
					LastName:  new(string),
					Email:     "test123@gmail.com",
				},
			},
			mockFunc: func(m *listMock, args args) {
				wg.Add(2)
				go func() {
					defer wg.Done()
					s.mocks.userRepo.On("FindByUsername", mock.Anything, args.payload.Username).
						Return(nil, nil)
				}()

				go func() {
					defer wg.Done()
					s.mocks.userRepo.On("FindByEmail", mock.Anything, args.payload.Email).
						Return(&model.User{
							Username: "test1234",
							Status:   common.StatusUserActive,
							Email:    "test123@gmail.com",
						}, nil)
				}()
			},
			want: want{
				err: common.ErrEmailAlreadyTaken,
			},
		},
		{
			name: "register failed, error when FindByEmail",
			args: args{
				ctx: context.Background(),
				payload: model.RegistrationUserRequest{
					Username:  "test123",
					Password:  "test123",
					FirstName: "test",
					LastName:  new(string),
					Email:     "test123@gmail.com",
				},
			},
			mockFunc: func(m *listMock, args args) {
				wg.Add(2)
				go func() {
					defer wg.Done()
					s.mocks.userRepo.On("FindByUsername", mock.Anything, args.payload.Username).
						Return(nil, nil)
				}()
				go func() {
					defer wg.Done()
					s.mocks.userRepo.On("FindByEmail", mock.Anything, args.payload.Email).
						Return(nil, common.ErrSQLExec)
				}()
			},
			want: want{
				err: common.ErrSQLExec,
			},
		},
		{
			name: "register failed, error when Insert",
			args: args{
				ctx: context.Background(),
				payload: model.RegistrationUserRequest{
					Username:  "test123",
					Password:  "test123",
					FirstName: "test",
					LastName:  new(string),
					Email:     "test123@gmail.com",
				},
			},
			mockFunc: func(m *listMock, args args) {
				wg.Add(2)
				go func() {
					defer wg.Done()
					s.mocks.userRepo.On("FindByUsername", mock.Anything, args.payload.Username).
						Return(nil, nil)
				}()
				go func() {
					defer wg.Done()
					s.mocks.userRepo.On("FindByEmail", mock.Anything, args.payload.Email).
						Return(nil, nil)
				}()
				s.mocks.userRepo.On("Insert", mock.Anything, mock.Anything).
					Return(nil, common.ErrSQLExec)
			},
			want: want{
				err: common.ErrSQLExec,
			},
		},
		{
			name: "register success, return nil",
			args: args{
				ctx: context.Background(),
				payload: model.RegistrationUserRequest{
					Username:  "test123",
					Password:  "test123",
					FirstName: "test",
					LastName:  new(string),
					Email:     "test123@gmail.com",
				},
			},
			mockFunc: func(m *listMock, args args) {
				wg.Add(2)
				go func() {
					defer wg.Done()
					s.mocks.userRepo.On("FindByUsername", mock.Anything, args.payload.Username).
						Return(nil, nil)
				}()
				go func() {
					defer wg.Done()
					s.mocks.userRepo.On("FindByEmail", mock.Anything, args.payload.Email).
						Return(nil, nil)
				}()
				s.mocks.userRepo.On("Insert", mock.Anything, mock.Anything).
					Return(nil, nil)
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockFunc(&s.mocks, tc.args)
			wg.Wait()

			response := s.UserService.Register(tc.args.ctx, tc.args.payload)
			s.Equal(tc.want.err, response)
		})
	}
}

func (s *UserServiceSuite) TestDetail() {
	type (
		args struct {
			ctx    context.Context
			userId int64
		}
		want struct {
			*model.User
			err error
		}

		testCase struct {
			name     string
			args     args
			mockFunc func(listMock *listMock, args args)
			want     want
		}
	)
	testCases := []testCase{
		{
			name: "get Detail failed, data not found",
			args: args{
				ctx:    context.Background(),
				userId: 1,
			},
			mockFunc: func(m *listMock, args args) {
				s.mocks.userRepo.On("FindById", mock.Anything, args.userId).
					Return(nil, common.ErrDataNotFound)
			},
			want: want{
				err: common.ErrDataNotFound,
			},
		},
		{
			name: "get Detail failed, error when FindById",
			args: args{
				ctx:    context.Background(),
				userId: 1,
			},
			mockFunc: func(m *listMock, args args) {
				s.mocks.userRepo.On("FindById", mock.Anything, args.userId).
					Return(nil, common.ErrSQLExec)
			},
			want: want{
				err: common.ErrSQLExec,
			},
		},
		{
			name: "get Detail successfully",
			args: args{
				ctx:    context.Background(),
				userId: 1,
			},
			mockFunc: func(m *listMock, args args) {
				s.mocks.userRepo.On("FindById", mock.Anything, args.userId).
					Return(&model.User{
						ID:       1,
						Username: "test123",
						Status:   common.StatusUserActive,
					}, nil)
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockFunc(&s.mocks, tc.args)

			_, err := s.UserService.Detail(tc.args.ctx, tc.args.userId)

			s.Equal(tc.want.err, err)
		})
	}
}

func (s *UserServiceSuite) TestDelete() {
	type (
		args struct {
			ctx    context.Context
			userId int64
		}
		want struct {
			err error
		}

		testCase struct {
			name     string
			args     args
			mockFunc func(listMock *listMock, args args)
			want     want
		}
	)
	testCases := []testCase{
		{
			name: "Delete user failed, data not found",
			args: args{
				ctx:    context.Background(),
				userId: 1,
			},
			mockFunc: func(m *listMock, args args) {
				s.mocks.userRepo.On("FindById", mock.Anything, args.userId).
					Return(nil, common.ErrDataNotFound)
			},
			want: want{
				err: common.ErrDataNotFound,
			},
		},
		{
			name: "delete user failed, error when FindById",
			args: args{
				ctx:    context.Background(),
				userId: 1,
			},
			mockFunc: func(m *listMock, args args) {
				s.mocks.userRepo.On("FindById", mock.Anything, args.userId).
					Return(nil, common.ErrSQLExec)
			},
			want: want{
				err: common.ErrSQLExec,
			},
		},
		{
			name: "delete user failed, user already deleted",
			args: args{
				ctx:    context.Background(),
				userId: 1,
			},
			mockFunc: func(m *listMock, args args) {
				s.mocks.userRepo.On("FindById", mock.Anything, args.userId).
					Return(&model.User{ID: 1, Status: common.StatusUserInactive}, nil)
			},
			want: want{
				err: common.ErrUserAlreadyDeleted,
			},
		},
		{
			name: "Delete user failed, error when Update",
			args: args{
				ctx:    context.Background(),
				userId: 1,
			},
			mockFunc: func(m *listMock, args args) {
				s.mocks.userRepo.On("FindById", mock.Anything, args.userId).
					Return(&model.User{
						ID:       1,
						Username: "test123",
						Status:   common.StatusUserActive,
					}, nil)
				s.mocks.userRepo.On("Update", mock.Anything, mock.Anything).
					Return(common.ErrSQLExec)
			},
			want: want{
				err: common.ErrSQLExec,
			},
		},
		{
			name: "Delete user successfully",
			args: args{
				ctx:    context.Background(),
				userId: 1,
			},
			mockFunc: func(m *listMock, args args) {
				s.mocks.userRepo.On("FindById", mock.Anything, args.userId).
					Return(&model.User{
						ID:       1,
						Username: "test123",
						Status:   common.StatusUserActive,
					}, nil)
				s.mocks.userRepo.On("Update", mock.Anything, mock.Anything).
					Return(nil)
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockFunc(&s.mocks, tc.args)

			err := s.UserService.Delete(tc.args.ctx, tc.args.userId)

			s.Equal(tc.want.err, err)
		})
	}
}

func (s *UserServiceSuite) TestUpdate() {
	type (
		args struct {
			ctx     context.Context
			payload model.UpdateUserRequest
		}
		want struct {
			err error
		}

		testCase struct {
			name     string
			args     args
			mockFunc func(listMock *listMock, args args)
			want     want
		}
	)
	testCases := []testCase{
		{
			name: "Update user failed, data not found",
			args: args{
				ctx: context.Background(),
				payload: model.UpdateUserRequest{
					Id:        1,
					FirstName: "test",
					LastName:  new(string),
				},
			},
			mockFunc: func(m *listMock, args args) {
				s.mocks.userRepo.On("FindById", mock.Anything, args.payload.Id).
					Return(nil, common.ErrDataNotFound)
			},
			want: want{
				err: common.ErrDataNotFound,
			},
		},
		{
			name: "Update user failed, error when FindById",
			args: args{
				ctx: context.Background(),
				payload: model.UpdateUserRequest{
					Id:        1,
					FirstName: "test",
					LastName:  new(string),
				},
			},
			mockFunc: func(m *listMock, args args) {
				s.mocks.userRepo.On("FindById", mock.Anything, args.payload.Id).
					Return(nil, common.ErrSQLExec)
			},
			want: want{
				err: common.ErrSQLExec,
			},
		},
		{
			name: "Update user failed, error when Update",
			args: args{
				ctx: context.Background(),
				payload: model.UpdateUserRequest{
					Id:        1,
					FirstName: "test",
					LastName:  new(string),
				},
			},
			mockFunc: func(m *listMock, args args) {
				s.mocks.userRepo.On("FindById", mock.Anything, args.payload.Id).
					Return(&model.User{
						ID:       1,
						Username: "test123",
					}, nil)
				s.mocks.userRepo.On("Update", mock.Anything, mock.Anything).
					Return(common.ErrSQLExec)
			},
			want: want{
				err: common.ErrSQLExec,
			},
		},
		{
			name: "Update user successfully",
			args: args{
				ctx: context.Background(),
				payload: model.UpdateUserRequest{
					Id:        1,
					FirstName: "test",
					LastName:  new(string),
				},
			},
			mockFunc: func(m *listMock, args args) {
				s.mocks.userRepo.On("FindById", mock.Anything, args.payload.Id).
					Return(&model.User{
						ID:       1,
						Username: "test123",
					}, nil)
				s.mocks.userRepo.On("Update", mock.Anything, mock.Anything).
					Return(nil)
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockFunc(&s.mocks, tc.args)

			err := s.UserService.Update(tc.args.ctx, tc.args.payload)

			s.Equal(tc.want.err, err)
		})
	}
}

func (s *UserServiceSuite) TestLogin() {
	type (
		args struct {
			ctx     context.Context
			payload model.LoginRequest
		}
		want struct {
			err error
		}

		testCase struct {
			name     string
			args     args
			mockFunc func(listMock *listMock, args args)
			want     want
		}
	)
	testCases := []testCase{
		{
			name: "Login failed, data not found",
			args: args{
				ctx: context.Background(),
				payload: model.LoginRequest{
					Username: "test",
					Password: "test",
				},
			},
			mockFunc: func(m *listMock, args args) {
				s.mocks.userRepo.On("FindByUsername", mock.Anything, args.payload.Username).
					Return(nil, common.ErrDataNotFound)
			},
			want: want{
				err: common.ErrDataNotFound,
			},
		},
		{
			name: "Login failed, error when FindByUsername",
			args: args{
				ctx: context.Background(),
				payload: model.LoginRequest{
					Username: "test",
					Password: "test",
				},
			},
			mockFunc: func(m *listMock, args args) {
				s.mocks.userRepo.On("FindByUsername", mock.Anything, args.payload.Username).
					Return(nil, common.ErrSQLExec)
			},
			want: want{
				err: common.ErrSQLExec,
			},
		},
		{
			name: "Login failed, user already deleted",
			args: args{
				ctx: context.Background(),
				payload: model.LoginRequest{
					Username: "test",
					Password: "test",
				},
			},
			mockFunc: func(m *listMock, args args) {
				s.mocks.userRepo.On("FindByUsername", mock.Anything, args.payload.Username).
					Return(&model.User{
						ID:     1,
						Status: common.StatusUserInactive,
					}, nil)
			},
			want: want{
				err: common.ErrUserAlreadyDeleted,
			},
		},
		{
			name: "Login successfully",
			args: args{
				ctx: context.Background(),
				payload: model.LoginRequest{
					Username: "test",
					Password: "test",
				},
			},
			mockFunc: func(m *listMock, args args) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(args.payload.Password), bcrypt.DefaultCost)

				s.mocks.userRepo.On("FindByUsername", mock.Anything, args.payload.Username).
					Return(&model.User{
						ID:       1,
						Password: string(hashedPassword),
						Status:   common.StatusUserActive,
					}, nil)

			},
			want: want{
				err: nil,
			},
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockFunc(&s.mocks, tc.args)

			_, _, err := s.UserService.Login(tc.args.ctx, tc.args.payload)

			s.Equal(tc.want.err, err)
		})
	}
}

func (s *UserServiceSuite) TestAuthenticate() {
	type (
		args struct {
			ctx   context.Context
			token string
		}
		want struct {
			err error
		}

		testCase struct {
			name     string
			args     args
			mockFunc func(listMock *listMock, args args)
			want     want
		}
	)

	testCases := []testCase{
		{
			name: "Authenticate failed, missing token",
			args: args{
				ctx:   context.Background(),
				token: "",
			},
			mockFunc: func(m *listMock, args args) {
			},
			want: want{
				err: common.ErrInvalidToken,
			},
		},
		{
			name: "Authenticate failed, invalid token format",
			args: args{
				ctx:   context.Background(),
				token: "invalid_token_format",
			},
			mockFunc: func(m *listMock, args args) {
			},
			want: want{
				err: common.ErrInvalidToken,
			},
		},
		{
			name: "Authenticate failed, invalid signing method",
			args: args{
				ctx:   context.Background(),
				token: "Bearer <invalid-token>",
			},
			mockFunc: func(m *listMock, args args) {
			},
			want: want{
				err: common.ErrInvalidToken,
			},
		},
		{
			name: "Authenticate failed, expired token",
			args: args{
				ctx:   context.Background(),
				token: generateToken("testUser", jwt.NumericDate{Time: time.Now().Add(-time.Hour)}), // Expired token (1 hour ago)
			},
			mockFunc: func(m *listMock, args args) {

			},
			want: want{
				err: common.ErrInvalidToken,
			},
		},
		{
			name: "Authenticate failed, empty username in token",
			args: args{
				ctx:   context.Background(),
				token: generateToken("", jwt.NumericDate{Time: time.Now().Add(time.Hour * 24)}), // Empty username

			},
			mockFunc: func(m *listMock, args args) {

			},
			want: want{
				err: common.ErrInvalidToken,
			},
		},
		{
			name: "Authenticate success",
			args: args{
				ctx:   context.Background(),
				token: generateToken("userTest", jwt.NumericDate{Time: time.Now().Add(time.Hour * 24)}), // Empty username
			},
			mockFunc: func(m *listMock, args args) {

			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.mockFunc(&s.mocks, tc.args)
			err := s.UserService.Authenticate(tc.args.ctx, tc.args.token)
			s.Equal(tc.want.err, err)
		})
	}
}

func generateToken(username string, time jwt.NumericDate) string {
	claims := &model.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &time,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validToken, _ := token.SignedString([]byte(config.Cold.SecretKey))
	return validToken
}
