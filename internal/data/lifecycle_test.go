package data

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jsmit257/huautla/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var (
	_lc = lifecycle{
		UUID:       "30313233-3435-3637-3839-616263646566",
		Location:   "location",
		StrainCost: 0,
		GrainCost:  0,
		BulkCost:   0,
		Yield:      0,
		Count:      0,
		Gross:      0,
		MTime:      wwtbn,
		CTime:      wwtbn,
		Strain:     types.Strain(_strain),
		GrainSubstrate: types.Substrate{
			UUID: "gs",
			Name: "gs",
			Type: types.GrainType,
			Vendor: types.Vendor{
				UUID:    "1",
				Name:    "vendor 1",
				Website: "website",
			},
		},
		BulkSubstrate: types.Substrate{
			UUID: "bs",
			Name: "bs",
			Type: types.BulkType,
			Vendor: types.Vendor{
				UUID:    "2",
				Name:    "vendor 2",
				Website: "website",
			},
		},
	}
	lcFields = row{
		"uuid",
		"location",
		"straincost",
		"graincost",
		"bulkcost",
		"yield",
		"count",
		"gross",
		"mtime",
		"ctime",
		"strain_uuid",
		"strain_species",
		"strain_name",
		"generation_uuid",
		"strain_ctime",
		"strain_dtime",
		"strain_vendor_uuid",
		"strain_vendor_name",
		"strain_vendor_website",
		"grain_substrate_uuid",
		"grain_substrate_name",
		"grain_substrate_type",
		"grain_vendor_uuid",
		"grain_vendor_name",
		"grain_vendor_website",
		"bulk_substrate_uuid",
		"bulk_substrate_name",
		"bulk_substrate_type",
		"bulk_vendor_uuid",
		"bulk_vendor_name",
		"bulk_vendor_website",
	}
	lcValues = []driver.Value{
		_lc.UUID,
		_lc.Location,
		_lc.StrainCost,
		_lc.GrainCost,
		_lc.BulkCost,
		_lc.Yield,
		_lc.Count,
		_lc.Gross,
		_lc.MTime,
		_lc.CTime,
		_lc.Strain.UUID,
		_lc.Strain.Species,
		_lc.Strain.Name,
		nil,
		_lc.Strain.CTime,
		_lc.Strain.DTime,
		_lc.Strain.Vendor.UUID,
		_lc.Strain.Vendor.Name,
		_lc.Strain.Vendor.Website,
		_lc.GrainSubstrate.UUID,
		_lc.GrainSubstrate.Name,
		_lc.GrainSubstrate.Type,
		_lc.GrainSubstrate.Vendor.UUID,
		_lc.GrainSubstrate.Vendor.Name,
		_lc.GrainSubstrate.Vendor.Website,
		_lc.BulkSubstrate.UUID,
		_lc.BulkSubstrate.Name,
		_lc.BulkSubstrate.Type,
		_lc.BulkSubstrate.Vendor.UUID,
		_lc.BulkSubstrate.Vendor.Name,
		_lc.BulkSubstrate.Vendor.Website,
	}
)

func Test_SelectLifecycleIndex(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "SelectLifecycleIndex")

	tcs := map[string]struct {
		db     getMockDB
		result []types.Lifecycle
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{
							"uuid",
							"location",
							"mtime",
							"ctime",
							"strain_uuid",
							"strain_species",
							"strain_name",
							"strain_ctime",
							"vendor_uuid",
							"vendor_name",
							"vendor_website",
							"event_uuid",
							"temp",
							"humidity",
							"event_mtime",
							"event_ctime",
							"et_uuid",
							"et_name",
							"et_sev",
							"stage_uuid",
							"stage_name"}).
						AddRow(
							"0",
							"happy_path",
							wwtbn,
							wwtbn,
							"strain 0",
							"strain 0",
							"strain 0",
							wwtbn,
							"vendor 0",
							"vendor 0",
							"vendor 0",
							"event 0",
							0,
							0,
							wwtbn,
							wwtbn,
							"type 0",
							"type 0",
							"type 0",
							"stage 0",
							"stage 0",
						).
						AddRow(
							"1",
							"happy_path 2",
							wwtbn,
							wwtbn,
							"strain 0",
							"strain 0",
							"strain 0",
							wwtbn,
							"vendor 0",
							"vendor 0",
							"vendor 0",
							"event 0",
							0,
							0,
							wwtbn,
							wwtbn,
							"type 0",
							"type 0",
							"type 0",
							"stage 0",
							"stage 0",
						).
						AddRow(
							"1",
							"happy_path 2",
							wwtbn,
							wwtbn,
							"strain 0",
							"strain 0",
							"strain 0",
							wwtbn,
							"vendor 0",
							"vendor 0",
							"vendor 0",
							"event 1",
							0,
							0,
							wwtbn,
							wwtbn,
							"type 0",
							"type 0",
							"type 0",
							"stage 0",
							"stage 0",
						))
				return db
			},
			result: []types.Lifecycle{
				{
					UUID:     "0",
					Location: "happy_path",
					MTime:    wwtbn,
					CTime:    wwtbn,
					Strain: types.Strain{
						UUID:    "strain 0",
						Name:    "strain 0",
						Species: "strain 0",
						CTime:   wwtbn,
						Vendor: types.Vendor{
							UUID:    "vendor 0",
							Name:    "vendor 0",
							Website: "vendor 0",
						},
					},
					Events: []types.Event{{
						UUID:  "event 0",
						MTime: wwtbn,
						CTime: wwtbn,
						EventType: types.EventType{
							UUID:     "type 0",
							Name:     "type 0",
							Severity: "type 0",
							Stage: types.Stage{
								UUID: "stage 0",
								Name: "stage 0",
							},
						},
					}},
				},
				{
					UUID:     "1",
					Location: "happy_path 2",
					MTime:    wwtbn,
					CTime:    wwtbn,
					Strain: types.Strain{
						UUID:    "strain 0",
						Name:    "strain 0",
						Species: "strain 0",
						CTime:   wwtbn,
						Vendor: types.Vendor{
							UUID:    "vendor 0",
							Name:    "vendor 0",
							Website: "vendor 0",
						},
					},
					Events: []types.Event{
						{
							UUID:  "event 0",
							MTime: wwtbn,
							CTime: wwtbn,
							EventType: types.EventType{
								UUID:     "type 0",
								Name:     "type 0",
								Severity: "type 0",
								Stage: types.Stage{
									UUID: "stage 0",
									Name: "stage 0",
								},
							},
						},
						{
							UUID:  "event 1",
							MTime: wwtbn,
							CTime: wwtbn,
							EventType: types.EventType{
								UUID:     "type 0",
								Name:     "type 0",
								Severity: "type 0",
								Stage: types.Stage{
									UUID: "stage 0",
									Name: "stage 0",
								},
							},
						},
					},
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

			result, err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).SelectLifecycleIndex(context.Background(), "Test_SelectLifecycleIndex")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_SelectLifecycle(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "SelectLifecycle")

	tcs := map[string]struct {
		db     getMockDB
		noid   bool
		result types.Lifecycle
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				lcFields.mock(mock, lcValues)
				eventFields.mock(mock, eventValues...)

				return db
			},
			result: func(lc lifecycle) types.Lifecycle {
				lc.Events = []types.Event{
					types.Event(_events[0]),
					types.Event(_events[1]),
					types.Event(_events[2]),
				}

				return types.Lifecycle(lc)
			}(_lc),
		},
		"too_much_of_a_good_thing": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				lcFields.mock(mock, lcValues, lcValues)
				eventFields.mock(mock, eventValues...)
				eventFields.mock(mock, eventValues...)

				return db
			},
			err: fmt.Errorf("too many rows returned for SelectLifecycle"),
		},
		"get_events_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock,
					lcFields.set(lcValues),
					eventFields.fail())

				return db
			},
			err: eventFields.err(),
		},
		"no_rows": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				lcFields.mock(mock)
				return db
			},
			err: sql.ErrNoRows,
		},
		"missing_id": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				return db
			},
			noid: true,
			err:  fmt.Errorf("request doesn't contain at least 1 required field"),
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

			var id types.UUID
			if !tc.noid {
				id = "0"
			}

			result, err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).SelectLifecycle(context.Background(), id, "Test_SelectLifecycle")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_InsertLifecycle(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "InsertLifecycle")

	tcs := map[string]struct {
		db  getMockDB
		err error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				newBuilder(mock, lcFields.set(lcValues), eventFields.set(eventValues...))

				return db
			},
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("lifecycle was not added: 0"),
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

			lc, err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).InsertLifecycle(
				context.Background(),
				types.Lifecycle{},
				"Test_InsertLifecycle")

			require.Equal(t, tc.err, err)
			require.Equal(t, lc.UUID, types.UUID(mockUUIDGen().String()))
		})
	}
}

func Test_UpdateLifecycle(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "UpdateLifecycle")

	tcs := map[string]struct {
		db  getMockDB
		err error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("lifecycle was not updated"),
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

			// there's no good way to test the returned lifecycle, to start, the
			// timestamps are non-deterministic; system tests will vet the rest
			_, err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).UpdateLifecycle(context.Background(), types.Lifecycle{}, "Test_UpdateLifecycle")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_UpdateLifecycleMTime(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "UpdateLifecycleMTime")

	now := time.Now()

	tcs := map[string]struct {
		db       getMockDB
		modified time.Time
		err      error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			modified: now,
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("mtime was not updated"),
		},
		"row_error": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			err: fmt.Errorf("some error"),
		},
		"db_error": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
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
			}).UpdateLifecycleMTime(context.Background(), &types.Lifecycle{}, time.Now(), "Test_UpdateLifecycle")

			require.Equal(t, tc.err, err)
			// require.Equal(t, tc.modified, lc.MTime)
		})
	}
}

func Test_DeleteLifecycle(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "DeleteLifecycle")

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
			err: fmt.Errorf("lifecycle could not be deleted: '0'"),
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
			}).DeleteLifecycle(
				context.Background(),
				tc.id,
				"Test_DeleteLifecycle")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_LifecycleReport(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "LifecycleReport")

	tcs := map[string]struct {
		db     getMockDB
		noid   bool
		result types.Entity
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock,
					lcFields.set(lcValues),
					eventFields.set(eventValues...),
					attrFields.set(attrValues...),
					ingFields.set(ingValues...),
					ingFields.set(ingValues...),
					napFields.set(),
					noteFields.set(noteValues...),
					photoFields.set())

				return db
			},
			result: func(lc types.Entity) types.Entity {
				lc["strain"].(map[string]interface{})["attributes"] = attributes
				lc["grain_substrate"].(map[string]interface{})["ingredients"] = ingredients
				lc["bulk_substrate"].(map[string]interface{})["ingredients"] = ingredients
				lc["events"] = events
				lc["notes"] = notes

				return lc
			}(mustEntity(_lc)),
		},
		"happy_photo_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock,
					lcFields.set(lcValues),
					eventFields.set(eventValues...),
					attrFields.set(attrValues...),
					ingFields.set(ingValues...),
					ingFields.set(ingValues...),
					napFields.set(),
					noteFields.set(noteValues...),
					photoFields.set(photoValues...))

				return db
			},
			result: func(lc types.Entity) types.Entity {
				strain := lc["strain"].(map[string]interface{})
				strain["attributes"] = attributes
				strain["photos"] = album

				lc["grain_substrate"].(map[string]interface{})["ingredients"] = ingredients
				lc["bulk_substrate"].(map[string]interface{})["ingredients"] = ingredients
				lc["events"] = events
				lc["notes"] = notes

				return lc
			}(mustEntity(_lc)),
		},
		"photo_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock,
					lcFields.set(lcValues),
					eventFields.set(eventValues...),
					attrFields.set(attrValues...),
					ingFields.set(ingValues...),
					ingFields.set(ingValues...),
					napFields.set(),
					noteFields.set(noteValues...),
					photoFields.fail())

				return db
			},
			err: photoFields.err(),
		},
		"notes_fail": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock,
					lcFields.set(lcValues),
					eventFields.set(),
					attrFields.set(),
					ingFields.set(),
					ingFields.set(),
					napFields.set(),
					noteFields.fail())

				return db
			},
			err: noteFields.err(),
		},
		"get_events_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock, lcFields.set(lcValues), eventFields.fail())
				return db
			},
			err: eventFields.err(),
		},
		"notes_and_photos_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock,
					lcFields.set(lcValues),
					eventFields.set(),
					attrFields.set(),
					ingFields.set(),
					ingFields.set(),
					napFields.fail())

				return db
			},
			err: napFields.err(),
		},
		"all_attrs_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock,
					lcFields.set(lcValues),
					eventFields.set(),
					attrFields.fail())

				return db
			},
			err: attrFields.err(),
		},
		"get_grain_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock,
					lcFields.set(lcValues),
					eventFields.set(),
					attrFields.set(),
					ingFields.fail())

				return db
			},
			err: ingFields.err(),
		},
		"get_bulk_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock,
					lcFields.set(lcValues),
					eventFields.set(),
					attrFields.set(),
					ingFields.set(),
					ingFields.fail())

				return db
			},
			err: ingFields.err(),
		},
		"no_rows": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock, lcFields.set())
				return db
			},
			err: sql.ErrNoRows,
		},
		"no_id": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				lcFields.mock(mock)
				return db
			},
			noid: true,
			err:  fmt.Errorf("failed to find param values in the following fields: [lifecycle-id]"),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock, lcFields.fail())
				return db
			},
			err: lcFields.err(),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var id types.UUID
			if !tc.noid {
				id = "0"
			}

			result, err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).LifecycleReport(context.Background(), id, "Test_LifecycleReport")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}
