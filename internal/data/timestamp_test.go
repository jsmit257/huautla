package data

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jsmit257/huautla/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func Test_UpdateTimestamps(t *testing.T) {
	t.Parallel()

	l := logrus.WithField("test", "UpdateSubstrate")

	tcs := map[string]struct {
		db    getMockDB
		table string
		id    types.UUID
		flds  []string
		org   *time.Time
		err   error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id:   "0",
			flds: []string{"mtime"},
			org:  &wwtbn,
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:   "0",
			flds: []string{"mtime"},
			org:  &wwtbn,
			err:  fmt.Errorf("timestamps were not updated"),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			id:   "0",
			flds: []string{"mtime"},
			org:  &wwtbn,
			err:  fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			id:   "0",
			flds: []string{"mtime"},
			org:  &wwtbn,
			err:  fmt.Errorf("some error"),
		},
		"validate_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id:  "0",
			err: fmt.Errorf("no fields specified for update"),
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
			}).UpdateTimestamps(
				context.Background(),
				tc.table,
				tc.id,
				types.Timestamp{
					Fields: tc.flds,
					Origin: tc.org,
				})

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_Undelete(t *testing.T) {
	t.Parallel()

	l := logrus.WithField("test", "DeleteSubstrate")

	tcs := map[string]struct {
		db    getMockDB
		table string
		id    types.UUID
		err   error
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
			err: fmt.Errorf("record could not be undeleted"),
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
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).Undelete(
				context.Background(),
				tc.table,
				tc.id)

			require.Equal(t, tc.err, err)
		})
	}
}
