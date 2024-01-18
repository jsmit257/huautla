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

var sqls = readSQL("pgsql.yaml")

func Test_SelectAllVendors(t *testing.T) {
	t.Parallel()

	querypat, l := sqls["select-all-vendors"],
		log.WithField("test", "SelectVendor")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result []types.Vendor
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery(querypat).
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name"}).
						AddRow("0", "vendor 0").
						AddRow("1", "vendor 1").
						AddRow("2", "vendor 2"))
				return db
			},
			result: []types.Vendor{
				types.Vendor{"0", "vendor 0"},
				types.Vendor{"1", "vendor 1"},
				types.Vendor{"2", "vendor 2"},
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
			}).SelectAllVendors(context.Background(), "Test_SelectAllVendors")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_SelectVendor(t *testing.T) {
	t.Parallel()

	var querypat = sqls["select-vendor"]

	l := log.WithField("test", "SelectVendor")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result types.Vendor
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery(querypat).
					WillReturnRows(sqlmock.
						NewRows([]string{"name"}).
						AddRow("vendor 0").
						AddRow("vendor 1").
						AddRow("vendor 2"))
				return db
			},
			id:     "0",
			result: types.Vendor{"0", "vendor 0"},
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
			}).SelectVendor(context.Background(), tc.id, "Test_SelectVendors")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_InsertVendor(t *testing.T) {
	t.Parallel()

	var querypat = sqls["insert-vendor"]

	l := log.WithField("test", "InsertVendor")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result types.Vendor
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
			result: types.Vendor{"30313233-3435-3637-3839-616263646566", "vendor 0"},
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
			result: types.Vendor{"30313233-3435-3637-3839-616263646566", "vendor 0"},
			err:    fmt.Errorf("vendor was not added"),
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
			result: types.Vendor{"30313233-3435-3637-3839-616263646566", "vendor 0"},
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
			}).InsertVendor(
				context.Background(),
				types.Vendor{tc.id, "vendor " + string(tc.id)},
				"Test_InsertVendors")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_UpdateVendor(t *testing.T) {
	t.Parallel()

	var querypat = sqls["update-vendor"]

	l := log.WithField("test", "UpdateVendor")

	tcs := map[string]struct {
		db  getMockDB
		id  types.UUID
		v   types.Vendor
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
			err: fmt.Errorf("vendor was not updated: '0'"),
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
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).UpdateVendor(
				context.Background(),
				tc.id,
				types.Vendor{Name: "vendor " + string(tc.id)},
				"Test_UpdateVendors")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_DeleteVendor(t *testing.T) {
	t.Parallel()

	var querypat = sqls["delete-vendor"]

	l := log.WithField("test", "deleteVendor")

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
			err: fmt.Errorf("vendor could not be deleted: '0'"),
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
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).DeleteVendor(
				context.Background(),
				tc.id,
				"Test_DeleteVendors")

			require.Equal(t, tc.err, err)
		})
	}
}
