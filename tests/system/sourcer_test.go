package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

func Test_AddStrainSource(t *testing.T) {
	t.Parallel()

	generation, err := db.SelectGeneration(context.Background(), "add source", "Test_AddStrainSource")
	require.Nil(t, err)
	failsourcecheck, err := db.SelectGeneration(context.Background(), "fail source check", "Test_AddStrainSource")
	require.Nil(t, err)

	type v struct {
		s      types.Source
		g      *types.Generation
		result []types.Source
		err    error
	}
	set := []struct {
		k string
		v v
	}{
		{
			k: "duplicate_key_violation",
			v: v{
				s: types.Source{
					Type:   "Spore",
					Strain: strains[2],
				},
				err: fmt.Errorf(uniqueKeyViolation, "sources_progenitor_uuid_generation_uuid_key"),
			},
		},
		{
			k: "fail_mixed_sources",
			v: v{
				s: types.Source{
					Type:   "Clone",
					Strain: strains[2],
				},
				err: fmt.Errorf("pq: source types can't be mixed"),
			},
		},
		{
			k: "fail_type",
			v: v{
				s: types.Source{
					Type:   "Fail",
					Strain: strains[2],
				},
				g:   &failsourcecheck,
				err: fmt.Errorf(checkConstraintViolation, "sources", "sources_type_check"),
			},
		},
		{
			k: "no_rows_affected_strain",
			v: v{
				s: types.Source{
					Type:   "Clone",
					Strain: types.Strain{UUID: "missing"},
				},
				err: fmt.Errorf("pq: no existing progenitor"),
			},
		},
		{
			k: "happy_path",
			v: v{
				s: types.Source{
					Type:   "Spore",
					Strain: strains[3],
				},
				result: append(generation.Sources, types.Source{Type: "Spore", Strain: strains[1]}),
			},
		},
		{
			k: "too_many_spore_sources",
			v: v{
				s: types.Source{
					Type:   "Spore",
					Strain: strains[4],
				},
				err: fmt.Errorf("pq: too many sources for this generation"),
			},
		},
	}

	for _, tc := range set {
		k, v, generation := tc.k, tc.v, generation
		t.Run(k, func(t *testing.T) {
			if v.g == nil {
				v.g = &generation
			}
			_, err := db.InsertSource(context.Background(), v.g.UUID, "strain", v.s, types.CID(k))
			equalErrorMessages(t, v.err, err)
		})
	}
}

func Test_AddEventSource(t *testing.T) {
	t.Parallel()
	t.Skip()

	generation, err := db.SelectGeneration(context.Background(), "add event source", "Test_AddEventSource")
	require.Nil(t, err)

	type v struct {
		e      types.Event
		g      *types.Generation
		result []types.Source
		err    error
	}

	set := []struct {
		k string
		v v
	}{
		{
			k: "fail_mixed_sources",
			v: v{
				e:   events[6],
				err: fmt.Errorf("pq: source types can't be mixed"),
			},
		},
		{
			k: "fail_non_generation_event",
			v: v{
				e:   events[0],
				err: fmt.Errorf("pq: event is not a generation type"),
			},
		},
		{
			k: "duplicate_key_violation",
			v: v{
				e:   events[3],
				err: fmt.Errorf(uniqueKeyViolation, "sources_progenitor_uuid_generation_uuid_key"),
			},
		},
		{
			k: "happy_path",
			v: v{
				e:      events[4],
				result: append(generation.Sources, types.Source{Type: "Spore", Strain: strains[1]}),
			},
		},
		{
			k: "too_many_spore_sources",
			v: v{
				e:   events[5],
				err: fmt.Errorf("pq: too many sources for this generation"),
			},
		},
	}

	for _, tc := range set {
		k, v, generation := tc.k, tc.v, generation
		t.Run(k, func(t *testing.T) {
			if v.g == nil {
				v.g = &generation
			}
			_, err := db.InsertSource(context.Background(), v.g.UUID, "event", types.Source{
				Lifecycle: &types.Lifecycle{
					Events: []types.Event{v.e},
				},
			}, types.CID(k))
			equalErrorMessages(t, v.err, err)
		})
	}
}

func Test_ChangeSource(t *testing.T) {
	t.Parallel()

	g, err := db.SelectGeneration(context.Background(), "change source", "Test_ChangeSource")
	require.Nil(t, err)
	g2, err := db.SelectGeneration(context.Background(), "change_source_fail_type", "Test_ChangeSource")
	require.Nil(t, err)

	set := map[string]struct {
		s      types.Source
		origin string
		result []types.Source
		err    error
	}{
		"happy_path": {
			origin: "strain",
			s: func(s types.Source) types.Source {
				// s.Type = "Clone"
				s.Strain = strains[3]
				return s
			}(g.Sources[0]),
			result: append(
				[]types.Source{
					func(s types.Source) types.Source {
						// s.Type = "Clone"
						s.Strain = strains[3]
						return s
					}(g.Sources[0]),
				},
				g.Sources[1],
			),
		},
		"cant_mix_types": {
			origin: "strain",
			s: func(s types.Source) types.Source {
				s.Type = "Clone"
				return s
			}(g.Sources[0]),
			result: g.Sources,
			err:    fmt.Errorf("pq: source types can't be mixed"),
		},
		"fail_type": {
			origin: "strain",
			s: func(s types.Source) types.Source {
				s.Type = "Fail"
				return s
			}(g2.Sources[0]),
			result: g2.Sources,
			err:    fmt.Errorf(checkConstraintViolation, "sources", "sources_type_check"),
		},
	}
	for name, tc := range set {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := db.UpdateSource(context.Background(), tc.origin, tc.s, types.CID(name))
			equalErrorMessages(t, tc.err, err)
		})
	}
}

func Test_RemoveSource(t *testing.T) {
	t.Parallel()

	g, err := db.SelectGeneration(context.Background(), "remove source", "Test_RemoveIngredient")
	require.Nil(t, err)

	set := map[string]struct {
		id     types.UUID
		result []types.Source
		err    error
	}{
		"happy_path": { // happy path has to run first
			id:     "delete me!",
			result: []types.Source{},
		},
		"no_rows_affected_ingredient": {
			id:     "missing",
			result: g.Sources,
			err:    fmt.Errorf("source could not be deleted: 'missing'"),
		},
	}
	for k, v := range set {
		k, v, g := k, v, g
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.RemoveSource(context.Background(), &g, v.id, types.CID(k))
			equalErrorMessages(t, v.err, err)
			require.ElementsMatch(t, v.result, g.Sources)
		})
	}
}
