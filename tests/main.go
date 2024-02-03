package main

//package tests ??

import (
	"context"
	"testing"

	"github.com/jsmit257/huautla/internal/data"
	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

func main() {}

var db types.DB

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
	tcs := map[string]func(t *testing.T){
		"SelectAllIngredients": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, cid types.CID) ([]types.Ingredient, error)
				result []types.Ingredient
				err    error
			}{
				"happy_path": {
					fn: db.SelectAllIngredients,
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
		"SelectIngredient": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, id types.UUID, cid types.CID) (types.Ingredient, error)
				id     types.UUID
				result types.Ingredient
				err    error
			}{
				"happy_path": {
					fn: db.SelectIngredient,
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
		"InsertIngredient": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, i types.Ingredient, cid types.CID) (types.Ingredient, error)
				i      types.Ingredient
				result []types.Ingredient
				err    error
			}{
				"happy_path": {
					fn: db.InsertIngredient,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					result, err := v.fn(context.Background(), v.i, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"UpdateIngredient": func(t *testing.T) {
			set := map[string]struct {
				fn  func(ctx context.Context, id types.UUID, i types.Ingredient, cid types.CID) error
				id  types.UUID
				i   types.Ingredient
				err error
			}{
				"happy_path": {
					fn: db.UpdateIngredient,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), v.id, v.i, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
		"DeleteIngredient": func(t *testing.T) {
			set := map[string]struct {
				fn  func(ctx context.Context, id types.UUID, cid types.CID) error
				id  types.UUID
				err error
			}{
				"happy_path": {
					fn: db.DeleteIngredient,
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
	tcs := map[string]func(t *testing.T){
		"SelectAllStages": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, cid types.CID) ([]types.Stage, error)
				result []types.Stage
				err    error
			}{
				"happy_path": {
					fn: db.SelectAllStages,
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
		"SelectStage": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, id types.UUID, cid types.CID) (types.Stage, error)
				id     types.UUID
				result types.Stage
				err    error
			}{
				"happy_path": {
					fn: db.SelectStage,
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
		"InsertStage": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, s types.Stage, cid types.CID) (types.Stage, error)
				s      types.Stage
				result types.Stage
				err    error
			}{
				"happy_path": {
					fn: db.InsertStage,
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
		"UpdateStage": func(t *testing.T) {
			set := map[string]struct {
				fn  func(ctx context.Context, id types.UUID, s types.Stage, cid types.CID) error
				id  types.UUID
				s   types.Stage
				err error
			}{
				"happy_path": {
					fn: db.UpdateStage,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), v.id, v.s, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
		"DeleteStage": func(t *testing.T) {
			set := map[string]struct {
				fn  func(ctx context.Context, id types.UUID, cid types.CID) error
				id  types.UUID
				err error
			}{
				"happy_path": {
					fn: db.DeleteStage,
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
				fn     func(ctx context.Context, s *types.Substrate, cid types.CID) error
				s      types.Substrate
				result []types.Ingredient
				err    error
			}{
				"happy_path": {
					fn: db.GetAllIngredients,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), &v.s, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, v.s.Ingredients)
				})
			}
		},
		"AddIngredient": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, s *types.Substrate, i types.Ingredient, cid types.CID) error
				s      types.Substrate
				i      types.Ingredient
				result []types.Ingredient
				err    error
			}{
				"happy_path": {
					fn: db.AddIngredient,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), &v.s, v.i, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, v.s.Ingredients)
				})
			}
		},
		"ChangeIngredient": func(t *testing.T) {
			set := map[string]struct {
				fn         func(ctx context.Context, s *types.Substrate, oldI, newI types.Ingredient, cid types.CID) error
				s          types.Substrate
				oldI, newI types.Ingredient
				result     []types.Ingredient
				err        error
			}{
				"happy_path": {
					fn: db.ChangeIngredient,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), &v.s, v.oldI, v.newI, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, v.s.Ingredients)
				})
			}
		},
		"RemoveIngredient": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, s *types.Substrate, i types.Ingredient, cid types.CID) error
				s      types.Substrate
				i      types.Ingredient
				result []types.Ingredient
				err    error
			}{
				"happy_path": {
					fn: db.RemoveIngredient,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), &v.s, v.i, types.CID(k))
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
				fn     func(ctx context.Context, cid types.CID) ([]types.Substrate, error)
				result []types.Substrate
				err    error
			}{
				"happy_path": {
					fn: db.SelectAllSubstrates,
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
		"SelectSubstrate": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, id types.UUID, cid types.CID) (types.Substrate, error)
				id     types.UUID
				result types.Substrate
				err    error
			}{
				"happy_path": {
					fn: db.SelectSubstrate,
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
		"InsertSubstrate": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, s types.Substrate, cid types.CID) (types.Substrate, error)
				s      types.Substrate
				result types.Substrate
				err    error
			}{
				"happy_path": {
					fn: db.InsertSubstrate,
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
		"UpdateSubstrate": func(t *testing.T) {
			set := map[string]struct {
				fn  func(ctx context.Context, id types.UUID, s types.Substrate, cid types.CID) error
				id  types.UUID
				s   types.Substrate
				err error
			}{
				"happy_path": {
					fn: db.UpdateSubstrate,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), v.id, v.s, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
		"DeleteSubstrate": func(t *testing.T) {
			set := map[string]struct {
				fn  func(ctx context.Context, id types.UUID, cid types.CID) error
				id  types.UUID
				err error
			}{
				"happy_path": {
					fn: db.DeleteSubstrate,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), v.id, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
		"GetAllIngredients": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, s *types.Substrate, cid types.CID) error
				s      types.Substrate
				result []types.Ingredient
				err    error
			}{
				"happy_path": {
					fn: db.GetAllIngredients,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), &v.s, types.CID(k))
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

func Test_Vendorer(t *testing.T) {
	tcs := map[string]func(t *testing.T){
		"SelectAllVendors": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, cid types.CID) ([]types.Vendor, error)
				result []types.Vendor
				err    error
			}{
				"happy_path": {
					fn: db.SelectAllVendors,
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
		"SelectVendor": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, id types.UUID, cid types.CID) (types.Vendor, error)
				id     types.UUID
				result types.Vendor
				err    error
			}{
				"happy_path": {
					fn: db.SelectVendor,
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
		"InsertVendor": func(t *testing.T) {
			set := map[string]struct {
				fn     func(ctx context.Context, v types.Vendor, cid types.CID) (types.Vendor, error)
				v      types.Vendor
				result types.Vendor
				err    error
			}{
				"happy_path": {
					fn: db.InsertVendor,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					result, err := v.fn(context.Background(), v.v, types.CID(k))
					require.Equal(t, v.err, err)
					require.Equal(t, v.result, result)
				})
			}
		},
		"UpdateVendor": func(t *testing.T) {
			set := map[string]struct {
				fn  func(ctx context.Context, id types.UUID, v types.Vendor, cid types.CID) error
				id  types.UUID
				v   types.Vendor
				err error
			}{
				"happy_path": {
					fn: db.UpdateVendor,
				},
			}
			for k, v := range set {
				t.Run(k, func(t *testing.T) {
					err := v.fn(context.Background(), v.id, v.v, types.CID(k))
					require.Equal(t, v.err, err)
				})
			}
		},
		"DeleteVendor": func(t *testing.T) {
			set := map[string]struct {
				fn  func(ctx context.Context, id types.UUID, cid types.CID) error
				id  types.UUID
				err error
			}{
				"happy_path": {
					fn: db.DeleteVendor,
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
