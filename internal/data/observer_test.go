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
		{UUID: e0.UUID, Temperature: e0.Temperature, Humidity: e0.Humidity, MTime: wwtbn, CTime: wwtbn, EventType: types.EventType{UUID: e0.EventType.UUID, Name: e0.EventType.Name, Severity: e0.EventType.Severity, Stage: types.Stage{UUID: e0.EventType.Stage.UUID, Name: e0.EventType.Stage.Name}}},
		{UUID: e1.UUID, Temperature: e1.Temperature, Humidity: e1.Humidity, MTime: wwtbn, CTime: wwtbn, EventType: types.EventType{UUID: e1.EventType.UUID, Name: e1.EventType.Name, Severity: e1.EventType.Severity, Stage: types.Stage{UUID: e1.EventType.Stage.UUID, Name: e1.EventType.Stage.Name}}},
		{UUID: e2.UUID, Temperature: e2.Temperature, Humidity: e2.Humidity, MTime: wwtbn, CTime: wwtbn, EventType: types.EventType{UUID: e2.EventType.UUID, Name: e2.EventType.Name, Severity: e2.EventType.Severity, Stage: types.Stage{UUID: e2.EventType.Stage.UUID, Name: e2.EventType.Stage.Name}}},
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

	// nap == NotesAndPhotos; it's note really implemented for test
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

	e0, e1, e2 := types.Event{UUID: "0"},
		types.Event{UUID: "1"},
		types.Event{UUID: "2"}

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
						NewRows([]string{"id", "temperature", "humidity", "mtime", "ctime", "eventtype_uuid", "event_severity", "eventtype_name", "stage_uuid", "stage_name"}).
						AddRow(e0.UUID, e0.Temperature, e0.Humidity, e0.MTime, e0.CTime, e0.EventType.UUID, e0.EventType.Name, e0.EventType.Severity, e0.EventType.Stage.UUID, e0.EventType.Stage.Name).
						AddRow(e1.UUID, e1.Temperature, e1.Humidity, e1.MTime, e1.CTime, e1.EventType.UUID, e1.EventType.Name, e1.EventType.Severity, e1.EventType.Stage.UUID, e1.EventType.Stage.Name).
						AddRow(e2.UUID, e2.Temperature, e2.Humidity, e2.MTime, e2.CTime, e2.EventType.UUID, e2.EventType.Name, e2.EventType.Severity, e2.EventType.Stage.UUID, e2.EventType.Stage.Name))
				return db
			},
			result: []types.Event{e0, e1, e2},
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectQuery("").
					WillReturnError(fmt.Errorf("some error"))
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
				query:        tc.db(),
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

	e0 := types.Event{UUID: "0"}

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result types.Event
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"temperature", "humidity", "mtime", "ctime", "eventtype_uuid", "eventtype_severity", "eventtype_name", "stage_uuid", "stage_name"}).
						AddRow(e0.Temperature, e0.Humidity, e0.MTime, e0.CTime, e0.EventType.UUID, e0.EventType.Name, e0.EventType.Severity, e0.EventType.Stage.UUID, e0.EventType.Stage.Name))
				return db
			},
			id:     "0",
			result: e0,
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
			}).SelectEvent(context.Background(), tc.id, "Test_SelectEvent")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
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
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
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
			db: func() *sql.DB {
				db, _, _ := sqlmock.New()
				return db
			},
		},
		"event_error": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectQuery("").
					WillReturnError(fmt.Errorf("some error"))
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
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).notesAndPhotos(context.Background(), tc.in, "0", "Test_notesAndPhotos")
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.out, tc.in)
		})
	}
}
