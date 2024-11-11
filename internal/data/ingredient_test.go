package data

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"

	// "fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/require"

	"github.com/jsmit257/huautla/types"
)

var (
	_ingredients = []ingredient{
		{UUID: "0", Name: "ingredient 0"},
		{UUID: "1", Name: "ingredient 1"},
		{UUID: "2", Name: "ingredient 2"},
	}
	ingFields = row{"id", "name"}
	ingValues = [][]driver.Value{
		{_ingredients[0].UUID, _ingredients[0].Name},
		{_ingredients[1].UUID, _ingredients[1].Name},
		{_ingredients[2].UUID, _ingredients[2].Name},
	}
)

func Test_SelectAllIngredients(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "SelectAllIngredients")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result []types.Ingredient
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				ingFields.mock(mock, ingValues...)
				return db
			},
			result: []types.Ingredient{
				types.Ingredient(_ingredients[0]),
				types.Ingredient(_ingredients[1]),
				types.Ingredient(_ingredients[2]),
			},
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))
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
				query:        tc.db(sqlmock.New()),
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

	l := log.WithField("test", "SelectIngredient")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result types.Ingredient
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"name"}).
						AddRow("ingredient 0"))
				return db
			},
			id:     "0",
			result: types.Ingredient{UUID: "0", Name: "ingredient 0"},
		},
		"no_rows_returned": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				ingFields.mock(mock)
				return db
			},
			err: fmt.Errorf("sql: no rows in result set"),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				ingFields.fail()(&mocker{mock}) // kinda rude?
				return db
			},
			err: ingFields.err(),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := (&Conn{
				query:        tc.db(sqlmock.New()),
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

	l := log.WithField("test", "InsertIngredient")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result types.Ingredient
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id:     "0",
			result: types.Ingredient{UUID: "30313233-3435-3637-3839-616263646566", Name: "ingredient 0"},
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:     "0",
			result: types.Ingredient{UUID: "30313233-3435-3637-3839-616263646566", Name: "ingredient 0"},
			err:    fmt.Errorf("ingredient was not added"),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			id:     "0",
			result: types.Ingredient{UUID: "30313233-3435-3637-3839-616263646566", Name: "ingredient 0"},
			err:    fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			id:     "0",
			result: types.Ingredient{UUID: "30313233-3435-3637-3839-616263646566", Name: "ingredient 0"},
			err:    fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).InsertIngredient(
				context.Background(),
				types.Ingredient{UUID: tc.id, Name: "ingredient " + string(tc.id)},
				"Test_InsertIngredients")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_UpdateIngredients(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "UpdateIngredient")

	tcs := map[string]struct {
		db  getMockDB
		id  types.UUID
		err error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id: "0",
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:  "0",
			err: fmt.Errorf("ingredient was not updated: '0'"),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
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
				query:        tc.db(sqlmock.New()),
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

	l := log.WithField("test", "DeleteIngredient")

	tcs := map[string]struct {
		db  getMockDB
		id  types.UUID
		err error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id: "0",
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:  "0",
			err: fmt.Errorf("ingredient could not be deleted: '0'"),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").
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
				query:        tc.db(sqlmock.New()),
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
