package postgres

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"app/internal/core/entity"
	domainErrors "app/internal/core/error"
	"app/pkg/transactor"
	pgxTransactor "app/pkg/transactor/pgx"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestMock(t *testing.T) (transactor.Transactor, pgxTransactor.DBGetter, pgxmock.PgxPoolIface) {
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)

	txManager, getter := pgxTransactor.New(mockPool)

	return txManager, getter, mockPool
}

func TestUserRepository_Create(t *testing.T) {
	t.Parallel()

	mockUser := &entity.User{
		ID:       uuid.New(),
		Username: "testuser",
	}
	mockTime := time.Now()

	duplicateKeyErr := &pgconn.PgError{
		Code:    duplicateKeyErrorCode,
		Message: "duplicate key",
	}
	genericErr := errors.New("something went wrong")

	sql, _, err := psql.
		Insert("users").
		Columns("id", "username").
		Values(mockUser.ID, mockUser.Username).
		Suffix("RETURNING created_at").
		ToSql()
	require.NoError(t, err)

	testCases := []struct {
		name        string
		inputUser   *entity.User
		setupMock   func(mock pgxmock.PgxPoolIface)
		expectedErr error
	}{
		{
			name:      "Success",
			inputUser: mockUser,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"created_at"}).AddRow(mockTime)
				mock.ExpectQuery(regexp.QuoteMeta(sql)).
					WithArgs(mockUser.ID, mockUser.Username).
					WillReturnRows(rows)
			},
			expectedErr: nil,
		},
		{
			name:      "Duplicate Key Error",
			inputUser: mockUser,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(sql)).
					WithArgs(mockUser.ID, mockUser.Username).
					WillReturnError(duplicateKeyErr)
			},
			expectedErr: domainErrors.ErrUserAlreadyExists,
		},
		{
			name:      "Generic DB Error",
			inputUser: mockUser,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(sql)).
					WithArgs(mockUser.ID, mockUser.Username).
					WillReturnError(genericErr)
			},
			expectedErr: genericErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, dbGetter, mockPool := newTestMock(t)
			repo := NewUserRepository(dbGetter)

			tc.setupMock(mockPool)

			err := repo.Create(context.Background(), tc.inputUser)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr, "expected error: %v, got: %v", tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, mockTime, tc.inputUser.CreatedAt)
			}

			assert.NoError(t, mockPool.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	t.Parallel()

	mockID := uuid.New()
	mockUser := &entity.User{
		ID:        mockID,
		Username:  "testuser",
		CreatedAt: time.Now(),
	}

	genericErr := errors.New("something went wrong")
	scanErr := fmt.Errorf("scan user:")

	sql, args, err := psql.
		Select("id", "username", "created_at").
		From("users").
		Where(sq.Eq{"id": mockID}).
		ToSql()
	require.NoError(t, err)

	testCases := []struct {
		name         string
		inputID      uuid.UUID
		setupMock    func(mock pgxmock.PgxPoolIface)
		expectedUser *entity.User
		expectedErr  error
	}{
		{
			name:    "Success",
			inputID: mockID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "username", "created_at"}).
					AddRow(mockUser.ID, mockUser.Username, mockUser.CreatedAt)
				mock.ExpectQuery(regexp.QuoteMeta(sql)).
					WithArgs(args...).
					WillReturnRows(rows)
			},
			expectedUser: mockUser,
			expectedErr:  nil,
		},
		{
			name:    "Not Found",
			inputID: mockID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(sql)).
					WithArgs(args...).
					WillReturnError(pgx.ErrNoRows)
			},
			expectedUser: nil,
			expectedErr:  domainErrors.ErrUserNotFound,
		},
		{
			name:    "Scan Error",
			inputID: mockID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(mockUser.ID)
				mock.ExpectQuery(regexp.QuoteMeta(sql)).
					WithArgs(args...).
					WillReturnRows(rows)
			},
			expectedUser: nil,
			expectedErr:  scanErr,
		},
		{
			name:    "Generic DB Error",
			inputID: mockID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(sql)).
					WithArgs(args...).
					WillReturnError(genericErr)
			},
			expectedUser: nil,
			expectedErr:  genericErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, dbGetter, mockPool := newTestMock(t)
			repo := NewUserRepository(dbGetter)

			tc.setupMock(mockPool)

			user, err := repo.GetByID(context.Background(), tc.inputID)

			if tc.expectedErr != nil {
				assert.Error(t, err)

				if errors.Is(tc.expectedErr, domainErrors.ErrUserNotFound) {
					assert.ErrorIs(t, err, tc.expectedErr)
				} else {
					assert.ErrorContains(t, err, tc.expectedErr.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedUser, user)
			assert.NoError(t, mockPool.ExpectationsWereMet())
		})
	}
}
