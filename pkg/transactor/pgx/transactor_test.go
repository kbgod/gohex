package pgx

import (
	"context"
	"errors"
	"testing"

	"app/pkg/transactor"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func InitTestMock(t *testing.T) (transactor.Transactor, DBGetter, pgxmock.PgxPoolIface) {
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)

	txManager, getter := New(mockPool)

	return txManager, getter, mockPool
}

type testRepo struct {
	dbGetter DBGetter
}

func (r testRepo) CreateTestTable(ctx context.Context) error {
	_, err := r.dbGetter(ctx).Exec(ctx, "CREATE TABLE tmp (id SERIAL PRIMARY KEY);")

	return err
}

func (r testRepo) DropTestTable(ctx context.Context) error {
	_, err := r.dbGetter(ctx).Exec(ctx, "DROP TABLE tmp;")

	return err
}

func (r testRepo) CheckTestTable(ctx context.Context) (bool, error) {
	var exists bool

	err := r.dbGetter(ctx).QueryRow(ctx,
		"SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'tmp');",
	).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func Test_PgxPool_AtomicityRollback(t *testing.T) {
	txManager, dbGetter, mock := InitTestMock(t)
	defer mock.Close()

	ctx := context.Background()
	tr := testRepo{dbGetter: dbGetter}
	errRollback := errors.New("rollback")

	mock.ExpectBeginTx(pgx.TxOptions{})
	mock.ExpectQuery("SELECT EXISTS").WillReturnRows(
		pgxmock.NewRows([]string{"exists"}).AddRow(false),
	)
	mock.ExpectExec("CREATE TABLE tmp").WillReturnResult(pgxmock.NewResult("CREATE", 1))
	mock.ExpectQuery("SELECT EXISTS").WillReturnRows(
		pgxmock.NewRows([]string{"exists"}).AddRow(true),
	)
	mock.ExpectRollback()

	err := txManager.Do(ctx, func(ctx context.Context) error {
		exists, err := tr.CheckTestTable(ctx)
		require.NoError(t, err)
		require.False(t, exists)

		err = tr.CreateTestTable(ctx)
		require.NoError(t, err)

		exists, err = tr.CheckTestTable(ctx)
		require.NoError(t, err)
		require.True(t, exists)

		return errRollback
	})

	require.ErrorIs(t, err, errRollback)
	require.NoError(t, mock.ExpectationsWereMet())
}

func Test_PgxPool_AtomicityCommit(t *testing.T) {
	txManager, dbGetter, mock := InitTestMock(t)
	defer mock.Close()

	ctx := context.Background()
	tr := testRepo{dbGetter: dbGetter}

	mock.ExpectBeginTx(pgx.TxOptions{})
	mock.ExpectQuery("SELECT EXISTS").WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))
	mock.ExpectExec("CREATE TABLE tmp").WillReturnResult(pgxmock.NewResult("CREATE", 1))
	mock.ExpectQuery("SELECT EXISTS").WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectCommit()

	err := txManager.Do(ctx, func(ctx context.Context) error {
		exists, err := tr.CheckTestTable(ctx)
		require.NoError(t, err)
		require.False(t, exists)

		err = tr.CreateTestTable(ctx)
		require.NoError(t, err)

		exists, err = tr.CheckTestTable(ctx)
		require.NoError(t, err)
		require.True(t, exists)

		return nil
	})

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func Test_PgxPool_AtomicityNested(t *testing.T) {
	txManager, dbGetter, mock := InitTestMock(t)
	defer mock.Close()

	ctx := context.Background()
	tr := testRepo{dbGetter: dbGetter}

	mock.ExpectBeginTx(pgx.TxOptions{})
	mock.ExpectQuery("SELECT EXISTS").WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))
	mock.ExpectExec("CREATE TABLE tmp").WillReturnResult(pgxmock.NewResult("CREATE", 1))
	mock.ExpectQuery("SELECT EXISTS").WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectBeginTx(pgx.TxOptions{})
	mock.ExpectExec("DROP TABLE tmp").WillReturnResult(pgxmock.NewResult("DROP", 1))
	mock.ExpectCommit()

	mock.ExpectQuery("SELECT EXISTS").WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))
	mock.ExpectCommit()

	err := txManager.Do(ctx, func(ctx context.Context) error {
		exists, err := tr.CheckTestTable(ctx)
		require.NoError(t, err)
		require.False(t, exists)

		err = tr.CreateTestTable(ctx)
		require.NoError(t, err)

		exists, err = tr.CheckTestTable(ctx)
		require.NoError(t, err)
		require.True(t, exists)

		err = txManager.Do(ctx, tr.DropTestTable)
		require.NoError(t, err)

		exists, err = tr.CheckTestTable(ctx)
		require.NoError(t, err)
		require.False(t, exists)

		return nil
	})

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func Test_PgxPool_BeginError(t *testing.T) {
	txManager, _, mock := InitTestMock(t)
	defer mock.Close()

	ctx := context.Background()
	expectedErr := errors.New("begin failed")

	mock.ExpectBeginTx(pgx.TxOptions{}).WillReturnError(expectedErr)

	err := txManager.Do(ctx, func(ctx context.Context) error {
		t.Fatal("transaction body should NOT be executed when Begin fails")
		return nil
	})

	require.Error(t, err)
	require.ErrorContains(t, err, "begin failed")
	require.NoError(t, mock.ExpectationsWereMet())
}

func Test_PgxPool_CommitError(t *testing.T) {
	txManager, dbGetter, mock := InitTestMock(t)
	defer mock.Close()

	ctx := context.Background()
	tr := testRepo{dbGetter: dbGetter}
	expectedErr := errors.New("commit failed")

	mock.ExpectBeginTx(pgx.TxOptions{})

	mock.ExpectQuery("SELECT EXISTS").WillReturnRows(
		pgxmock.NewRows([]string{"exists"}).AddRow(false),
	)
	mock.ExpectExec("CREATE TABLE tmp").WillReturnResult(pgxmock.NewResult("CREATE", 1))
	mock.ExpectQuery("SELECT EXISTS").WillReturnRows(
		pgxmock.NewRows([]string{"exists"}).AddRow(true),
	)

	mock.ExpectCommit().WillReturnError(expectedErr)

	err := txManager.Do(ctx, func(ctx context.Context) error {
		exists, err := tr.CheckTestTable(ctx)
		require.NoError(t, err)
		require.False(t, exists)

		err = tr.CreateTestTable(ctx)
		require.NoError(t, err)

		exists, err = tr.CheckTestTable(ctx)
		require.NoError(t, err)
		require.True(t, exists)

		return nil
	})

	require.Error(t, err)
	require.ErrorContains(t, err, "commit failed")
	require.ErrorContains(t, err, "failed to commit transaction")
	require.NoError(t, mock.ExpectationsWereMet())
}

func Test_NewTransactorFromPool_ReturnsPool_WhenNoTxInContext(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	txManager, dbGetter := New(mock)

	ctx := txManager.Skip(context.Background())

	db := dbGetter(ctx)
	pgxTx, ok := db.(pgxmock.PgxPoolIface)

	require.Equal(t, mock, pgxTx, "pgxTx must be pgxpool")
	require.Equal(t, true, ok, "ok must be true")
}

func Test_IsWithinTransaction(t *testing.T) {
	type testCase struct {
		name string
		ctx  context.Context
		want bool
	}

	tx := struct{}{}
	key := transactorKey{}

	tests := []testCase{
		{
			name: "no value in context",
			ctx:  context.Background(),
			want: false,
		},
		{
			name: "nil value under correct key",
			ctx:  context.WithValue(context.Background(), key, nil),
			want: false,
		},
		{
			name: "non-nil value under correct key",
			ctx:  context.WithValue(context.Background(), key, tx),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsWithinTransaction(tt.ctx)
			require.Equal(t, tt.want, got)
		})
	}
}
