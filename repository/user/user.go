package user

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"user-management-service/common"
	"user-management-service/common/logger"
	"user-management-service/model"
)

type IUserRepository interface {
	Insert(ctx context.Context, payload model.User) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindById(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, payload model.User) error
	GetAll(ctx context.Context) ([]model.User, error)
}

type userRepository struct {
	common   common.IRegistry
	database *sqlx.DB
}

func NewUserRepository(common common.IRegistry, database *sqlx.DB) *userRepository {
	return &userRepository{
		common:   common,
		database: database,
	}
}

func (r *userRepository) Insert(ctx context.Context, payload model.User) (*model.User, error) {
	insertQuery := "INSERT INTO pengguna(username, password, first_name, last_name, email, status, created_at, updated_at) values (?,?,?,?,?,?, now(), now())"

	stmtx, err := r.database.PreparexContext(ctx, insertQuery)
	if err != nil {
		logger.Error(ctx, err.Error(), err, logger.Tag{Key: "logCtx", Value: ctx})
		return nil, common.WrapWithErr(err, common.ErrSQLQueryBuilder)
	}

	_, err = stmtx.ExecContext(
		ctx,
		payload.Username,
		payload.Password,
		payload.FirstName,
		payload.LastName,
		payload.Email,
		payload.Status,
	)
	if err != nil {
		logger.Error(ctx, err.Error(), err, logger.Tag{Key: "logCtx", Value: ctx})
		return nil, common.WrapWithErr(err, common.ErrSQLExec)
	}

	// for retrieve inserted user
	var lastInsertId int64
	err = r.database.GetContext(ctx, &lastInsertId, "SELECT LAST_INSERT_ID()")
	if err != nil {
		return nil, common.WrapWithErr(err, common.ErrSQLQueryBuilder)
	}

	var user model.User
	selectQuery := "SELECT id, username, email, first_name, last_name, status, created_at, updated_at FROM pengguna WHERE id = ?"
	err = r.database.GetContext(
		ctx,
		&user,
		selectQuery,
		lastInsertId,
	)
	if err != nil {
		logger.Error(ctx, err.Error(), err, logger.Tag{Key: "logCtx", Value: ctx})
		return nil, common.WrapWithErr(err, common.ErrSQLExec)
	}

	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	selectQuery := "SELECT id, username, password, first_name, last_name, status, email, created_at FROM pengguna WHERE username = ?"

	stmtx, err := r.database.PreparexContext(ctx, selectQuery)
	if err != nil {
		logger.Error(ctx, err.Error(), err, logger.Tag{Key: "logCtx", Value: ctx})
		return nil, common.WrapWithErr(err, common.ErrSQLQueryBuilder)
	}

	var user model.User
	err = stmtx.
		QueryRowContext(ctx, username).
		Scan(&user.ID, &user.Username, &user.Password, &user.FirstName, &user.LastName, &user.Status, &user.Email, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrDataNotFound
		}
		logger.Error(ctx, err.Error(), err, logger.Tag{Key: "logCtx", Value: ctx})
		return nil, common.WrapWithErr(err, common.ErrSQLExec)
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	selectQuery := "SELECT id, username, first_name, last_name, email, status, created_at FROM pengguna WHERE email = ?"
	stmtx, err := r.database.PreparexContext(ctx, selectQuery)
	if err != nil {
		return nil, err
	}
	var user model.User
	err = stmtx.
		QueryRowContext(ctx, email).
		Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Status, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrDataNotFound
		}
		logger.Error(ctx, err.Error(), err, logger.Tag{Key: "logCtx", Value: ctx})
		return nil, common.WrapWithErr(err, common.ErrSQLExec)
	}

	return &user, nil
}

func (r *userRepository) FindById(ctx context.Context, id int64) (*model.User, error) {
	selectQuery := "SELECT id, username, first_name, last_name, email, status, created_at FROM pengguna WHERE id = ?"
	stmtx, err := r.database.PreparexContext(ctx, selectQuery)
	if err != nil {
		logger.Error(ctx, err.Error(), err, logger.Tag{Key: "logCtx", Value: ctx})
		return nil, common.WrapWithErr(err, common.ErrSQLQueryBuilder)
	}
	var user model.User
	err = stmtx.
		QueryRowContext(ctx, id).
		Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Status, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrDataNotFound
		}
		logger.Error(ctx, err.Error(), err, logger.Tag{Key: "logCtx", Value: ctx})
		return nil, common.WrapWithErr(err, common.ErrSQLExec)
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, payload model.User) error {

	updateQuery := "UPDATE pengguna SET first_name = ?, last_name = ?, status = ?, updated_at = now() WHERE id = ?"

	stmtx, err := r.database.PreparexContext(ctx, updateQuery)
	if err != nil {
		logger.Error(ctx, err.Error(), err, logger.Tag{Key: "logCtx", Value: ctx})
		return common.WrapWithErr(err, common.ErrSQLQueryBuilder)
	}

	_, err = stmtx.ExecContext(
		ctx,
		payload.FirstName,
		payload.LastName,
		payload.Status,
		payload.ID,
	)
	if err != nil {
		logger.Error(ctx, err.Error(), err, logger.Tag{Key: "logCtx", Value: ctx})
		return common.WrapWithErr(err, common.ErrSQLExec)
	}

	return nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]model.User, error) {
	selectQuery := "SELECT id, username, first_name, last_name, email, status FROM pengguna"

	var users []model.User
	err := r.database.SelectContext(ctx, &users, selectQuery)
	if err != nil {
		logger.Error(ctx, err.Error(), err, logger.Tag{Key: "logCtx", Value: ctx})
		return nil, common.WrapWithErr(err, common.ErrSQLExec)
	}

	return users, nil
}
