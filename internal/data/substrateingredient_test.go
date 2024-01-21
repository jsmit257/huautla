package data

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jsmit257/huautla/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func Test_GetAllIngredients(t *testing.T) {
	t.Parallel()

	querypat, l := sqls["all-ingredients"],
		log.WithField("test", "GetAllIngredients")

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
			result: []types.Ingredient{},
			err:    fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := &types.Substrate{UUID: tc.id}

			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).GetAllIngredients(context.Background(), s, "Test_GetAllIngredients")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, s.Ingredients)
		})
	}
}

func Test_AddIngredient(t *testing.T) {
	t.Parallel()

	var querypat = sqls["add-ingredient"]

	l := log.WithField("test", "AddIngredient")

	vermiculite := types.Ingredient{"0", "Vermiculite"}

	tcs := map[string]struct {
		db     getMockDB
		result []types.Ingredient
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
			result: []types.Ingredient{vermiculite},
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("substrateingredient was not added"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
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
			err: fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := &types.Substrate{}
			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).AddIngredient(
				context.Background(),
				s,
				vermiculite,
				"Test_AddIngredient")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, s.Ingredients)
		})
	}
}

func Test_ChangeIngredient(t *testing.T) {
	t.Parallel()

	var querypat = sqls["change-ingredient"]

	l := log.WithField("test", "ChangeIngredient")

	vermiculite, millet, popcorn :=
		types.Ingredient{"0", "Vermiculite"},
		types.Ingredient{"1", "Millet"},
		types.Ingredient{"2", "Popcorn"}

	tcs := map[string]struct {
		db getMockDB
		before,
		after []types.Ingredient
		from,
		to types.Ingredient
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
			before: []types.Ingredient{vermiculite, popcorn},
			from:   popcorn,
			to:     millet,
			after:  []types.Ingredient{vermiculite, millet},
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("substrateingredient was not changed"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
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
			err: fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := &types.Substrate{Ingredients: tc.before}

			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).ChangeIngredient(
				context.Background(),
				s,
				tc.from,
				tc.to,
				"Test_ChangeIngredient")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.after, s.Ingredients)
		})
	}
}

func Test_RemoveIngredient(t *testing.T) {
	t.Parallel()

	var querypat = sqls["remove-ingredient"]

	l := log.WithField("test", "RemoveIngredient")

	vermiculite, millet, popcorn :=
		types.Ingredient{"0", "Vermiculite"},
		types.Ingredient{"1", "Millet"},
		types.Ingredient{"2", "Popcorn"}

	tcs := map[string]struct {
		db     getMockDB
		i      []types.Ingredient
		result []types.Ingredient
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
			i:      []types.Ingredient{millet, popcorn, vermiculite},
			result: []types.Ingredient{millet, popcorn},
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("substrateingredient was not removed"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
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
			err: fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := &types.Substrate{Ingredients: tc.i}

			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).RemoveIngredient(
				context.Background(),
				s,
				vermiculite,
				"Test_RemoveIngredient")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, s.Ingredients)
		})
	}
}
