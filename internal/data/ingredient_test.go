package data

import (
	"context"
	"database/sql"
	"fmt"

	// "fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/require"

	"github.com/jsmit257/huautla/types"
)

func Test_SelectAllIngredients(t *testing.T) {
	t.Parallel()

	querypat, l := sqls["select-all"],
		log.WithField("test", "SelectAllIngredients")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result []types.Ingredient
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery(querypat).
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name"}).
						AddRow("0", "ingredient 0").
						AddRow("1", "ingredient 1").
						AddRow("2", "ingredient 2"))
				return db
			},
			result: []types.Ingredient{
				types.Ingredient{"0", "ingredient 0"},
				types.Ingredient{"1", "ingredient 1"},
				types.Ingredient{"2", "ingredient 2"},
			},
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectQuery(querypat).
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err: fmt.Errorf("some error"),
		},
		// "query_result_nil": {}, // FIXME: how to mock?
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).SelectAllIngredients(context.Background(), "Test_SelectAllIngredients")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_SelectIngredient(t *testing.T) {
	t.Parallel()

	var querypat = sqls["select"]

	l := log.WithField("test", "SelectIngredient")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result types.Ingredient
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery(querypat).
					WillReturnRows(sqlmock.
						NewRows([]string{"name"}).
						AddRow("ingredient 0"))
				return db
			},
			id:     "0",
			result: types.Ingredient{"0", "ingredient 0"},
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectQuery(querypat).
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err: fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).SelectIngredient(context.Background(), tc.id, "Test_SelectIngredients")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_InsertIngredients(t *testing.T) {
	t.Parallel()

	var querypat = sqls["insert"]

	l := log.WithField("test", "InsertIngredient")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result types.Ingredient
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id:     "0",
			result: types.Ingredient{"30313233-3435-3637-3839-616263646566", "ingredient 0"},
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:     "0",
			result: types.Ingredient{"30313233-3435-3637-3839-616263646566", "ingredient 0"},
			err:    fmt.Errorf("ingredient was not added"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			id:     "0",
			result: types.Ingredient{"30313233-3435-3637-3839-616263646566", "ingredient 0"},
			err:    fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			id:     "0",
			result: types.Ingredient{"30313233-3435-3637-3839-616263646566", "ingredient 0"},
			err:    fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).InsertIngredient(
				context.Background(),
				types.Ingredient{tc.id, "ingredient " + string(tc.id)},
				"Test_InsertIngredients")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_UpdateIngredients(t *testing.T) {
	t.Parallel()

	var querypat = sqls["update"]

	l := log.WithField("test", "UpdateIngredient")

	tcs := map[string]struct {
		db  getMockDB
		id  types.UUID
		err error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id: "0",
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:  "0",
			err: fmt.Errorf("ingredient was not updated: '0'"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).UpdateIngredient(
				context.Background(),
				tc.id,
				types.Ingredient{Name: "ingredient " + string(tc.id)},
				"Test_UpdateIngredients")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_DeleteIngredients(t *testing.T) {
	t.Parallel()

	var querypat = sqls["delete"]

	l := log.WithField("test", "DeleteIngredient")

	tcs := map[string]struct {
		db  getMockDB
		id  types.UUID
		err error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id: "0",
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:  "0",
			err: fmt.Errorf("ingredient could not be deleted: '0'"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).DeleteIngredient(
				context.Background(),
				tc.id,
				"Test_DeleteIngredients")

			require.Equal(t, tc.err, err)
		})
	}
}
