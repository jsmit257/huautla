package test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

var generations []types.Generation

func init() {
	for _, id := range []types.UUID{"0", "1", "2", "3", "4"} {
		if g, err := db.SelectGeneration(context.Background(), id, "generation_init"); err != nil {
			panic(fmt.Errorf("select error: %v, %w", id, err))
		} else {
			generations = append(generations, g)
		}
	}
}

func Test_SelectGenerationIndex(t *testing.T) {
	t.Parallel()

	result, err := db.SelectGenerationIndex(context.Background(), types.CID("Test_SelectGenerationIndex"))
	require.Nil(t, err)
	require.LessOrEqual(t, 5, len(result))
}

func Test_SelectGeneration(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		err error
	}{
		// "happy_path": { // XXX: kinda redundant
		// 	id:     generation[0].UUID,
		// 	result: generation[0],
		// },
		"no_rows_returned": {
			err: sql.ErrNoRows,
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.SelectGeneration(context.Background(), "abc", types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, types.UUID(""), result.UUID)
			// require.Equal(t, v.result, result)
		})
	}
}

func Test_InsertGeneration(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		g      types.Generation
		result types.Generation
		err    error
	}{
		"happy_path": {
			g: types.Generation{
				PlatingSubstrate: substrates[types.PlatingType][0],
				LiquidSubstrate:  substrates[types.LiquidType][0],
			},
			result: types.Generation{
				PlatingSubstrate: substrates[types.PlatingType][0],
				LiquidSubstrate:  substrates[types.LiquidType][0],
			},
		},
		"no_rows_affected_plating": {
			g: types.Generation{
				// PlatingSubstrate: substrates[types.PlatingType][0],
				LiquidSubstrate: substrates[types.LiquidType][0],
			},
			err: fmt.Errorf("generation was not added: 0"),
		},
		"no_rows_affected_liquid": {
			g: types.Generation{
				PlatingSubstrate: substrates[types.PlatingType][0],
				// LiquidSubstrate:  substrates[types.LiquidType][0],
			},
			err: fmt.Errorf("generation was not added: 0"),
		},
		"no_rows_affected_check_plating_type": {
			g: types.Generation{
				PlatingSubstrate: substrates[types.GrainType][0],
				LiquidSubstrate:  substrates[types.LiquidType][0],
			},
			err: fmt.Errorf("generation was not added: 0"),
		},
		"no_rows_affected_check_liquid_type": {
			g: types.Generation{
				PlatingSubstrate: substrates[types.PlatingType][0],
				LiquidSubstrate:  substrates[types.BulkType][0],
			},
			err: fmt.Errorf("generation was not added: 0"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.InsertGeneration(context.Background(), v.g, types.CID(k))
			equalErrorMessages(t, v.err, err)
			if err != nil {
				require.NotEmpty(t, result.UUID)
			}
		})
	}
}

func Test_UpdateGeneration(t *testing.T) {
	t.Parallel()

	updated, err := db.SelectGeneration(context.Background(), "update me!", "Test_UpdateGeneration")
	require.Nil(t, err)

	set := map[string]struct {
		xform func(types.Generation) types.Generation
		err   error
	}{
		"happy_path": {
			xform: func(g types.Generation) types.Generation {
				g.PlatingSubstrate = substrates[types.PlatingType][0]
				g.LiquidSubstrate = substrates[types.LiquidType][1]
				return g
			},
		},
		"no_rows_affected_plating": {
			xform: func(g types.Generation) types.Generation {
				g.PlatingSubstrate = types.Substrate{UUID: "missing"}
				return g
			},
			err: fmt.Errorf("generation was not updated"),
		},
		"no_rows_affected_liquid": {
			xform: func(g types.Generation) types.Generation {
				g.LiquidSubstrate = types.Substrate{UUID: "missing"}
				return g
			},
			err: fmt.Errorf("generation was not updated"),
		},
		"check_plating_type": {
			xform: func(g types.Generation) types.Generation {
				g.PlatingSubstrate = substrates[types.GrainType][0]
				return g
			},
			err: fmt.Errorf("generation was not updated"),
		},
		"check_liquid_type": {
			xform: func(g types.Generation) types.Generation {
				g.LiquidSubstrate = substrates[types.BulkType][0]
				return g
			},
			err: fmt.Errorf("generation was not updated"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			g := v.xform(updated)
			_, err := db.UpdateGeneration(context.Background(), g, types.CID(k))
			equalErrorMessages(t, v.err, err)
		})
	}
}

func Test_DeleteGeneration(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id  types.UUID
		err error
	}{
		"happy_path": {
			id: "delete me!",
		},
		"no_rows_affected": {
			id:  "missing",
			err: fmt.Errorf("generation could not be deleted: 'missing'"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.DeleteGeneration(context.Background(), v.id, types.CID(k))
			equalErrorMessages(t, v.err, err)
		})
	}
}

func Test_GeneratedStrain(t *testing.T) {
	t.Skip()
	t.Parallel()

	err := db.UpdateGeneratedStrain(context.Background(), &generations[0].UUID, strains[0].UUID, types.CID("LinkingStrains"))
	require.Nil(t, err, "linking strain to generation")

	strains[0].Generation = &types.Generation{UUID: generations[0].UUID}

	set := map[string]struct {
		id     types.UUID
		result types.Strain
		err    error
	}{
		"happy_path": {
			id: generations[0].UUID,
			result: func(s types.Strain) types.Strain {
				s.Attributes = nil
				return s
			}(strains[0]),
		},
		"no_strain_found": {
			id:  generations[1].UUID,
			err: sql.ErrNoRows,
		},
		"no_generation_found": {
			id:  "missing",
			err: sql.ErrNoRows,
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.GeneratedStrain(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err, "result: %#v", result)
			require.Equal(t, v.result, result)
		})
	}
}
func Test_UpdateGeneratedStrain(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		gid,
		sid types.UUID
		err error
	}{
		"happy_update_new": {
			gid: generations[2].UUID,
			sid: strains[1].UUID,
		},
		"happy_update_existing": {
			gid: generations[1].UUID,
			sid: strains[2].UUID,
		},
		"happy_delete": {
			sid: strains[2].UUID,
		},
		// FIXME: decide if these need to work
		"no_rows_affected_strain": {
			gid: generations[0].UUID,
			err: sql.ErrNoRows,
		},
		"no_rows_affected_generation": { // vendors aren't part of the update (yet)
			gid: generations[0].UUID,
			err: sql.ErrNoRows,
		},
		// "unique_key_violation": {
		// 	gid: generations[2].UUID,
		// 	sid: strains[4].UUID,
		// 	err: fmt.Errorf(uniqueKeyViolation, "strains_name_vendor_uuid_ctime_key"),
		// },
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			// t.Parallel()

			tmp := (*types.UUID)(nil)
			if v.gid != "" {
				tmp = &v.gid
			}

			err := db.UpdateGeneratedStrain(context.Background(), tmp, v.sid, types.CID(k))
			equalErrorMessages(t, v.err, err)
		})
	}
}
