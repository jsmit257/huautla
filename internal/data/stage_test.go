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

func Test_SelectAllStages(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "SelectAllStages")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result []types.Stage
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name"}).
						AddRow("0", "stage 0").
						AddRow("1", "stage 1").
						AddRow("2", "stage 2"))
				return db
			},
			result: []types.Stage{
				{UUID: "0", Name: "stage 0"},
				{UUID: "1", Name: "stage 1"},
				{UUID: "2", Name: "stage 2"},
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
			}).SelectAllStages(context.Background(), "Test_SelectAllStages")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_SelectStage(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "SelectStage")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result types.Stage
		err    error
	}{
		"happy_path": {db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
			mock.ExpectQuery("").WillReturnRows(
				sqlmock.NewRows([]string{"name"}).AddRow("stage 0"))
			return db
		}, id: "0", result: types.Stage{UUID: "0", Name: "stage 0"}},
		"no_row_returned": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows([]string{"name"}))
				return db
			},
			id:     "0",
			result: types.Stage{UUID: "0", Name: ""},
			err:    fmt.Errorf("sql: no rows in result set"),
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
			}).SelectStage(context.Background(), tc.id, "Test_SelectStages")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_InsertStage(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "InsertStage")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result types.Stage
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id:     "0",
			result: types.Stage{UUID: "30313233-3435-3637-3839-616263646566", Name: "stage 0"},
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:     "0",
			result: types.Stage{UUID: "30313233-3435-3637-3839-616263646566", Name: "stage 0"},
			err:    fmt.Errorf("stage was not added"),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			id:     "0",
			result: types.Stage{UUID: "30313233-3435-3637-3839-616263646566", Name: "stage 0"},
			err:    fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			id:     "0",
			result: types.Stage{UUID: "30313233-3435-3637-3839-616263646566", Name: "stage 0"},
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
			}).InsertStage(
				context.Background(),
				types.Stage{UUID: tc.id, Name: "stage " + string(tc.id)},
				"Test_InsertStages")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_UpdateStage(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "UpdateStage")

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
			err: fmt.Errorf("stage was not updated: '0'"),
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
			}).UpdateStage(
				context.Background(),
				tc.id,
				types.Stage{Name: "stage " + string(tc.id)},
				"Test_UpdateStages")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_DeleteStage(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "DeleteStage")

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
			err: fmt.Errorf("stage could not be deleted: '0'"),
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
			}).DeleteStage(
				context.Background(),
				tc.id,
				"Test_DeleteStages")

			require.Equal(t, tc.err, err)
		})
	}
}
