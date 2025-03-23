package data

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jsmit257/huautla/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var (
	_src = types.Source{
		UUID:      "sourceuuid",
		Type:      "sourcetype",
		Lifecycle: nil,
		Strain:    types.Strain(_strain),
	}
	srcFields = row{"uuid", "type", "progenitor_uuid", "lifecycle_uuid", "strain_uuid", "strain_name", "&strain_species", "strain_ctime", "strain_dtime", "strain_vendor_id", "strain_vendor_name", "strain_vendor_website"}
	srcValues = [][]driver.Value{{_src.UUID, _src.Type, "pgid", nil, _src.Strain.UUID, _src.Strain.Name, _src.Strain.Species, _strain.CTime, _strain.DTime, _strain.Vendor.UUID, _strain.Vendor.Name, _strain.Vendor.Website}}
)

func Test_GetSources(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "GetSources")

	tcs := map[string]struct {
		db     getMockDB
		result []types.Source
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				srcFields.mock(mock, [][]driver.Value{
					{"uuid", "type", "progenitor_uuid", "lifecycle_uuid", "strain_uuid", "strain_name", "strain_species", wwtbn, nil, "strain_vendor_id", "strain_vendor_name", "strain_vendor_website"},
					{"uuid", "type", "progenitor_uuid", nil, "strain_uuid", "strain_name", "strain_species", wwtbn, nil, "strain_vendor_id", "strain_vendor_name", "strain_vendor_website"},
				}...)
				lcFields.mock(mock, []driver.Value{"uuid", "location", 0, 0, 0, 0, 0, 0, wwtbn, wwtbn, "0", "X.species", "strain 0", nil, wwtbn, nil, "x", "vendor x", "website", "gs", "gs", types.GrainType, "1", "vendor 1", "website", "bs", "bs", types.BulkType, "2", "vendor 2", "website"})
				eventFields.mock(mock, eventValues...)
				attrFields.mock(mock, attrValues...)
				ingFields.mock(mock, ingValues...)
				ingFields.mock(mock, ingValues...)

				return db
			},
			result: []types.Source{
				{
					UUID:      "uuid",
					Type:      "type",
					Lifecycle: &types.Lifecycle{UUID: "lifecycle_uuid"},
					Strain:    types.Strain(_strain),
				},
				{
					UUID:   "uuid",
					Type:   "type",
					Strain: types.Strain(_strain),
				},
			},
		},
		"db_error": {
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

			g := types.Generation{}

			err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).GetSources(context.Background(), &g, "Test_GetSources")

			require.Equal(t, tc.err, err)
			require.Equal(t, mustObject(tc.result), mustObject(g.Sources))
		})
	}
}

func Test_InsertSource(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "InsertSource")

	src := types.Source{
		UUID:      types.UUID(mockUUIDGen().String()),
		Lifecycle: &types.Lifecycle{Events: []types.Event{{}}}}

	tcs := map[string]struct {
		db     getMockDB
		origin string
		s      types.Source
		result types.Source
		err    error
	}{
		"happy_event_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			origin: "event",
			s:      src,
			result: src,
		},
		"happy_strain_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			origin: "strain",
			s:      src,
			result: src,
		},
		"origin_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				return nil
			},
			origin: "bad origin",
			err:    fmt.Errorf("only origins of type 'strain' and 'event' are allowed: 'bad origin'"),
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			origin: "strain",
			err:    fmt.Errorf("source was not added"),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			origin: "strain",
			err:    fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			origin: "strain",
			err:    fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			source, err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).InsertSource(
				context.Background(),
				"generation id",
				tc.origin,
				tc.s,
				"Test_AddStrainSource")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, source)
		})
	}
}

func Test_UpdateSource(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "UpdateSource")

	tcs := map[string]struct {
		db     getMockDB
		origin string
		s      types.Source
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			origin: "strain",
			s:      types.Source{UUID: "1", Type: "Spore"},
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			origin: "event",
			s: types.Source{
				Lifecycle: &types.Lifecycle{Events: []types.Event{{}}},
			},
			err: fmt.Errorf("source was not changed"),
		},
		"origin_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				return nil
			},
			origin: "bad origin",
			err:    fmt.Errorf("only origins of type 'strain' and 'event' are allowed: 'bad origin'"),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			origin: "strain",
			err:    fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			origin: "strain",
			err:    fmt.Errorf("some error"),
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
			}).UpdateSource(
				context.Background(),
				tc.origin,
				tc.s,
				"Test_ChangeSource")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_RemoveSource(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "RemoveSource")

	tcs := map[string]struct {
		db getMockDB
		id types.UUID
		sources,
		result []types.Source
		err error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			sources: []types.Source{
				{UUID: "0"},
				{UUID: "1"},
				{UUID: "2"},
			},
			id: "1",
			result: []types.Source{
				{UUID: "0"},
				{UUID: "2"},
			},
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:  "0",
			err: fmt.Errorf("source could not be deleted: '0'"),
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

			g := types.Generation{Sources: tc.sources}

			err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).RemoveSource(
				context.Background(),
				&g,
				tc.id,
				"Test_RemoveSource")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, g.Sources)
		})
	}
}
