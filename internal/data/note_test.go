package data

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/jsmit257/huautla/types"
)

var (
	_notes = []note{
		{UUID: "noteuuid 0", Note: "notenote 0", MTime: wwtbn, CTime: wwtbn},
		{UUID: "noteuuid 1", Note: "notenote 1", MTime: wwtbn, CTime: wwtbn},
		{UUID: "noteuuid 2", Note: "notenote 2", MTime: wwtbn, CTime: wwtbn},
	}
	noteFields = row{"id", "note", "mtime", "ctime"}
	noteValues = [][]driver.Value{
		{_notes[0].UUID, _notes[0].Note, _notes[0].MTime, _notes[0].CTime},
		{_notes[1].UUID, _notes[1].Note, _notes[1].MTime, _notes[1].CTime},
		{_notes[2].UUID, _notes[2].Note, _notes[2].MTime, _notes[2].CTime},
	}
)

func Test_GetNotes(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "Test_GetNotes")

	whenwillthenbenow := time.Now().UTC()

	set := map[string]struct {
		db     func() *sql.DB
		result []types.Note
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "note", "mtime", "ctime"}).
						AddRow("id-0", "note 0", whenwillthenbenow, whenwillthenbenow).
						AddRow("id-1", "note 1", whenwillthenbenow, whenwillthenbenow).
						AddRow("id-2", "note 2", whenwillthenbenow, whenwillthenbenow))
				return db
			},
			result: []types.Note{
				{UUID: "id-0", Note: "note 0", MTime: whenwillthenbenow, CTime: whenwillthenbenow},
				{UUID: "id-1", Note: "note 1", MTime: whenwillthenbenow, CTime: whenwillthenbenow},
				{UUID: "id-2", Note: "note 2", MTime: whenwillthenbenow, CTime: whenwillthenbenow},
			},
		},
		"db_error": {
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

	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			result, err := (&Conn{
				query:        v.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", k),
			}).GetNotes(context.Background(), "0", "Test_GetNotes")

			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}

func Test_AddNote(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "Test_AddNote")

	set := map[string]struct {
		db     func() *sql.DB
		result int
		n      types.Note
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
			result: 1,
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("note was not added"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
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
			err: fmt.Errorf("some error"),
		},
	}

	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			result, err := (&Conn{
				query:        v.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", k),
			}).AddNote(context.Background(), "0", []types.Note{}, v.n, "Test_AddNote")

			require.Equal(t, v.err, err)
			require.Equal(t, v.result, len(result))
		})
	}
}

func Test_ChangeNote(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "Test_ChangeNote")

	set := map[string]struct {
		db  func() *sql.DB
		n   types.Note
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
			n: types.Note{UUID: "0", Note: "note"},
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("note was not changed"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
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
			err: fmt.Errorf("some error"),
		},
	}

	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			_, err := (&Conn{
				query:        v.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", k),
			}).ChangeNote(
				context.Background(),
				[]types.Note{
					{UUID: "1"},
					{UUID: "0"},
				},
				v.n,
				"Test_ChangeNote")

			require.Equal(t, v.err, err)
		})
	}
}

func Test_RemoveNote(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "Test_RemoveNote")

	set := map[string]struct {
		db     func() *sql.DB
		result []types.Note
		id     types.UUID
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
			result: []types.Note{{UUID: "1"}},
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			result: []types.Note{
				{UUID: "1"},
				{UUID: "0"},
			},
			err: fmt.Errorf("note could not be removed"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			result: []types.Note{
				{UUID: "1"},
				{UUID: "0"},
			},
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
			result: []types.Note{
				{UUID: "1"},
				{UUID: "0"},
			},
			err: fmt.Errorf("some error"),
		},
	}

	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			result, err := (&Conn{
				query:        v.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", k),
			}).RemoveNote(
				context.Background(),
				[]types.Note{
					{UUID: "1"},
					{UUID: "0"},
				},
				"0",
				"Test_RemoveNote")

			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}
