package postgres

import (
	"context"
	"errors"
	"fmt"

	"app/internal/domain/entity"
	domainErrors "app/internal/domain/error"
	pgxTransactor "app/pkg/transactor/pgx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepository struct {
	dbGetter pgxTransactor.DBGetter
}

func NewUserRepository(dbGetter pgxTransactor.DBGetter) *UserRepository {
	return &UserRepository{dbGetter: dbGetter}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	sql, args, err := psql.
		Insert("users").
		Columns("id", "username").
		Values(user.ID, user.Username).
		Suffix("RETURNING created_at").
		ToSql()

	if err != nil {
		return fmt.Errorf("make query: %w", err)
	}

	err = r.dbGetter(ctx).QueryRow(ctx, sql, args...).Scan(&user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == duplicateKeyErrorCode {
				return domainErrors.ErrUserAlreadyExists
			}
		}

		return fmt.Errorf("execute query: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	sql, args, err := psql.
		Select("id", "username", "created_at").
		From("users").
		Where("id", id).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("make query: %w", err)
	}

	row := r.dbGetter(ctx).QueryRow(ctx, sql, args...)
	user, err := r.scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}

	return user, nil
}

func (r *UserRepository) scanUser(row pgx.Row) (*entity.User, error) {
	user := &entity.User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErrors.ErrUserNotFound
		}

		return nil, fmt.Errorf("scan user: %w", err)
	}

	return user, nil
}
