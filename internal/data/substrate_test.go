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

	l := log.WithField("test", "SelectAllSubstrates")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result []types.Substrate
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name", "type", "vendor_uuid", "vendor_name"}).
						AddRow("0", "substrate 0", types.GrainType, "0", "vendor 0").
						AddRow("1", "substrate 1", types.GrainType, "1", "vendor 1").
						AddRow("2", "substrate 2", types.GrainType, "1", "vendor 1"))
				return db
			},
			result: []types.Substrate{
				{UUID: "0", Name: "substrate 0", Type: types.GrainType, Vendor: types.Vendor{UUID: "0", Name: "vendor 0"}, Ingredients: nil},
				{UUID: "1", Name: "substrate 1", Type: types.GrainType, Vendor: types.Vendor{UUID: "1", Name: "vendor 1"}, Ingredients: nil},
				{UUID: "2", Name: "substrate 2", Type: types.GrainType, Vendor: types.Vendor{UUID: "1", Name: "vendor 1"}, Ingredients: nil},
			},
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectQuery("").
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
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"name", "type", "vendor_uuid", "vendor_name"}).
						AddRow("substrate 0", types.GrainType, "0", "vendor 0"))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name"}).
						AddRow("0", "ingredient 0").
						AddRow("1", "ingredient 1").
						AddRow("2", "ingredient 2"))
				return db
			},
			id: "0",
			result: types.Substrate{
				UUID:   "0",
				Name:   "substrate 0",
				Type:   types.GrainType,
				Vendor: types.Vendor{UUID: "0", Name: "vendor 0"},
				Ingredients: []types.Ingredient{
					{UUID: "0", Name: "ingredient 0"},
					{UUID: "1", Name: "ingredient 1"},
					{UUID: "2", Name: "ingredient 2"},
				}},
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectQuery("").
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
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id:     "0",
			tp:     types.GrainType,
			result: types.Substrate{UUID: "30313233-3435-3637-3839-616263646566", Name: "substrate 0", Type: types.GrainType, Vendor: types.Vendor{}, Ingredients: nil},
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:     "0",
			tp:     types.GrainType,
			result: types.Substrate{UUID: "30313233-3435-3637-3839-616263646566", Name: "substrate 0", Type: types.GrainType, Vendor: types.Vendor{}, Ingredients: nil},
			err:    fmt.Errorf("substrate was not added"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			id:     "0",
			tp:     types.GrainType,
			result: types.Substrate{UUID: "30313233-3435-3637-3839-616263646566", Name: "substrate 0", Type: types.GrainType, Vendor: types.Vendor{}, Ingredients: nil},
			err:    fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			id:     "0",
			tp:     types.GrainType,
			result: types.Substrate{UUID: "30313233-3435-3637-3839-616263646566", Name: "substrate 0", Type: types.GrainType, Vendor: types.Vendor{}, Ingredients: nil},
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
				types.Substrate{UUID: tc.id, Name: "substrate " + string(tc.id), Type: tc.tp, Vendor: types.Vendor{}, Ingredients: nil},
				"Test_InsertSubstrate")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_UpdateSubstrate(t *testing.T) {
	t.Parallel()

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
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id: "0",
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
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
					ExpectExec("").
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
					ExpectExec("").
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

func Test_DeleteSubstrate(t *testing.T) {
	t.Parallel()

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
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id: "0",
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
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
					ExpectExec("").
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
					ExpectExec("").
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
