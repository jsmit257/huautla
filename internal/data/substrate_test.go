package data

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"github.com/jsmit257/huautla/types"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/require"
)

var (
	_subs = []Substrate{
		{UUID: "0", Name: "substrate 0", Type: types.PlatingType, Vendor: types.Vendor{UUID: "0", Name: "vendor 0", Website: "website 0"}, Ingredients: []types.Ingredient{}},
		{UUID: "1", Name: "substrate 1", Type: types.GrainType, Vendor: types.Vendor{UUID: "1", Name: "vendor 1", Website: "website 1"}, Ingredients: []types.Ingredient{}},
		{UUID: "2", Name: "substrate 2", Type: types.GrainType, Vendor: types.Vendor{UUID: "1", Name: "vendor 1", Website: "website 1"}, Ingredients: []types.Ingredient{}},
	}
	subFields = row{"id", "name", "type", "vendor_uuid", "vendor_name", "vendor_website"}
	subValues = [][]driver.Value{
		{_subs[0].UUID, _subs[0].Name, _subs[0].Type, _subs[0].Vendor.UUID, _subs[0].Vendor.Name, _subs[0].Vendor.Website},
		{_subs[1].UUID, _subs[1].Name, _subs[1].Type, _subs[1].Vendor.UUID, _subs[1].Vendor.Name, _subs[1].Vendor.Website},
		{_subs[2].UUID, _subs[2].Name, _subs[2].Type, _subs[2].Vendor.UUID, _subs[2].Vendor.Name, _subs[2].Vendor.Website},
	}
)

func Test_SelectAllSubstrates(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "SelectAllSubstrates")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result []types.Substrate
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				subFields.mock(mock, subValues...)
				for range subValues {
					ingFields.mock(mock)
				}

				return db
			},
			result: []types.Substrate{
				types.Substrate(_subs[0]),
				types.Substrate(_subs[1]),
				types.Substrate(_subs[2]),
			},
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
			}).SelectAllSubstrates(context.Background(), "Test_SelectAllSubstrates")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_SelectSubstrate(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "SelectSubstrate")

	tcs := map[string]struct {
		db     getMockDB
		result types.Substrate
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				subFields.mock(mock, subValues[0])
				ingFields.mock(mock, ingValues...)

				return db
			},
			result: types.Substrate{
				UUID:   _subs[0].UUID,
				Name:   _subs[0].Name,
				Type:   _subs[0].Type,
				Vendor: _subs[0].Vendor,
				Ingredients: []types.Ingredient{
					types.Ingredient(_ingredients[0]),
					types.Ingredient(_ingredients[1]),
					types.Ingredient(_ingredients[2]),
				}},
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
			}).SelectSubstrate(context.Background(), "tc.id", "Test_SelectSubstrate")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_InsertSubstrate(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "InsertSubstrate")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		tp     types.SubstrateType
		result types.Substrate
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id:     "0",
			tp:     types.GrainType,
			result: types.Substrate{UUID: "30313233-3435-3637-3839-616263646566", Name: "substrate 0", Type: types.GrainType, Vendor: types.Vendor{}},
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:     "0",
			tp:     types.GrainType,
			result: types.Substrate{UUID: "30313233-3435-3637-3839-616263646566", Name: "substrate 0", Type: types.GrainType, Vendor: types.Vendor{}},
			err:    fmt.Errorf("substrate was not added"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			id:     "0",
			tp:     types.GrainType,
			result: types.Substrate{UUID: "30313233-3435-3637-3839-616263646566", Name: "substrate 0", Type: types.GrainType, Vendor: types.Vendor{}, Ingredients: nil},
			err:    fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			id:     "0",
			tp:     types.GrainType,
			result: types.Substrate{UUID: "30313233-3435-3637-3839-616263646566", Name: "substrate 0", Type: types.GrainType, Vendor: types.Vendor{}, Ingredients: nil},
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
			}).InsertSubstrate(
				context.Background(),
				types.Substrate{Name: "substrate " + string(tc.id), Type: tc.tp, Vendor: types.Vendor{}, Ingredients: nil},
				"Test_InsertSubstrate")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_UpdateSubstrate(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "UpdateSubstrate")

	tcs := map[string]struct {
		db  getMockDB
		id  types.UUID
		err error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id: "0",
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:  "0",
			err: fmt.Errorf("substrate was not updated: '0'"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
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
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).UpdateSubstrate(
				context.Background(),
				tc.id,
				types.Substrate{Name: "substrate " + string(tc.id)},
				"Test_UpdateSubstrates")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_DeleteSubstrate(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "DeleteSubstrate")

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
			err: fmt.Errorf("substrate could not be deleted: '0'"),
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
			}).DeleteSubstrate(
				context.Background(),
				tc.id,
				"Test_DeleteSubstrate")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_SubstrateReport(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "SubstrateReport")

	tcs := map[string]struct {
		db                      getMockDB
		result, other, expected types.Entity
		err                     error
	}{
		"happy_generation_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				subFields.mock(mock, subValues[0])
				ingFields.mock(mock, ingValues...)
				genFields.mock(mock, genValues)
				eventFields.mock(mock, eventValues...)
				srcFields.mock(mock, srcValues[0])
				ingFields.mock(mock, ingValues...)
				ingFields.mock(mock, ingValues...)
				napFields.mock(mock)
				noteFields.mock(mock, noteValues...)
				strainFields.mock(mock)

				return db
			},

			result: func(s types.Entity) types.Entity {
				s["ingredients"] = []interface{}{
					mustObject(_ingredients[0]),
					mustObject(_ingredients[1]),
					mustObject(_ingredients[2]),
				}

				s["generations"] = []types.Entity{mustEntity(_gen)}

				g := s["generations"].([]types.Entity)[0]

				g["events"] = []interface{}{
					map[string]interface{}{"ctime": wwtbn.Format(time.RFC3339Nano), "event_type": map[string]interface{}{"id": "", "name": "", "severity": "", "stage": map[string]interface{}{"id": "", "name": ""}}, "id": "0", "mtime": wwtbn.Format(time.RFC3339Nano), "temperature": 0.0},
					map[string]interface{}{"ctime": wwtbn.Format(time.RFC3339Nano), "event_type": map[string]interface{}{"id": "", "name": "", "severity": "", "stage": map[string]interface{}{"id": "", "name": ""}}, "id": "1", "mtime": wwtbn.Format(time.RFC3339Nano), "temperature": 0.0},
					map[string]interface{}{"ctime": wwtbn.Format(time.RFC3339Nano), "event_type": map[string]interface{}{"id": "", "name": "", "severity": "", "stage": map[string]interface{}{"id": "", "name": ""}}, "id": "2", "mtime": wwtbn.Format(time.RFC3339Nano), "temperature": 0.0},
				}
				g["plating_substrate"].(map[string]interface{})["ingredients"] = []interface{}{
					mustObject(_ingredients[0]),
					mustObject(_ingredients[1]),
					mustObject(_ingredients[2]),
				}
				g["liquid_substrate"].(map[string]interface{})["ingredients"] = []interface{}{
					mustObject(_ingredients[0]),
					mustObject(_ingredients[1]),
					mustObject(_ingredients[2]),
				}
				g["notes"] = []types.Entity{
					mustEntity(_notes[0]),
					mustEntity(_notes[1]),
					mustEntity(_notes[2]),
				}
				g["sources"] = []interface{}{
					map[string]interface{}{
						"id":     _src.UUID,
						"strain": mustEntity(_src.Strain),
						"type":   _src.Type,
					},
				}

				return s
			}(mustEntity(_subs[0])),
		},
		"generation_path_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				subFields.mock(mock, subValues[0])
				ingFields.mock(mock, ingValues...)

				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))

				return db
			},
			err: fmt.Errorf("some error"),
		},
		"happy_lifecycle_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				subFields.mock(mock, subValues[1])
				ingFields.mock(mock, ingValues...)
				lcFields.mock(mock, lcValues)
				eventFields.mock(mock, eventValues...)
				attrFields.mock(mock, attrValues...)
				ingFields.mock(mock, ingValues...)
				ingFields.mock(mock, ingValues...)
				napFields.mock(mock)
				noteFields.mock(mock, noteValues...)
				genFields.mock(mock)

				return db
			},
			result: types.Entity{
				"id": _subs[1].UUID,
				"ingredients": []interface{}{
					mustObject(_ingredients[0]),
					mustObject(_ingredients[1]),
					mustObject(_ingredients[2]),
				},
				"lifecycles": []types.Entity{{
					"bulk_substrate": map[string]interface{}{
						"id": "bs",
						"ingredients": []interface{}{
							mustObject(_ingredients[0]),
							mustObject(_ingredients[1]),
							mustObject(_ingredients[2]),
						},
						"name":   "bs",
						"type":   "bulk",
						"vendor": mustObject(_lc.BulkSubstrate.Vendor),
					},
					"ctime": wwtbn.Format(time.RFC3339Nano),
					"events": []interface{}{
						map[string]interface{}{
							"ctime": wwtbn.Format(time.RFC3339Nano),
							"event_type": map[string]interface{}{
								"id":       "",
								"name":     "",
								"severity": "",
								"stage": map[string]interface{}{
									"id":   "",
									"name": "",
								},
							},
							"id":          "0",
							"mtime":       wwtbn.Format(time.RFC3339Nano),
							"temperature": 0.0,
						},
						map[string]interface{}{
							"ctime": wwtbn.Format(time.RFC3339Nano),
							"event_type": map[string]interface{}{
								"id":       "",
								"name":     "",
								"severity": "",
								"stage": map[string]interface{}{
									"id":   "",
									"name": "",
								},
							},
							"id":          "1",
							"mtime":       wwtbn.Format(time.RFC3339Nano),
							"temperature": 0.0},
						map[string]interface{}{
							"ctime": wwtbn.Format(time.RFC3339Nano),
							"event_type": map[string]interface{}{
								"id":       "",
								"name":     "",
								"severity": "",
								"stage": map[string]interface{}{
									"id":   "",
									"name": "",
								},
							},
							"id":          "2",
							"mtime":       wwtbn.Format(time.RFC3339Nano),
							"temperature": 0.0},
					},
					"grain_substrate": map[string]interface{}{
						"id": "gs",
						"ingredients": []interface{}{
							mustObject(_ingredients[0]),
							mustObject(_ingredients[1]),
							mustObject(_ingredients[2]),
						},
						"name":   "gs",
						"type":   "grain",
						"vendor": mustObject(_lc.GrainSubstrate.Vendor),
					},
					"id":       "30313233-3435-3637-3839-616263646566",
					"location": "location",
					"mtime":    wwtbn.Format(time.RFC3339Nano),
					"notes": []types.Entity{
						mustEntity(_notes[0]),
						mustEntity(_notes[1]),
						mustEntity(_notes[2]),
					},
					"strain": map[string]interface{}{
						"attributes": []interface{}{
							mustObject(_attrs[0]),
							mustObject(_attrs[1]),
							mustObject(_attrs[2]),
						},
						"ctime":   wwtbn.Format(time.RFC3339Nano),
						"id":      "strainuuid 0",
						"name":    "strainname 0",
						"species": "X.species",
						"vendor":  mustObject(_lc.Strain.Vendor),
					},
				}},
				"name":   _subs[1].Name,
				"type":   _subs[1].Type,
				"vendor": mustObject(_subs[1].Vendor),
			},
		},
		"lifecycle_path_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				subFields.mock(mock, subValues[1])
				ingFields.mock(mock, ingValues...)

				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))

				return db
			},
			err: fmt.Errorf("some error"),
		},
		"happy_path_no_children": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				subFields.mock(mock, subValues[0])
				ingFields.mock(mock, ingValues...)
				genFields.mock(mock)

				return db
			},
			result: types.Entity{
				"id":   "0",
				"name": "substrate 0",
				"type": string(types.PlatingType),
				"vendor": map[string]interface{}{
					"id":      "0",
					"name":    "vendor 0",
					"website": "website 0",
				},
				"ingredients": []interface{}{
					map[string]interface{}{"id": "0", "name": "ingredient 0"},
					map[string]interface{}{"id": "1", "name": "ingredient 1"},
					map[string]interface{}{"id": "2", "name": "ingredient 2"},
				}},
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
			}).SubstrateReport(context.Background(), "tc.id", "Test_SubstrateReport")

			js, _ := json.Marshal(result)

			require.Equal(t, tc.err, err)
			require.Equal(t, mustObject(tc.result), mustObject(result), string(js))
		})
	}
}
