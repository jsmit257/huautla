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

func Test_GetGenerationEvents(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "Test_GetGenerationEvents")

	e0, e1, e2 := types.Event{UUID: "0"},
		types.Event{UUID: "1"},
		types.Event{UUID: "2"}

	tcs := map[string]struct {
		db     getMockDB
		result []types.Event
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "temperature", "humidity", "mtime", "ctime", "eventtype_uuid", "event_severity", "eventtype_name", "stage_uuid", "stage_name"}).
						AddRow(e0.UUID, e0.Temperature, e0.Humidity, e0.MTime, e0.CTime, e0.EventType.UUID, e0.EventType.Name, e0.EventType.Severity, e0.EventType.Stage.UUID, e0.EventType.Stage.Name).
						AddRow(e1.UUID, e1.Temperature, e1.Humidity, e1.MTime, e1.CTime, e1.EventType.UUID, e1.EventType.Name, e1.EventType.Severity, e1.EventType.Stage.UUID, e1.EventType.Stage.Name).
						AddRow(e2.UUID, e2.Temperature, e2.Humidity, e2.MTime, e2.CTime, e2.EventType.UUID, e2.EventType.Name, e2.EventType.Severity, e2.EventType.Stage.UUID, e2.EventType.Stage.Name))
				return db
			},
			result: []types.Event{e0, e1, e2},
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err:    fmt.Errorf("some error"),
			result: []types.Event{},
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		lc := &types.Generation{}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).GetGenerationEvents(context.Background(), lc, "Test_SelectAllEventTypes")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, lc.Events)
		})
	}
}

func Test_AddGenerationEvent(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "AddGenerationEvent")

	e0, e1, e2 := types.Event{UUID: "0"},
		types.Event{UUID: "1"},
		types.Event{UUID: "2"}

	tcs := map[string]struct {
		db     getMockDB
		evts   []types.Event
		evt    types.Event
		result []types.Event
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				etFields.mock(mock, etValues[0])
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			evts:   []types.Event{e0, e1},
			evt:    e2,
			result: []types.Event{e0, e1, e2},
		},
		"modified_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				etFields.mock(mock, etValues[0])
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			evts:   []types.Event{e0, e1},
			evt:    e2,
			result: []types.Event{e0, e1, e2},
			err:    fmt.Errorf("couldn't update Generation.mtime"),
		},
		"eventtype_error": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err: fmt.Errorf("couldn't fetch eventtype"),
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("event was not added"),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
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
			}).AddGenerationEvent(
				context.Background(),
				&types.Generation{Events: tc.evts},
				tc.evt,
				"Test_AddGenerationEvent")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_ChangeGenerationEvent(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "ChangeGenerationEvent")

	e0, e1, e2 := types.Event{UUID: "0"},
		types.Event{UUID: "1"},
		types.Event{UUID: "2"}

	modifyevent := func(e types.Event) types.Event {
		e.Temperature = 100.0
		return e
	}

	tcs := map[string]struct {
		db     getMockDB
		evts   []types.Event
		evt    types.Event
		result []types.Event
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				etFields.mock(mock, etValues[0])
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			evts:   []types.Event{e0, e1, e2},
			evt:    modifyevent(e1),
			result: []types.Event{e0, modifyevent(e1), e2},
		},
		"modified_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				etFields.mock(mock, etValues[0])
				mock.ExpectExec("").WillReturnError(fmt.Errorf("couldn't update Generation.mtime"))
				return db
			},
			evts:   []types.Event{e0, e1, e2},
			evt:    modifyevent(e1),
			result: []types.Event{e0, modifyevent(e1), e2},
			err:    fmt.Errorf("couldn't update Generation.mtime"),
		},
		"eventtype_error": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err: fmt.Errorf("couldn't fetch eventtype"),
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("event was not changed"),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			err: fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).ChangeGenerationEvent(
				context.Background(),
				&types.Generation{Events: tc.evts},
				tc.evt,
				"Test_ChangeGenerationEvent")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_RemoveGenerationEvent(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "RemoveGenerationEvent")

	e0, e1, e2 := types.Event{UUID: "0"},
		types.Event{UUID: "1"},
		types.Event{UUID: "2"}

	tcs := map[string]struct {
		db     getMockDB
		evts   []types.Event
		id     types.UUID
		result []types.Event
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			evts:   []types.Event{e0, e1, e2},
			id:     e1.UUID,
			result: []types.Event{e0, e2},
		},
		"modified_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("").WillReturnError(fmt.Errorf("couldn't update Generation.mtime"))
				return db
			},
			evts:   []types.Event{e0, e1, e2},
			id:     e1.UUID,
			result: []types.Event{e0, e2},
			err:    fmt.Errorf("couldn't update Generation.mtime"),
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("event could not be removed"),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
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
			}).RemoveGenerationEvent(
				context.Background(),
				&types.Generation{Events: tc.evts},
				tc.id,
				"Test_RemoveGenerationEvent")

			require.Equal(t, tc.err, err)
		})
	}
}
