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

func Test_SelectAllSubstrates(t *testing.T) {
	t.Parallel()

	querypat, l := sqls["select-all"],
		log.WithField("test", "SelectAllSubstrates")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result []types.Substrate
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery(querypat).
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name", "type", "vendor_uuid", "vendor_name"}).
						AddRow("0", "substrate 0", types.GrainType, "0", "vendor 0").
						AddRow("1", "substrate 1", types.GrainType, "1", "vendor 1").
						AddRow("2", "substrate 2", types.GrainType, "1", "vendor 1"))
				return db
			},
			result: []types.Substrate{
				types.Substrate{"0", "substrate 0", types.GrainType, types.Vendor{"0", "vendor 0"}, nil},
				types.Substrate{"1", "substrate 1", types.GrainType, types.Vendor{"1", "vendor 1"}, nil},
				types.Substrate{"2", "substrate 2", types.GrainType, types.Vendor{"1", "vendor 1"}, nil},
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
			}).SelectAllSubstrates(context.Background(), "Test_SelectAllSubstrates")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_SelectSubstrate(t *testing.T) {
	t.Parallel()

	var querypat = sqls["select"]

	l := log.WithField("test", "SelectSubstrate")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result types.Substrate
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery(querypat).
					WillReturnRows(sqlmock.
						NewRows([]string{"name", "type", "vendor_uuid", "vendor_name"}).
						AddRow("substrate 0", types.GrainType, "0", "vendor 0"))
				return db
			},
			id:     "0",
			result: types.Substrate{"0", "substrate 0", types.GrainType, types.Vendor{"0", "vendor 0"}, nil},
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
			}).SelectSubstrate(context.Background(), tc.id, "Test_SelectSubstrate")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_InsertSubstrate(t *testing.T) {
	t.Parallel()

	var querypat = sqls["insert"]

	l := log.WithField("test", "InsertSubstrate")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		tp     types.SubstrateType
		result types.Substrate
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
			tp:     types.GrainType,
			result: types.Substrate{"30313233-3435-3637-3839-616263646566", "substrate 0", types.GrainType, types.Vendor{}, nil},
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
			tp:     types.GrainType,
			result: types.Substrate{"30313233-3435-3637-3839-616263646566", "substrate 0", types.GrainType, types.Vendor{}, nil},
			err:    fmt.Errorf("substrate was not added"),
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
			tp:     types.GrainType,
			result: types.Substrate{"30313233-3435-3637-3839-616263646566", "substrate 0", types.GrainType, types.Vendor{}, nil},
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
			tp:     types.GrainType,
			result: types.Substrate{"30313233-3435-3637-3839-616263646566", "substrate 0", types.GrainType, types.Vendor{}, nil},
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
			}).InsertSubstrate(
				context.Background(),
				types.Substrate{tc.id, "substrate " + string(tc.id), tc.tp, types.Vendor{}, nil},
				"Test_InsertSubstrate")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_UpdateSubstrate(t *testing.T) {
	t.Parallel()

	var querypat = sqls["update"]

	l := log.WithField("test", "UpdateSubstrate")

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
			err: fmt.Errorf("substrate was not updated: '0'"),
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
			}).UpdateSubstrate(
				context.Background(),
				tc.id,
				types.Substrate{Name: "substrate " + string(tc.id)},
				"Test_UpdateSubstrates")

			require.Equal(t, tc.err, err)
		})
	}
}

func foo(t *testing.T, querypat string) {

}

func Test_DeleteSubstrate(t *testing.T) {
	t.Parallel()

	var querypat = sqls["delete"]

	l := log.WithField("test", "DeleteSubstrate")

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
			err: fmt.Errorf("substrate could not be deleted: '0'"),
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
			}).DeleteSubstrate(
				context.Background(),
				tc.id,
				"Test_DeleteSubstrate")

			require.Equal(t, tc.err, err)
		})
	}
}

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
