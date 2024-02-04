package tests

//package tests ??

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/internal/data"
	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

func main() {}

var db types.DB
var noRows error = fmt.Errorf("sql: no rows in result set")

func init() {
	var err error
	if db, err = data.New(&types.Config{}, nil); err != nil {
		panic(err)
	}
}

func Test_Eventer(t *testing.T) {
	tcs := map[string]func(t *testing.T){
		"GetLifecycleEvents": func(t *testing.T) {
			gle := map[string]struct {
				fn  func(context.Context, *types.Lifecycle, types.CID) error
				lc  types.Lifecycle
				err error
			}{
				"happy_path": {
					fn: db.GetLifecycleEvents,
					lc: types.Lifecycle{},
				},
				"no_events_to_get": {
					lc: types.Lifecycle{},
				},
			}
			for k, v := range gle {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), &v.lc, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
		"SelectByEventType": func(t *testing.T) {
			set := map[string]struct {
				fn     func(context.Context, types.EventType, types.CID) ([]types.Event, error)
				e      types.EventType
				result []types.Event
				err    error
			}{
				"happy_path": {
					fn: db.SelectByEventType,
				},
			}
			for k, v := range set {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					result, err := v.fn(context.Background(), v.e, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"SelectEvent": func(t *testing.T) {
			set := map[string]struct {
				fn     func(context.Context, types.UUID, types.CID) (types.Event, error)
				id     types.UUID
				err    error
				result types.Event
			}{
				"happy_path": {
					fn: db.SelectEvent,
				},
			}
			for k, v := range set {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					result, err := v.fn(context.Background(), v.id, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"AddEvent": func(t *testing.T) {
			set := map[string]struct {
				fn  func(context.Context, *types.Lifecycle, types.Event, types.CID) error
				lc  types.Lifecycle
				e   types.Event
				err error
			}{
				"happy_path": {
					fn: db.AddEvent,
				},
			}
			for k, v := range set {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), &v.lc, v.e, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
		"ChangeEvent": func(t *testing.T) {
			set := map[string]struct {
				fn  func(context.Context, *types.Lifecycle, types.Event, types.CID) error
				lc  types.Lifecycle
				e   types.Event
				err error
			}{
				"happy_path": {
					fn: db.ChangeEvent,
				},
			}
			for k, v := range set {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), &v.lc, v.e, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
		"RemoveEvent": func(t *testing.T) {
			set := map[string]struct {
				fn  func(context.Context, *types.Lifecycle, types.UUID, types.CID) error
				lc  types.Lifecycle
				id  types.UUID
				err error
			}{
				"happy_path": {
					fn: db.RemoveEvent,
				},
			}
			for k, v := range set {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), &v.lc, v.id, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
	}
	for k, v := range tcs {
		// k, v := k, v
		t.Run(k, v)
	}
}

func Test_EventTyper(t *testing.T) {
	tcs := map[string]func(t *testing.T){
		"SelectAllEventTypes": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, cid types.CID) ([]types.EventType, error)
				result []types.EventType
				err    error
			}{
				"happy_path": {
					fn: db.SelectAllEventTypes,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					result, err := v.fn(context.Background(), types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		}, //
		"SelectEventType": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, id types.UUID, cid types.CID) (types.EventType, error)
				id     types.UUID
				result types.EventType
				err    error
			}{
				"happy_path": {
					fn: db.SelectEventType,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					result, err := v.fn(context.Background(), v.id, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		}, //
		"InsertEventType": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, e types.EventType, cid types.CID) (types.EventType, error)
				e      types.EventType
				result types.EventType
				err    error
			}{
				"happy_path": {
					fn: db.InsertEventType,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					result, err := v.fn(context.Background(), v.e, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"UpdateEventType": func(t *testing.T) {
			set := map[string]struct {
				fn  func(ctx context.Context, id types.UUID, e types.EventType, cid types.CID) error
				id  types.UUID
				e   types.EventType
				err error
			}{
				"happy_path": {
					fn: db.UpdateEventType,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), v.id, v.e, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
		"DeleteEventType": func(t *testing.T) {
			set := map[string]struct {
				fn  func(ctx context.Context, id types.UUID, cid types.CID) error
				id  types.UUID
				err error
			}{
				"happy_path": {
					fn: db.DeleteEventType,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), v.id, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
	}
	for k, v := range tcs {
		t.Run(k, v)
	}
}

func Test_Ingredienter(t *testing.T) {
	t.Parallel()
	tcs := map[string]func(t *testing.T){
		"SelectAllIngredients": func(t *testing.T) {
			set := map[string]struct {
				result []types.Ingredient
				err    error
			}{
				"happy_path": {
					result: []types.Ingredient{
						{UUID: "0", Name: "Vermiculite"},
						{UUID: "1", Name: "Maltodexterin"},
						{UUID: "2", Name: "Rye"},
						{UUID: "3", Name: "White Millet"},
						{UUID: "4", Name: "Popcorn"},
						{UUID: "5", Name: "Manure"},
						{UUID: "6", Name: "Coir"},
						{UUID: "7", Name: "Honey"},
						{UUID: "8", Name: "Agar"},
						{UUID: "9", Name: "Rice Flour"},
						{UUID: "10", Name: "White Milo"},
						{UUID: "11", Name: "Red Milo"},
						{UUID: "12", Name: "Red Millet"},
						{UUID: "13", Name: "Gypsum"},
						{UUID: "14", Name: "Calcium phosphate"},
						{UUID: "15", Name: "Diammonium phosphate"},
					},
				},
			}
			for k, v := range set {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					t.Parallel()
					result, err := db.SelectAllIngredients(context.Background(), types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"SelectIngredient": func(t *testing.T) {
			set := map[string]struct {
				id     types.UUID
				result types.Ingredient
				err    error
			}{
				"happy_path": {
					id:     "13",
					result: types.Ingredient{UUID: "13", Name: "Gypsum"},
				},
				"no_rows_returned": {
					id:     "foobar",
					result: types.Ingredient{UUID: "foobar", Name: ""},
					err:    fmt.Errorf("sql: no rows in result set"),
				},
			}
			for k, v := range set {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					t.Parallel()
					result, err := db.SelectIngredient(context.Background(), v.id, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"InsertIngredient": func(t *testing.T) {
			set := map[string]struct {
				i      types.Ingredient
				result []types.Ingredient
				err    error
			}{
				"happy_path": {
					i: types.Ingredient{Name: "bogus!"},
				},
				"no_rows_affected": {}, // ???
				"duplicate_name_violation": {
					i:      types.Ingredient{Name: "Coir"},
					result: []types.Ingredient{},
					err:    fmt.Errorf(""),
				},
			}
			for k, v := range set {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					t.Parallel()
					result, err := db.InsertIngredient(context.Background(), v.i, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"UpdateIngredient": func(t *testing.T) {
			set := map[string]struct {
				id  types.UUID
				i   types.Ingredient
				err error
			}{
				"happy_path": {
					id: "0",
					i:  types.Ingredient{Name: "renamed"},
				},
				"no_rows_affected": {
					id:  "foobar",
					err: fmt.Errorf("ingredient was not updated: '0'"),
				},
				"duplicate_name_violation": {
					id: "0",
					i:  types.Ingredient{Name: "Honey"},
				},
			}
			for k, v := range set {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					t.Parallel()
					err := db.UpdateIngredient(context.Background(), v.id, v.i, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
		"DeleteIngredient": func(t *testing.T) {
			set := map[string]struct {
				id  types.UUID
				err error
			}{
				"happy_path": {
					id: "-1",
				},
				"no_rows_affected": {
					id:  "foobar",
					err: fmt.Errorf("ingredient was not deleted: '0'"),
				},
				"query_fails": {
					id:  "01234567890123456789012345678901234567891",
					err: fmt.Errorf("ingredient was not deleted: '0'"),
				},
				"referential_violation": {
					id:  "2",
					err: fmt.Errorf("referential constraint"),
				},
			}
			for k, v := range set {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					t.Parallel()
					err := db.DeleteIngredient(context.Background(), v.id, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
	}
	for k, v := range tcs {
		t.Run(k, v)
	}
}

func Test_Lifecycler(t *testing.T) {
	tcs := map[string]func(t *testing.T){
		"SelectLifecycle": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, id types.UUID, cid types.CID) (types.Lifecycle, error)
				id     types.UUID
				result types.Lifecycle
				err    error
			}{
				"happy_path": {
					fn: db.SelectLifecycle,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					result, err := v.fn(context.Background(), v.id, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"InsertLifecycle": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, lc types.Lifecycle, cid types.CID) (types.Lifecycle, error)
				lc     types.Lifecycle
				result types.Lifecycle
				err    error
			}{
				"happy_path": {
					fn: db.InsertLifecycle,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					result, err := v.fn(context.Background(), v.lc, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"UpdateLifecycle": func(t *testing.T) {
			set := map[string]struct {
				fn  func(ctx context.Context, lc types.Lifecycle, cid types.CID) error
				lc  types.Lifecycle
				err error
			}{
				"happy_path": {
					fn: db.UpdateLifecycle,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), v.lc, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
		"DeleteLifecycle": func(t *testing.T) {
			set := map[string]struct {
				fn  func(ctx context.Context, id types.UUID, cid types.CID) error
				id  types.UUID
				err error
			}{
				"happy_path": {
					fn: db.DeleteLifecycle,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), v.id, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
	}
	for k, v := range tcs {
		t.Run(k, v)
	}
}

func Test_Stager(t *testing.T) {
	t.Parallel()
	tcs := map[string]func(t *testing.T){
		"SelectAllStages": func(t *testing.T) {
			set := map[string]struct {
				result []types.Stage
				err    error
			}{
				"happy_path": {
					result: []types.Stage{
						{UUID: "0", Name: "Gestation"},
						{UUID: "1", Name: "Colonization"},
						{UUID: "2", Name: "Majority"},
						{UUID: "3", Name: "Vacation"},
						{UUID: "4", Name: "System Test - Delete Me!!!"},
					},
				},
			}
			for k, v := range set {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					t.Parallel()
					result, err := db.SelectAllStages(context.Background(), types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"SelectStage": func(t *testing.T) {
			set := map[string]struct {
				id     types.UUID
				result types.Stage
				err    error
			}{
				"happy_path": {
					id:     "2",
					result: types.Stage{UUID: "2", Name: "Majority"},
				},
				"no_row_returned": {
					id:     "8",
					result: types.Stage{UUID: "0", Name: ""},
					err:    fmt.Errorf("sql: no rows in result set"),
				},
				"query_fails": {
					id:     "8888888888888888888888888888888888888888888888888888888888888888",
					result: types.Stage{UUID: "0", Name: ""},
					err:    fmt.Errorf("sql: no rows in result set"), // XXX
				},
			}
			for k, v := range set {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					t.Parallel()
					result, err := db.SelectStage(context.Background(), v.id, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"InsertStage": func(t *testing.T) {
			set := map[string]struct {
				s   types.Stage
				err error
			}{
				"happy_path": {
					s: types.Stage{Name: "bogus!"},
				},
				"no_rows_affected": {}, // ???
				"query_fails": {
					s:   types.Stage{Name: "01234567890123456789012345678901234567891"},
					err: fmt.Errorf(""),
				},
				"duplicate_name_violation": {
					s:   types.Stage{Name: "Colonization"},
					err: fmt.Errorf(""),
				},
			}
			for k, v := range set {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					t.Parallel()
					result, err := db.InsertStage(context.Background(), v.s, types.CID(k))
					require.Equal(t, v.err, err)
					require.NotEmpty(t, result.UUID)
				})
			}
		},
		"UpdateStage": func(t *testing.T) {
			set := map[string]struct {
				id  types.UUID
				s   types.Stage
				err error
			}{
				"happy_path": {
					id: "4",
					s:  types.Stage{Name: "Renamed"},
				},
				"no_rows_affected": {
					id:  "12",
					err: fmt.Errorf("stage was not updated: '12'"),
				},
				"query_fails": {
					id: "12121212121212121212121212121212121212121",
				},
				"duplicate_name_violation": {
					id: "0",
					s:  types.Stage{Name: "Vacation"},
				},
			}
			for k, v := range set {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					t.Parallel()
					err := db.UpdateStage(context.Background(), v.id, v.s, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
		"DeleteStage": func(t *testing.T) {
			set := map[string]struct {
				id  types.UUID
				err error
			}{
				"happy_path": {
					id: "4",
				},
				"no_rows_affected": {
					id:  "9",
					err: fmt.Errorf("stage could not be deleted: '9'"),
				},
				"query_fails": {
					id:  "00000000000000000000000000000000000000000000",
					err: fmt.Errorf(""),
				},
			}
			for k, v := range set {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					t.Parallel()
					err := db.DeleteStage(context.Background(), v.id, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
	}
	for k, v := range tcs {
		t.Run(k, v)
	}
}

func Test_StrainAttributer(t *testing.T) {
	tcs := map[string]func(t *testing.T){
		"KnownAttributeNames": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, cid types.CID) ([]string, error)
				result []string
				err    error
			}{
				"happy_path": {
					fn: db.KnownAttributeNames,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					result, err := v.fn(context.Background(), types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"GetAllAttributes": func(t *testing.T) {
			set := map[string]struct {
				fn  func(ctx context.Context, s *types.Strain, cid types.CID) error
				s   types.Strain
				err error
			}{
				"happy_path": {
					fn: db.GetAllAttributes,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), &v.s, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
		"AddAttribute": func(t *testing.T) {
			set := map[string]struct {
				fn   func(ctx context.Context, s *types.Strain, n, v string, cid types.CID) error
				s    types.Strain
				n, v string
				err  error
			}{
				"happy_path": {
					fn: db.AddAttribute,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), &v.s, v.n, v.v, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
		"ChangeAttribute": func(t *testing.T) {
			set := map[string]struct {
				fn   func(ctx context.Context, s *types.Strain, id types.UUID, n, v string, cid types.CID) error
				s    types.Strain
				id   types.UUID
				n, v string
				err  error
			}{
				"happy_path": {
					fn: db.ChangeAttribute,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), &v.s, v.id, v.n, v.v, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
		"RemoveAttribute": func(t *testing.T) {
			set := map[string]struct {
				fn  func(ctx context.Context, s *types.Strain, id types.UUID, cid types.CID) error
				s   types.Strain
				id  types.UUID
				err error
			}{
				"happy_path": {
					fn: db.RemoveAttribute,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), &v.s, v.id, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
	}
	for k, v := range tcs {
		t.Run(k, v)
	}
}

func Test_Strainer(t *testing.T) {
	tcs := map[string]func(t *testing.T){
		"SelectAllStrains": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, cid types.CID) ([]types.Strain, error)
				result []types.Strain
				err    error
			}{
				"happy_path": {
					fn: db.SelectAllStrains,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					result, err := v.fn(context.Background(), types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"SelectStrain": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, id types.UUID, cid types.CID) (types.Strain, error)
				id     types.UUID
				result types.Strain
				err    error
			}{
				"happy_path": {
					fn: db.SelectStrain,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					result, err := v.fn(context.Background(), v.id, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"InsertStrain": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, s types.Strain, cid types.CID) (types.Strain, error)
				s      types.Strain
				result types.Strain
				err    error
			}{
				"happy_path": {
					fn: db.InsertStrain,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					result, err := v.fn(context.Background(), v.s, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"UpdateStrain": func(t *testing.T) {
			set := map[string]struct {
				fn  func(ctx context.Context, id types.UUID, s types.Strain, cid types.CID) error
				id  types.UUID
				s   types.Strain
				err error
			}{
				"happy_path": {
					fn: db.UpdateStrain,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), v.id, v.s, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
		"DeleteStrain": func(t *testing.T) {
			set := map[string]struct {
				fn  func(ctx context.Context, id types.UUID, cid types.CID) error
				id  types.UUID
				err error
			}{
				"happy_path": {
					fn: db.DeleteStrain,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), v.id, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
	}
	for k, v := range tcs {
		t.Run(k, v)
	}
}

func Test_SubstrateIngredienter(t *testing.T) {
	tcs := map[string]func(t *testing.T){
		"GetAllIngredients": func(t *testing.T) {
			set := map[string]struct {
				s      types.Substrate
				result []types.Ingredient
				err    error
			}{
				"happy_path": {
					s: types.Substrate{UUID: "1"},
					result: []types.Ingredient{
						{UUID: "3", Name: "White Millet"},
						{UUID: "12", Name: "Red Millet"},
					},
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := db.GetAllIngredients(context.Background(), &v.s, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, v.s.Ingredients)
				})
			}
		},
		"AddIngredient": func(t *testing.T) {
			set := map[string]struct {
				s      types.Substrate
				i      types.Ingredient
				result []types.Ingredient
				err    error
			}{
				"happy_path": {
					s:      types.Substrate{UUID: "0"},
					i:      types.Ingredient{UUID: "9"},
					result: []types.Ingredient{{UUID: "9"}},
				},
				"duplicate_key_violation": {
					s:   types.Substrate{UUID: "0"},
					i:   types.Ingredient{UUID: "2"},
					err: fmt.Errorf("duplicate key violation"),
				},
				"no_rows_affected_ingredient": {
					s:   types.Substrate{UUID: "0"},
					i:   types.Ingredient{UUID: "-2"},
					err: fmt.Errorf("substrateingredient was not added"),
				},
				"no_rows_affected_substrate": {
					s:   types.Substrate{UUID: "-0"},
					i:   types.Ingredient{UUID: "2"},
					err: fmt.Errorf("substrateingredient was not added"),
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := db.AddIngredient(context.Background(), &v.s, v.i, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, v.s.Ingredients)
				})
			}
		},
		"ChangeIngredient": func(t *testing.T) {
			set := map[string]struct {
				s          types.Substrate
				oldI, newI types.Ingredient
				result     []types.Ingredient
				err        error
			}{
				"happy_path": {
					s: types.Substrate{UUID: "1", Ingredients: []types.Ingredient{
						{UUID: "3"},
						{UUID: "12"},
					}},
					oldI: types.Ingredient{UUID: "3"},
					newI: types.Ingredient{UUID: "4"},
					result: []types.Ingredient{
						{UUID: "4"},
						{UUID: "12"},
					},
				},
				"no_rows_affected": {
					s: types.Substrate{UUID: "1", Ingredients: []types.Ingredient{
						{UUID: "3"},
						{UUID: "12"},
					}},
					oldI: types.Ingredient{UUID: "4"},
					newI: types.Ingredient{UUID: "4"},
					result: []types.Ingredient{
						{UUID: "3"},
						{UUID: "12"},
					},
					err: fmt.Errorf("substrateingredient was not changed"),
				},
				"unique_key_violation": {
					s: types.Substrate{UUID: "1", Ingredients: []types.Ingredient{
						{UUID: "3"},
						{UUID: "12"},
					}},
					oldI: types.Ingredient{UUID: "3"},
					newI: types.Ingredient{UUID: "12"},
					result: []types.Ingredient{
						{UUID: "3"},
						{UUID: "12"},
					},
					err: fmt.Errorf("unique key violation"),
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := db.ChangeIngredient(context.Background(), &v.s, v.oldI, v.newI, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, v.s.Ingredients)
				})
			}
		},
		"RemoveIngredient": func(t *testing.T) {
			set := map[string]struct {
				s      types.Substrate
				i      types.Ingredient
				result []types.Ingredient
				err    error
			}{
				"happy_path": {
					s: types.Substrate{UUID: "0", Ingredients: []types.Ingredient{
						{UUID: "1"},
						{UUID: "2"},
						{UUID: "3"},
					}},
					i: types.Ingredient{UUID: "2"},
					result: []types.Ingredient{
						{UUID: "1"},
						{UUID: "3"},
					},
				},
				"no_rows_affected": {
					s:   types.Substrate{UUID: "0"},
					i:   types.Ingredient{UUID: "12"},
					err: fmt.Errorf("substrateingredient was not removed"),
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := db.RemoveIngredient(context.Background(), &v.s, v.i, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, v.s.Ingredients)
				})
			}
		},
	}
	for k, v := range tcs {
		t.Run(k, v)
	}
}

func Test_Substrater(t *testing.T) {
	tcs := map[string]func(t *testing.T){
		"SelectAllSubstrates": func(t *testing.T) {
			set := map[string]struct {
				result []types.Substrate
				err    error
			}{
				"happy_path": {
					result: []types.Substrate{
						{UUID: "0", Name: "Rye", Type: "Grain", Vendor: types.Vendor{}},
						{UUID: "1", Name: "Millet", Type: "Grain", Vendor: types.Vendor{}},
					},
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					result, err := db.SelectAllSubstrates(context.Background(), types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"SelectSubstrate": func(t *testing.T) {
			set := map[string]struct {
				id     types.UUID
				result types.Substrate
				err    error
			}{
				"happy_path": {
					id:     "0",
					result: types.Substrate{UUID: "0", Name: "Rye", Type: "Grain", Vendor: types.Vendor{}},
				},
				"no_rows_returned": {
					id:  "5",
					err: noRows,
				},
				"query_fails": {
					id:  "01234567890123456789012345678901234567891",
					err: noRows,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					result, err := db.SelectSubstrate(context.Background(), v.id, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"InsertSubstrate": func(t *testing.T) {
			set := map[string]struct {
				s      types.Substrate
				result types.Substrate
				err    error
			}{
				"happy_path": {
					s:      types.Substrate{Name: "Honey Solution", Type: "Bulk", Vendor: types.Vendor{UUID: "0"}},
					result: types.Substrate{},
				},
				"unique_key_violation": {
					s:      types.Substrate{Name: "Rye", Type: "Bulk", Vendor: types.Vendor{UUID: "0"}},
					result: types.Substrate{},
					err:    fmt.Errorf("duplicate key violation"),
				},
				"check_constraint_violation": {
					s:      types.Substrate{Name: "Maltodexterin", Type: "Stardust", Vendor: types.Vendor{UUID: "0"}},
					result: types.Substrate{},
					err:    fmt.Errorf("check constraint"),
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					result, err := db.InsertSubstrate(context.Background(), v.s, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"UpdateSubstrate": func(t *testing.T) {
			set := map[string]struct {
				id  types.UUID
				s   types.Substrate
				err error
			}{
				"happy_path": {
					s: types.Substrate{},
				},
				"unique_key_violation": {
					id:  "0",
					s:   types.Substrate{Name: "Millet", Type: "Bulk", Vendor: types.Vendor{UUID: "0"}},
					err: fmt.Errorf("duplicate key violation"),
				},
				"check_constraint_violation": {
					id:  "0",
					s:   types.Substrate{Name: "Maltodexterin", Type: "Stardust", Vendor: types.Vendor{UUID: "0"}},
					err: fmt.Errorf("check constraint"),
				},
				"query_fails": {
					id: "12121212121212121212121212121212121212121",
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := db.UpdateSubstrate(context.Background(), v.id, v.s, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
		"DeleteSubstrate": func(t *testing.T) {
			set := map[string]struct {
				id  types.UUID
				err error
			}{
				"happy_path": {
					id: "-1",
				},
				"no_rows_affected": {
					id:  "12",
					err: fmt.Errorf("vendor table was not deleted '12'"),
				},
				"query_fails": {
					id:  "01234567890123456789012345678901234567891",
					err: fmt.Errorf("some error"),
				},
				"referential_violation": {
					id:  "0",
					err: fmt.Errorf("referential constraint"),
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := db.DeleteSubstrate(context.Background(), v.id, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
	}
	for k, v := range tcs {
		t.Run(k, v)
	}
}

func Test_Vendorer(t *testing.T) {
	t.Parallel()
	tcs := map[string]func(t *testing.T){
		"SelectAllVendors": func(t *testing.T) {
			set := map[string]struct {
				result []types.Vendor
				err    error
			}{
				"happy_path": {
					result: []types.Vendor{{UUID: "0", Name: "127.0.0.1"}},
				},
			}
			for k, v := range set {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					t.Parallel()
					result, err := db.SelectAllVendors(context.Background(), types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"SelectVendor": func(t *testing.T) {
			set := map[string]struct {
				id     types.UUID
				result types.Vendor
				err    error
			}{
				"happy_path": {
					id:     "0",
					result: types.Vendor{UUID: "0", Name: "127.0.0.1"},
				},
				"no_row_returned": {
					id:     "8",
					result: types.Vendor{UUID: "8", Name: ""},
					err:    fmt.Errorf("sql: no rows in result set"),
				},
				"query_fails": {
					id:     "8888888888888888888888888888888888888888888888888888888888888888",
					result: types.Vendor{UUID: "0", Name: ""},
					err:    noRows, // XXX
				},
			}
			for k, v := range set {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					t.Parallel()
					result, err := db.SelectVendor(context.Background(), v.id, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"InsertVendor": func(t *testing.T) {
			set := map[string]struct {
				v      types.Vendor
				result types.Vendor
				err    error
			}{
				"happy_path": {
					v: types.Vendor{},
				},
				"duplicate_name_violation": {
					v: types.Vendor{Name: "127.0.0.1"},
				},
			}
			for k, v := range set {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					t.Parallel()
					result, err := db.InsertVendor(context.Background(), v.v, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"UpdateVendor": func(t *testing.T) {
			set := map[string]struct {
				id  types.UUID
				v   types.Vendor
				err error
			}{
				"happy_path": {
					id: "0",
					v:  types.Vendor{Name: "localhost"},
				},
				"duplicate_name_violation": {
					id: "1",
					v:  types.Vendor{Name: "127.0.0.1"},
				},
				"no_rows_affected": {
					id:  "12",
					err: fmt.Errorf("vendor table was not updated '12'"),
				},
				"query_fails": {
					id: "12121212121212121212121212121212121212121",
				},
			}
			for k, v := range set {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					t.Parallel()
					err := db.UpdateVendor(context.Background(), v.id, v.v, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
		"DeleteVendor": func(t *testing.T) {
			set := map[string]struct {
				id  types.UUID
				err error
			}{
				"happy_path": {
					id: "-1",
				},
				"no_rows_affected": {
					id:  "12",
					err: fmt.Errorf("vendor table was not deleted '12'"),
				},
				"query_fails": {
					id:  "01234567890123456789012345678901234567891",
					err: fmt.Errorf("some error"),
				},
				"referential_violation": {
					id:  "0",
					err: fmt.Errorf("referential constraint"),
				},
			}
			for k, v := range set {
				k, v := k, v
				t.Run(k, func(t *testing.T) {
					t.Parallel()
					err := db.DeleteVendor(context.Background(), v.id, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
	}
	for k, v := range tcs {
		t.Run(k, v)
	}
}
