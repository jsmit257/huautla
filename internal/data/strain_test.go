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
	_strain = strain{
		UUID:    "strainuuid 0",
		Name:    "strainname 0",
		Species: "X.species",
		CTime:   wwtbn,
		Vendor: types.Vendor{
			UUID:    "vendoruuid 0",
			Name:    "vendorname 0",
			Website: "vendorwebsite 0",
		},
	}
	strainFields = row{
		"uuid",
		"species",
		"name",
		"ctime",
		"dtime",
		"vendor_uuid",
		"vendor_name",
		"vendor_website",
		"generation_uuid",
	}
	strainValues = []driver.Value{
		_strain.UUID,
		_strain.Species,
		_strain.Name,
		_strain.CTime,
		_strain.DTime,
		_strain.Vendor.UUID,
		_strain.Vendor.Name,
		_strain.Vendor.Website,
		nil,
	}
)

func Test_SelectAllStrains(t *testing.T) {
	t.Parallel()

	whenwillthenbenow := time.Now() // time.Soon()

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

				strainFields.mock(mock,
					[]driver.Value{"0", "X.species", "strain 0", whenwillthenbenow, nil, "0", "vendor 0", "website", nil},
					[]driver.Value{"1", "X.species", "strain 1", whenwillthenbenow, nil, "1", "vendor 1", "website", nil},
					[]driver.Value{"2", "X.species", "strain 2", whenwillthenbenow, nil, "1", "vendor 1", "website", "0"})

				return db
			},
			result: []types.Strain{
				{UUID: "0", Species: "X.species", Name: "strain 0", CTime: whenwillthenbenow, Vendor: types.Vendor{UUID: "0", Name: "vendor 0", Website: "website"}, Attributes: nil},
				{UUID: "1", Species: "X.species", Name: "strain 1", CTime: whenwillthenbenow, Vendor: types.Vendor{UUID: "1", Name: "vendor 1", Website: "website"}, Attributes: nil},
				{UUID: "2", Species: "X.species", Name: "strain 2", CTime: whenwillthenbenow, Vendor: types.Vendor{UUID: "1", Name: "vendor 1", Website: "website"}, Attributes: nil, Generation: &types.Generation{UUID: "0"}},
			},
		},
		// "scan_error": {
		// 	db: func() *sql.DB {
		// 		db, mock, _ := sqlmock.New()
		// 		mock.ExpectQuery("").
		// 			WillReturnRows(sqlmock.
		// 				NewRows([]string{"id", "species", "name", "ctime", "vendor_uuid", "vendor_name", "vendor_website", "generation_uuid"}).
		// 				AddRow("0", "X.species", "strain 0", whenwillthenbenow, "0", "vendor 0", "website", nil).
		// 				RowError(0, fmt.Errorf("some error")))
		// 		return db
		// 	},
		// 	err: fmt.Errorf("some error"),
		// },
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
		result types.Strain
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				strainFields.mock(mock, strainValues)
				attrFields.mock(mock, attrValues...)

				return db
			},
			result: types.Strain{
				UUID:    "strainuuid 0",
				Species: "X.species",
				Name:    "strainname 0",
				Vendor: types.Vendor{
					UUID:    "vendoruuid 0",
					Name:    "vendorname 0",
					Website: "vendorwebsite 0",
				},
				Attributes: []types.StrainAttribute{
					{
						UUID:  "attruuid 0",
						Name:  "attrname 0",
						Value: "attrvalue 0",
					},
					{
						UUID:  "attruuid 1",
						Name:  "attrname 1",
						Value: "attrvalue 1",
					},
					{
						UUID:  "attruuid 2",
						Name:  "attrname 2",
						Value: "attrvalue 2",
					},
				},
				CTime: wwtbn,
			},
		},
		"no_results_found": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				strainFields.mock(mock)
				return db
			},
			result: types.Strain{},
			err:    fmt.Errorf("sql: no rows in result set"),
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
			}).SelectStrain(context.Background(), "tc.id", "Test_SelectStrain")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_InsertStrain(t *testing.T) {
	t.Parallel()

	whenwillthenbenow := time.Now() // time.Soon()

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
			result: types.Strain{UUID: "30313233-3435-3637-3839-616263646566", Name: "strain 0", CTime: whenwillthenbenow, Vendor: types.Vendor{}, Attributes: nil},
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

			_, err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).InsertStrain(
				context.Background(),
				types.Strain{UUID: tc.id, Name: "strain " + string(tc.id), Vendor: types.Vendor{}, Attributes: nil},
				"Test_InsertStrains")

			require.Equal(t, tc.err, err)
			// require.Equal(t, tc.result, result)  // TODO: tripped up by `time` again
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

func Test_GeneratedStrain(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "GeneratedStrain")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result types.Strain
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				strainFields[0:8].mock(mock, strainValues[0:8])
				return db
			},
			id: "0",
			result: types.Strain{
				UUID:    "strainuuid 0",
				Species: "X.species",
				Name:    "strainname 0",
				Vendor: types.Vendor{
					UUID:    "vendoruuid 0",
					Name:    "vendorname 0",
					Website: "vendorwebsite 0",
				},
				CTime: wwtbn,
			},
		},
		"no_results_found": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				strainFields.mock(mock)
				return db
			},
			id:     "0",
			result: types.Strain{},
			err:    sql.ErrNoRows,
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
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
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).GeneratedStrain(context.Background(), tc.id, "Test_GeneratedStrain")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_UpdateGeneratedStrain(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "UpdateGeneratedStrain")

	tcs := map[string]struct {
		db getMockDB
		gid,
		sid types.UUID
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
			gid: "0",
			sid: "0",
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			gid: "0",
			sid: "0",
			err: sql.ErrNoRows,
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			gid: "0",
			sid: "0",
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
			gid: "0",
			sid: "0",
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
			}).UpdateGeneratedStrain(
				context.Background(),
				&tc.gid,
				tc.sid,
				"Test_UpdateStrains")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_StrainReport(t *testing.T) {
	t.Parallel()

	var (
		attrs = []interface{}{
			mustObject(_attrs[0]),
			mustObject(_attrs[1]),
			mustObject(_attrs[2]),
		}
		album = []types.Entity{
			mustObject(_photos[0]),
			mustObject(_photos[1]),
			mustObject(_photos[2]),
		}
		ings = []interface{}{
			mustObject(_ingredients[0]),
			mustObject(_ingredients[1]),
			mustObject(_ingredients[2]),
		}

		l = log.WithField("test", "StrainReport")
	)

	tcs := map[string]struct {
		db                       getMockDB
		result, expected, actual types.Entity
		err                      error
	}{
		"happy_path_with_photos": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				strainFields.mock(mock, strainValues)
				attrFields.mock(mock, attrValues...)
				genFields.mock(mock)
				lcFields.mock(mock)
				photoFields.mock(mock, photoValues...)

				return db
			},
			result: func(s types.Entity) types.Entity {
				s["attributes"] = attrs
				s["photos"] = album
				return s
			}(mustEntity(_strain)),
		},
		"photos_report_error": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				strainFields.mock(mock, strainValues)
				attrFields.mock(mock, attrValues...)
				genFields.mock(mock)
				lcFields.mock(mock)

				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))

				return db
			},
			err: fmt.Errorf("some error"),
		},
		"happy_lifecycles_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				strainFields.mock(mock, strainValues)
				attrFields.mock(mock, attrValues...)
				genFields.mock(mock)
				lcFields.mock(mock, lcValues)
				eventFields.mock(mock, eventValues...)
				attrFields.mock(mock, attrValues...)
				ingFields.mock(mock, ingValues...)
				ingFields.mock(mock, ingValues...)
				napFields.mock(mock)
				noteFields.mock(mock, noteValues...)
				photoFields.mock(mock)
				photoFields.mock(mock)

				return db
			},
			result: func(s types.Entity) types.Entity {
				s["attributes"] = attrs

				s["lifecycles"] = []types.Entity{mustEntity(_lc)}

				lc := s["lifecycles"].([]types.Entity)[0]

				lc["strain"].(map[string]interface{})["attributes"] = []interface{}{
					mustObject(_attrs[0]),
					mustObject(_attrs[1]),
					mustObject(_attrs[2]),
				}

				lc["grain_substrate"].(map[string]interface{})["ingredients"] = ings

				lc["bulk_substrate"].(map[string]interface{})["ingredients"] = ings

				lc["events"] = []interface{}{
					mustObject(_events[0]),
					mustObject(_events[1]),
					mustObject(_events[2]),
				}

				lc["notes"] = []types.Entity{
					mustEntity(_notes[0]),
					mustEntity(_notes[1]),
					mustEntity(_notes[2]),
				}

				return s
			}(mustEntity(_strain)),
		},
		"lifecycles_report_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				strainFields.mock(mock, strainValues)
				attrFields.mock(mock, attrValues...)
				genFields.mock(mock)

				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))

				return db
			},
			err: fmt.Errorf("some error"),
		},
		"happy_generations_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				strainFields.mock(mock, strainValues)
				attrFields.mock(mock, attrValues...)
				genFields.mock(mock, genValues)
				eventFields.mock(mock, eventValues...)
				srcFields.mock(mock, srcValues...)
				ingFields.mock(mock, ingValues...)
				ingFields.mock(mock, ingValues...)
				napFields.mock(mock)
				noteFields.mock(mock, noteValues...)
				strainFields.mock(mock)
				lcFields.mock(mock)
				photoFields.mock(mock)

				return db
			},
			result: func(s types.Entity) types.Entity {
				s["attributes"] = attrs

				s["generations"] = []types.Entity{mustEntity(_gen)}

				gen := s["generations"].([]types.Entity)[0]

				gen["plating_substrate"].(map[string]interface{})["ingredients"] = ings

				gen["liquid_substrate"].(map[string]interface{})["ingredients"] = ings

				gen["events"] = []interface{}{
					mustObject(_events[0]),
					mustObject(_events[1]),
					mustObject(_events[2]),
				}

				gen["notes"] = []types.Entity{
					mustEntity(_notes[0]),
					mustEntity(_notes[1]),
					mustEntity(_notes[2]),
				}

				gen["sources"] = []types.Entity{
					mustEntity(types.Source{
						UUID: _src.UUID,
						Type: _src.Type,
						// Lifecycle: nil,
						Strain: types.Strain(_strain),
					}),
				}

				return s
			}(mustEntity(_strain)),
			expected: map[string]interface{}{
				"attributes": []interface{}{
					map[string]interface{}{"id": "attruuid 0", "name": "attrname 0", "value": "attrvalue 0"},
					map[string]interface{}{"id": "attruuid 1", "name": "attrname 1", "value": "attrvalue 1"},
					map[string]interface{}{"id": "attruuid 2", "name": "attrname 2", "value": "attrvalue 2"},
				},
				"ctime": wwtbn.Format(time.RFC3339Nano),
				"generation": map[string]interface{}{
					"ctime": wwtbn.Format(time.RFC3339Nano),
					"events": []interface{}{
						map[string]interface{}{"ctime": wwtbn.Format(time.RFC3339Nano), "event_type": map[string]interface{}{"id": "", "name": "", "severity": "", "stage": map[string]interface{}{"id": "", "name": ""}}, "id": "0", "mtime": wwtbn.Format(time.RFC3339Nano), "temperature": 0},
						map[string]interface{}{"ctime": wwtbn.Format(time.RFC3339Nano), "event_type": map[string]interface{}{"id": "", "name": "", "severity": "", "stage": map[string]interface{}{"id": "", "name": ""}}, "id": "1", "mtime": wwtbn.Format(time.RFC3339Nano), "temperature": 0},
						map[string]interface{}{"ctime": wwtbn.Format(time.RFC3339Nano), "event_type": map[string]interface{}{"id": "", "name": "", "severity": "", "stage": map[string]interface{}{"id": "", "name": ""}}, "id": "2", "mtime": wwtbn.Format(time.RFC3339Nano), "temperature": 0},
					},
					"id": "uuid",
					"liquid_substrate": map[string]interface{}{
						"id": "liquidsubstrate_uuid",
						"ingredients": []interface{}{
							map[string]interface{}{"id": "0", "name": "ingredient 0"},
							map[string]interface{}{"id": "1", "name": "ingredient 1"},
							map[string]interface{}{"id": "2", "name": "ingredient 2"},
						},
						"name":   "liquidsubstrate_name",
						"type":   "liquidsubstrate_type",
						"vendor": map[string]interface{}{"id": "liquidsubstrate_vendor_uuid", "name": "liquidsubstrate_vendor_name", "website": "liquidsubstrate_vendor_website"},
					},
					"mtime": wwtbn.Format(time.RFC3339Nano),
					"notes": []interface{}{
						map[string]interface{}{"ctime": wwtbn.Format(time.RFC3339Nano), "id": "noteuuid 0", "mtime": wwtbn.Format(time.RFC3339Nano), "note": "notenote 0"},
						map[string]interface{}{"ctime": wwtbn.Format(time.RFC3339Nano), "id": "noteuuid 1", "mtime": wwtbn.Format(time.RFC3339Nano), "note": "notenote 1"},
						map[string]interface{}{"ctime": wwtbn.Format(time.RFC3339Nano), "id": "noteuuid 2", "mtime": wwtbn.Format(time.RFC3339Nano), "note": "notenote 2"},
					},
					"plating_substrate": map[string]interface{}{
						"id": "platingsubstrate_uuid",
						"ingredients": []interface{}{
							map[string]interface{}{"id": "0", "name": "ingredient 0"},
							map[string]interface{}{"id": "1", "name": "ingredient 1"},
							map[string]interface{}{"id": "2", "name": "ingredient 2"},
						},
						"name":   "platingsubstrate_name",
						"type":   "platingsubstrate_type",
						"vendor": map[string]interface{}{"id": "platingsubstrate_vendor_uuid", "name": "platingsubstrate_vendor_name", "website": "platingsubstrate_vendor_website"},
					},
					"sources": []interface{}{
						map[string]interface{}{
							"id": "sourceuuid",
							"strain": map[string]interface{}{
								"ctime":   wwtbn.Format(time.RFC3339Nano),
								"id":      "strainuuid 0",
								"name":    "strainname 0",
								"species": "X.species",
								"vendor":  map[string]interface{}{"id": "vendoruuid 0", "name": "vendorname 0", "website": "vendorwebsite 0"}},
							"type": "sourcetype"},
					},
				},
				"id":      "strainuuid 0",
				"name":    "strainname 0",
				"species": "X.species",
				"vendor":  map[string]interface{}{"id": "vendoruuid 0", "name": "vendorname 0", "website": "vendorwebsite 0"},
			},
			actual: map[string]interface{}{
				"attributes": []interface{}{
					map[string]interface{}{"id": "attruuid 0", "name": "attrname 0", "value": "attrvalue 0"},
					map[string]interface{}{"id": "attruuid 1", "name": "attrname 1", "value": "attrvalue 1"},
					map[string]interface{}{"id": "attruuid 2", "name": "attrname 2", "value": "attrvalue 2"},
				},
				"ctime": wwtbn.Format(time.RFC3339Nano),
				"generations": []interface{}{
					map[string]interface{}{
						"ctime": wwtbn.Format(time.RFC3339Nano),
						"events": []interface{}{
							map[string]interface{}{"ctime": wwtbn.Format(time.RFC3339Nano), "event_type": map[string]interface{}{"id": "", "name": "", "severity": "", "stage": map[string]interface{}{"id": "", "name": ""}}, "id": "0", "mtime": wwtbn.Format(time.RFC3339Nano), "temperature": 0},
							map[string]interface{}{"ctime": wwtbn.Format(time.RFC3339Nano), "event_type": map[string]interface{}{"id": "", "name": "", "severity": "", "stage": map[string]interface{}{"id": "", "name": ""}}, "id": "1", "mtime": wwtbn.Format(time.RFC3339Nano), "temperature": 0},
							map[string]interface{}{"ctime": wwtbn.Format(time.RFC3339Nano), "event_type": map[string]interface{}{"id": "", "name": "", "severity": "", "stage": map[string]interface{}{"id": "", "name": ""}}, "id": "2", "mtime": wwtbn.Format(time.RFC3339Nano), "temperature": 0},
						},
						"id": "uuid",
						"liquid_substrate": map[string]interface{}{
							"id": "liquidsubstrate_uuid",
							"ingredients": []interface{}{
								map[string]interface{}{"id": "0", "name": "ingredient 0"},
								map[string]interface{}{"id": "1", "name": "ingredient 1"},
								map[string]interface{}{"id": "2", "name": "ingredient 2"},
							},
							"name":   "liquidsubstrate_name",
							"type":   "liquidsubstrate_type",
							"vendor": map[string]interface{}{"id": "liquidsubstrate_vendor_uuid", "name": "liquidsubstrate_vendor_name", "website": "liquidsubstrate_vendor_website"},
						},
						"mtime": wwtbn.Format(time.RFC3339Nano),
						"notes": []interface{}{
							map[string]interface{}{"ctime": wwtbn.Format(time.RFC3339Nano), "id": "noteuuid 0", "mtime": wwtbn.Format(time.RFC3339Nano), "note": "notenote 0"},
							map[string]interface{}{"ctime": wwtbn.Format(time.RFC3339Nano), "id": "noteuuid 1", "mtime": wwtbn.Format(time.RFC3339Nano), "note": "notenote 1"},
							map[string]interface{}{"ctime": wwtbn.Format(time.RFC3339Nano), "id": "noteuuid 2", "mtime": wwtbn.Format(time.RFC3339Nano), "note": "notenote 2"},
						},
						"plating_substrate": map[string]interface{}{
							"id": "platingsubstrate_uuid",
							"ingredients": []interface{}{
								map[string]interface{}{"id": "0", "name": "ingredient 0"},
								map[string]interface{}{"id": "1", "name": "ingredient 1"},
								map[string]interface{}{"id": "2", "name": "ingredient 2"},
							},
							"name":   "platingsubstrate_name",
							"type":   "platingsubstrate_type",
							"vendor": map[string]interface{}{"id": "platingsubstrate_vendor_uuid", "name": "platingsubstrate_vendor_name", "website": "platingsubstrate_vendor_website"},
						},
						"sources": []interface{}{
							map[string]interface{}{
								"id": "sourceuuid",
								"strain": map[string]interface{}{
									"ctime":   wwtbn.Format(time.RFC3339Nano),
									"id":      "strainuuid 0",
									"name":    "strainname 0",
									"species": "X.species",
									"vendor":  map[string]interface{}{"id": "vendoruuid 0", "name": "vendorname 0", "website": "vendorwebsite 0"},
								},
								"type": "sourcetype"},
						},
					},
				},
				"id":      "strainuuid 0",
				"name":    "strainname 0",
				"species": "X.species",
				"vendor":  map[string]interface{}{"id": "vendoruuid 0", "name": "vendorname 0", "website": "vendorwebsite 0"},
			},
		},
		"generation_report_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				strainFields.mock(mock, strainValues)
				attrFields.mock(mock, attrValues...)

				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))

				return db
			},
			err: fmt.Errorf("some error"),
		},
		"happy_progenitor_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				strainFields.mock(mock, xformer(strainValues).replace(xform{8: "not-nil"}))
				attrFields.mock(mock, attrValues...)
				genFields.mock(mock)
				lcFields.mock(mock)
				photoFields.mock(mock)
				// BEGIN: progenitor
				genFields.mock(mock, genValues)
				eventFields.mock(mock, eventValues...)
				srcFields.mock(mock, srcValues...)
				ingFields.mock(mock, ingValues...)
				ingFields.mock(mock, ingValues...)
				napFields.mock(mock)
				noteFields.mock(mock, noteValues...)
				strainFields.mock(mock)

				return db
			},
			result: func(s types.Entity) types.Entity {
				s["attributes"] = attrs

				s["generation"] = mustEntity(_gen)

				gen := s["generation"].(types.Entity)

				gen["plating_substrate"].(map[string]interface{})["ingredients"] = ings

				gen["liquid_substrate"].(map[string]interface{})["ingredients"] = ings

				gen["events"] = []interface{}{
					mustObject(_events[0]),
					mustObject(_events[1]),
					mustObject(_events[2]),
				}

				gen["notes"] = []types.Entity{
					mustEntity(_notes[0]),
					mustEntity(_notes[1]),
					mustEntity(_notes[2]),
				}

				gen["sources"] = []types.Entity{
					mustEntity(types.Source{
						UUID: _src.UUID,
						Type: _src.Type,
						// Lifecycle: nil,
						Strain: types.Strain(_strain),
					}),
				}

				return s
			}(mustEntity(_strain)),
		},
		"missing_progen_id": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				strainFields.mock(mock, xformer(strainValues).replace(xform{8: ""}))
				attrFields.mock(mock, attrValues...)
				genFields.mock(mock)
				lcFields.mock(mock)
				photoFields.mock(mock)

				return db
			},
			err: fmt.Errorf("failed to find param values in the following fields: [generation-id]"),
		},
		"progen_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				strainFields.mock(mock, xformer(strainValues).replace(xform{8: "not-nil"}))
				attrFields.mock(mock, attrValues...)
				genFields.mock(mock)
				lcFields.mock(mock)
				photoFields.mock(mock)

				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))

				return db
			},
			err: fmt.Errorf("some error"),
		},
		"progen_empty": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				strainFields.mock(mock, xformer(strainValues).replace(xform{8: "not-nil"}))
				attrFields.mock(mock, attrValues...)
				genFields.mock(mock)
				lcFields.mock(mock)
				photoFields.mock(mock)
				genFields.mock(mock)

				return db
			},
			err: fmt.Errorf("how does 'not-nil' not identify a generation?"),
		},
		"happy_path_no_children": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				strainFields.mock(mock, strainValues)
				attrFields.mock(mock, attrValues...)
				genFields.mock(mock)
				lcFields.mock(mock)
				photoFields.mock(mock)

				return db
			},
			result: func() types.Entity {
				s := mustEntity(_strain)
				s["attributes"] = attrs
				return s
			}(),
		},
		"no_results_found": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				strainFields.mock(mock)
				return db
			},
			err: fmt.Errorf("sql: no rows in result set"),
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
			}).StrainReport(context.Background(), "tc.id", "Test_StrainReport")

			js1 := mustEntity(result)
			js2 := mustEntity(tc.result)
			require.Equal(t, tc.err, err)
			// assert.Equal(t, js2, js1)
			require.Equal(t, mustObject(tc.result), mustObject(result), "\n%s\n%s", js1, js2)
		})
	}
}
