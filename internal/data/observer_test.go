package data

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/jsmit257/huautla/types"
)

var (
	_events = []event{
		{UUID: "eventuuid 0", Temperature: -40.0, Humidity: 0, MTime: wwtbn, CTime: wwtbn, EventType: types.EventType{UUID: "typeuuid 0", Name: "Clone", Severity: "e0.EventType.Severity", Stage: types.Stage{UUID: "e0.EventType.Stage.UUID", Name: "e0.EventType.Stage.Name"}}},
		{UUID: "eventuuid 1", Temperature: 451, Humidity: 50, MTime: wwtbn, CTime: wwtbn, EventType: types.EventType{UUID: "typeuuid 1", Name: "e1.EventType.Name", Severity: "e1.EventType.Severity", Stage: types.Stage{UUID: "e1.EventType.Stage.UUID", Name: "e1.EventType.Stage.Name"}}},
		{UUID: "eventuuid 2", Temperature: 10.0, Humidity: 100, MTime: wwtbn, CTime: wwtbn, EventType: types.EventType{UUID: "typeuuid 2", Name: "e2.EventType.Name", Severity: "e2.EventType.Severity", Stage: types.Stage{UUID: "e2.EventType.Stage.UUID", Name: "e2.EventType.Stage.Name"}}},
	}
	eventFields = row{
		"id",
		"temperature",
		"humidity",
		"mtime",
		"ctime",
		"eventtype_uuid",
		"event_severity",
		"eventtype_name",
		"stage_uuid",
		"stage_name",
	}
	eventValues = [][]driver.Value{
		{_events[0].UUID, _events[0].Temperature, _events[0].Humidity, _events[0].MTime, _events[0].CTime, _events[0].EventType.UUID, _events[0].EventType.Name, _events[0].EventType.Severity, _events[0].EventType.Stage.UUID, _events[0].EventType.Stage.Name},
		{_events[1].UUID, _events[1].Temperature, _events[1].Humidity, _events[1].MTime, _events[1].CTime, _events[1].EventType.UUID, _events[1].EventType.Name, _events[1].EventType.Severity, _events[1].EventType.Stage.UUID, _events[1].EventType.Stage.Name},
		{_events[2].UUID, _events[2].Temperature, _events[2].Humidity, _events[2].MTime, _events[2].CTime, _events[2].EventType.UUID, _events[2].EventType.Name, _events[2].EventType.Severity, _events[2].EventType.Stage.UUID, _events[2].EventType.Stage.Name},
	}

	// nap == NotesAndPhotos; it's not really implemented for test
	napFields = row{"uuid", "note_uuid", "note_note", "note_mtime", "note_ctime", "photo_uuid", "filename", "photo_mtime", "photo_ctime", "photonote_uuid", "photonote_note", "photonote_mtime", "photonote_ctime"}
	napValues = [][]driver.Value{
		{"0", "0", "note 0", wwtbn, wwtbn, "0", "photo 0", wwtbn, wwtbn, "0", "photonote 0", wwtbn, wwtbn},
		{"1", "1", "note 1", wwtbn, wwtbn, nil, nil, nil, nil, nil, nil, nil, nil},
		{"2", nil, nil, nil, nil, "2", "photo 2", wwtbn, wwtbn, "2", "photonote 2", wwtbn, wwtbn},
		{"2", nil, nil, nil, nil, "2", "photo 2", wwtbn, wwtbn, "3", "photonote 3", wwtbn, wwtbn},
	}
)

func Test_SelectByEventType(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "Test_GetLifecycleEvents")

	tcs := map[string]struct {
		db     getMockDB
		result []types.Event
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock, eventFields.set(eventValues...))
				return db
			},
			result: []types.Event{
				types.Event(_events[0]),
				types.Event(_events[1]),
				types.Event(_events[2]),
			},
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			result: []types.Event{},
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
			}).SelectByEventType(context.Background(), types.EventType{}, "Test_SelectAllEventTypes")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_SelectEvent(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "SelectEvent")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result types.Event
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock, eventFields.set(eventValues[0]))
				return db
			},
			id:     "0",
			result: types.Event(_events[0]),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock, eventFields.fail())
				return db
			},
			err: eventFields.err(),
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
			}).SelectEvent(context.Background(), tc.id, "Test_SelectEvent")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_UpdateEvent(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "UpdateEvent")

	e0 := types.Event{UUID: "0"}

	modifyevent := func(e types.Event) types.Event {
		e.Temperature = 100.0
		return e
	}

	tcs := map[string]struct {
		db  getMockDB
		evt types.Event
		err error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			evt: modifyevent(e0),
		},
		"exec_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			evt: modifyevent(e0),
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			evt: modifyevent(e0),
			err: fmt.Errorf("some error"),
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			evt: modifyevent(e0),
			err: fmt.Errorf("event was not changed"),
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
			}).UpdateEvent(
				context.Background(),
				tc.evt,
				"Test_UpdateEvent")

			require.Equal(t, tc.err, err)
		})
	}
}
func Test_notesAndPhotos(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "notesAndPhotos")

	tcs := map[string]struct {
		db  getMockDB
		in  []types.Event
		out []types.Event
		err error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				napFields.mock(mock, napValues...)
				return db
			},
			in: []types.Event{
				{UUID: "0"},
				{UUID: "1"},
				{UUID: "2"},
			},
			out: []types.Event{
				{
					UUID: "0",
					Notes: []types.Note{
						{
							UUID:  "0",
							Note:  "note 0",
							MTime: wwtbn,
							CTime: wwtbn,
						},
					},
					Photos: []types.Photo{
						{
							UUID:     "0",
							Filename: "photo 0",
							MTime:    wwtbn,
							CTime:    wwtbn,
							Notes: []types.Note{
								{
									UUID:  "0",
									Note:  "photonote 0",
									MTime: wwtbn,
									CTime: wwtbn,
								},
							},
						},
					},
				},
				{
					UUID: "1",
					Notes: []types.Note{
						{
							UUID:  "1",
							Note:  "note 1",
							MTime: wwtbn,
							CTime: wwtbn,
						},
					},
				},
				{
					UUID: "2",
					Photos: []types.Photo{
						{
							UUID:     "2",
							Filename: "photo 2",
							MTime:    wwtbn,
							CTime:    wwtbn,
							Notes: []types.Note{
								{
									UUID:  "3",
									Note:  "photonote 3",
									MTime: wwtbn,
									CTime: wwtbn,
								},
								{
									UUID:  "2",
									Note:  "photonote 2",
									MTime: wwtbn,
									CTime: wwtbn,
								},
							},
						},
					},
				},
			},
		},
		"events_nil": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				return db
			},
		},
		"event_error": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			in:  []types.Event{{}},
			out: []types.Event{{}},
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
			}).notesAndPhotos(context.Background(), tc.in, "0", "Test_notesAndPhotos")
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.out, tc.in)
		})
	}
}
