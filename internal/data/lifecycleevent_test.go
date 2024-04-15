package data

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jsmit257/huautla/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func Test_GetLifecycleEvents(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "Test_GetLifecycleEvents")

	whenwillthenbenow := time.Now().UTC()

	e0, e1, e2 := types.Event{UUID: "0"},
		types.Event{UUID: "1"},
		types.Event{UUID: "2"}

	oneNote := e1
	oneNote.UUID = "1 note"
	oneNote.Notes = []types.Note{
		{UUID: "note0", Note: "note 0", MTime: whenwillthenbenow, CTime: whenwillthenbenow},
	}
	twoNotes := e1
	twoNotes.UUID = "2 notes"
	twoNotes.Notes = []types.Note{
		{UUID: "note2", Note: "note 2", MTime: whenwillthenbenow, CTime: whenwillthenbenow},
		{UUID: "note1", Note: "note 1", MTime: whenwillthenbenow, CTime: whenwillthenbenow},
	}

	photo := e2
	photo.UUID = "photos"
	photo.Photos = []types.Photo{
		{UUID: "id-0", Filename: "photo 0", CTime: whenwillthenbenow, Notes: []types.Note{
			{UUID: "note-1", Note: "note 1", MTime: whenwillthenbenow, CTime: whenwillthenbenow},
			{UUID: "note-0", Note: "note 0", MTime: whenwillthenbenow, CTime: whenwillthenbenow},
		}},
		{UUID: "id-1", Filename: "photo 1", CTime: whenwillthenbenow},
		{UUID: "id-2", Filename: "photo 2", CTime: whenwillthenbenow},
	}
	tcs := map[string]struct {
		db     getMockDB
		result []types.Event
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "temperature", "humidity", "mtime", "ctime", "eventtype_uuid", "event_severity", "eventtype_name", "stage_uuid", "stage_name", "note_id", "note", "note_mtime", "note_ctime", "has_photos"}).
						AddRow(e0.UUID, e0.Temperature, e0.Humidity, e0.MTime, e0.CTime, e0.EventType.UUID, e0.EventType.Name, e0.EventType.Severity, e0.EventType.Stage.UUID, e0.EventType.Stage.Name, nil, nil, nil, nil, 0).
						AddRow(e1.UUID, e1.Temperature, e1.Humidity, e1.MTime, e1.CTime, e1.EventType.UUID, e1.EventType.Name, e1.EventType.Severity, e1.EventType.Stage.UUID, e1.EventType.Stage.Name, nil, nil, nil, nil, 0).
						AddRow("1 note", e1.Temperature, e1.Humidity, e1.MTime, e1.CTime, e1.EventType.UUID, e1.EventType.Name, e1.EventType.Severity, e1.EventType.Stage.UUID, e1.EventType.Stage.Name, "note0", "note 0", whenwillthenbenow, whenwillthenbenow, 0).
						AddRow("2 notes", e1.Temperature, e1.Humidity, e1.MTime, e1.CTime, e1.EventType.UUID, e1.EventType.Name, e1.EventType.Severity, e1.EventType.Stage.UUID, e1.EventType.Stage.Name, "note1", "note 1", whenwillthenbenow, whenwillthenbenow, 0).
						AddRow("2 notes", e1.Temperature, e1.Humidity, e1.MTime, e1.CTime, e1.EventType.UUID, e1.EventType.Name, e1.EventType.Severity, e1.EventType.Stage.UUID, e1.EventType.Stage.Name, "note2", "note 2", whenwillthenbenow, whenwillthenbenow, 0).
						AddRow(e2.UUID, e2.Temperature, e2.Humidity, e2.MTime, e2.CTime, e2.EventType.UUID, e2.EventType.Name, e2.EventType.Severity, e2.EventType.Stage.UUID, e2.EventType.Stage.Name, nil, nil, nil, nil, 0).
						AddRow("photos", e2.Temperature, e2.Humidity, e2.MTime, e2.CTime, e2.EventType.UUID, e2.EventType.Name, e2.EventType.Severity, e2.EventType.Stage.UUID, e2.EventType.Stage.Name, nil, nil, nil, nil, 1))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "filename", "ctime", "note_uuid", "note", "note_mtime", "note_ctime"}).
						AddRow("id-0", "photo 0", whenwillthenbenow, "note-0", "note 0", whenwillthenbenow, whenwillthenbenow).
						AddRow("id-0", "photo 0", whenwillthenbenow, "note-1", "note 1", whenwillthenbenow, whenwillthenbenow).
						AddRow("id-1", "photo 1", whenwillthenbenow, nil, nil, nil, nil).
						AddRow("id-2", "photo 2", whenwillthenbenow, nil, nil, nil, nil))
				return db
			},
			result: []types.Event{
				e0,
				e1,
				oneNote,
				twoNotes,
				e2,
				photo,
			},
		},
		"photo_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "temperature", "humidity", "mtime", "ctime", "eventtype_uuid", "event_severity", "eventtype_name", "stage_uuid", "stage_name", "note_id", "note", "note_mtime", "note_ctime", "has_photos"}).
						AddRow(e0.UUID, e0.Temperature, e0.Humidity, e0.MTime, e0.CTime, e0.EventType.UUID, e0.EventType.Name, e0.EventType.Severity, e0.EventType.Stage.UUID, e0.EventType.Stage.Name, nil, nil, nil, nil, 0).
						AddRow(e1.UUID, e1.Temperature, e1.Humidity, e1.MTime, e1.CTime, e1.EventType.UUID, e1.EventType.Name, e1.EventType.Severity, e1.EventType.Stage.UUID, e1.EventType.Stage.Name, nil, nil, nil, nil, 0).
						AddRow("1 note", e1.Temperature, e1.Humidity, e1.MTime, e1.CTime, e1.EventType.UUID, e1.EventType.Name, e1.EventType.Severity, e1.EventType.Stage.UUID, e1.EventType.Stage.Name, "note0", "note 0", whenwillthenbenow, whenwillthenbenow, 0).
						AddRow("2 notes", e1.Temperature, e1.Humidity, e1.MTime, e1.CTime, e1.EventType.UUID, e1.EventType.Name, e1.EventType.Severity, e1.EventType.Stage.UUID, e1.EventType.Stage.Name, "note1", "note 1", whenwillthenbenow, whenwillthenbenow, 0).
						AddRow("2 notes", e1.Temperature, e1.Humidity, e1.MTime, e1.CTime, e1.EventType.UUID, e1.EventType.Name, e1.EventType.Severity, e1.EventType.Stage.UUID, e1.EventType.Stage.Name, "note2", "note 2", whenwillthenbenow, whenwillthenbenow, 0).
						AddRow(e2.UUID, e2.Temperature, e2.Humidity, e2.MTime, e2.CTime, e2.EventType.UUID, e2.EventType.Name, e2.EventType.Severity, e2.EventType.Stage.UUID, e2.EventType.Stage.Name, nil, nil, nil, nil, 0).
						AddRow("photos", e2.Temperature, e2.Humidity, e2.MTime, e2.CTime, e2.EventType.UUID, e2.EventType.Name, e2.EventType.Severity, e2.EventType.Stage.UUID, e2.EventType.Stage.Name, nil, nil, nil, nil, 1))
				mock.ExpectQuery("").
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			result: []types.Event{
				e0,
				e1,
				oneNote,
				twoNotes,
				e2,
			},
			err: fmt.Errorf("some error"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectQuery("").
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err:    fmt.Errorf("some error"),
			result: []types.Event{},
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		lc := &types.Lifecycle{}
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).GetLifecycleEvents(context.Background(), lc, "Test_SelectAllEventTypes")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, lc.Events)
		})
	}
}

func Test_AddLifecycleEvent(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "AddLifecycleEvent")

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
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"name", "severity", "stage_uuid", "stage_name"}).
						AddRow("type 0", "Info", "0", "stage 0"))
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			evts:   []types.Event{e0, e1},
			evt:    e2,
			result: []types.Event{e0, e1, e2},
		},
		"modified_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"name", "severity", "stage_uuid", "stage_name"}).
						AddRow("type 0", "Info", "0", "stage 0"))
				mock.
					ExpectExec("").
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			evts:   []types.Event{e0, e1},
			evt:    e2,
			result: []types.Event{e0, e1, e2},
			err:    fmt.Errorf("couldn't update lifecycle.mtime"),
		},
		"eventtype_error": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.
					ExpectQuery("").
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err: fmt.Errorf("couldn't fetch eventtype"),
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("event was not added"),
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

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).AddLifecycleEvent(
				context.Background(),
				&types.Lifecycle{Events: tc.evts},
				tc.evt,
				"Test_AddLifecycleEvent")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_ChangeLifecycleEvent(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "ChangeLifecycleEvent")

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
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"name", "severity", "stage_uuid", "stage_name"}).
						AddRow("type 0", "Info", "0", "stage 0"))
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			evts:   []types.Event{e0, e1, e2},
			evt:    modifyevent(e1),
			result: []types.Event{e0, modifyevent(e1), e2},
		},
		"modified_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"name", "severity", "stage_uuid", "stage_name"}).
						AddRow("type 0", "Info", "0", "stage 0"))
				mock.
					ExpectExec("").
					WillReturnError(fmt.Errorf("couldn't update lifecycle.mtime"))
				return db
			},
			evts:   []types.Event{e0, e1, e2},
			evt:    modifyevent(e1),
			result: []types.Event{e0, modifyevent(e1), e2},
			err:    fmt.Errorf("couldn't update lifecycle.mtime"),
		},
		"eventtype_error": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.
					ExpectQuery("").
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err: fmt.Errorf("couldn't fetch eventtype"),
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("event was not changed"),
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

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).ChangeLifecycleEvent(
				context.Background(),
				&types.Lifecycle{Events: tc.evts},
				tc.evt,
				"Test_ChangeLifecycleEvent")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_RemoveLifecycleEvent(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "RemoveLifecycleEvent")

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
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			evts:   []types.Event{e0, e1, e2},
			id:     e1.UUID,
			result: []types.Event{e0, e2},
		},
		"modified_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.
					ExpectExec("").
					WillReturnError(fmt.Errorf("couldn't update lifecycle.mtime"))
				return db
			},
			evts:   []types.Event{e0, e1, e2},
			id:     e1.UUID,
			result: []types.Event{e0, e2},
			err:    fmt.Errorf("couldn't update lifecycle.mtime"),
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("event could not be removed"),
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

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).RemoveLifecycleEvent(
				context.Background(),
				&types.Lifecycle{Events: tc.evts},
				tc.id,
				"Test_RemoveLifecycleEvent")

			require.Equal(t, tc.err, err)
		})
	}
}
