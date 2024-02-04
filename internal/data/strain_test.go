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

func Test_SelectAllStrains(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "Test_SelectAllStrains")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result []types.Strain
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name", "vendor_uuid", "vendor_name"}).
						AddRow("0", "strain 0", "0", "vendor 0").
						AddRow("1", "strain 1", "1", "vendor 1").
						AddRow("2", "strain 2", "1", "vendor 1"))
				return db
			},
			result: []types.Strain{
				{UUID: "0", Name: "strain 0", Vendor: types.Vendor{UUID: "0", Name: "vendor 0"}, Attributes: nil},
				{UUID: "1", Name: "strain 1", Vendor: types.Vendor{UUID: "1", Name: "vendor 1"}, Attributes: nil},
				{UUID: "2", Name: "strain 2", Vendor: types.Vendor{UUID: "1", Name: "vendor 1"}, Attributes: nil},
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
			}).SelectAllStrains(context.Background(), "Test_SelectAllStrains")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_SelectStrain(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "SelectStrain")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result types.Strain
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"name", "vendor_uuid", "vendor_name"}).
						AddRow("strain 0", "0", "vendor 0"))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name", "value"}).
						AddRow("0", "name 0", "value 0").
						AddRow("1", "name 1", "value 1").
						AddRow("2", "name 2", "value 2"))
				return db
			},
			id: "0",
			result: types.Strain{UUID: "0", Name: "strain 0", Vendor: types.Vendor{UUID: "0", Name: "vendor 0"}, Attributes: []types.StrainAttribute{
				{UUID: "0", Name: "name 0", Value: "value 0"},
				{UUID: "1", Name: "name 1", Value: "value 1"},
				{UUID: "2", Name: "name 2", Value: "value 2"},
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
			}).SelectStrain(context.Background(), tc.id, "Test_SelectStrain")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_InsertStrain(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "InsertStrain")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result types.Strain
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
			result: types.Strain{UUID: "30313233-3435-3637-3839-616263646566", Name: "strain 0", Vendor: types.Vendor{}, Attributes: nil},
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
			result: types.Strain{UUID: "30313233-3435-3637-3839-616263646566", Name: "strain 0", Vendor: types.Vendor{}, Attributes: nil},
			err:    fmt.Errorf("strain was not added"),
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
			result: types.Strain{UUID: "30313233-3435-3637-3839-616263646566", Name: "strain 0", Vendor: types.Vendor{}, Attributes: nil},
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
			result: types.Strain{UUID: "30313233-3435-3637-3839-616263646566", Name: "strain 0", Vendor: types.Vendor{}, Attributes: nil},
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
			}).InsertStrain(
				context.Background(),
				types.Strain{UUID: tc.id, Name: "strain " + string(tc.id), Vendor: types.Vendor{}, Attributes: nil},
				"Test_InsertStrains")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_UpdateStrain(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "UpdateStrain")

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
			err: fmt.Errorf("strain was not updated: '0'"),
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
			}).UpdateStrain(
				context.Background(),
				tc.id,
				types.Strain{Name: "strain " + string(tc.id)},
				"Test_UpdateStrains")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_DeleteStrain(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "DeleteStrain")

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
			err: fmt.Errorf("strain could not be deleted: '0'"),
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
			}).DeleteStrain(
				context.Background(),
				tc.id,
				"Test_DeleteStrain")

			require.Equal(t, tc.err, err)
		})
	}
}
