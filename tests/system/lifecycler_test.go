package test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

var (
	lifecycles []types.Lifecycle
	// imp        = "!impossible!"
)

func init() {
	for _, id := range []string{"0", "1"} {
		if l, err := db.SelectLifecycle(context.Background(), types.UUID(id), "lifecycle_init"); err != nil {
			panic(fmt.Errorf("failed with err: %v", err))
		} else {
			lifecycles = append(lifecycles, l)
		}
	}
}

func Test_SelectLifecycleIndex(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		result []types.Lifecycle
		err    error
	}{
		"happy_path": {
			result: []types.Lifecycle{
				{UUID: lifecycles[0].UUID, Location: lifecycles[0].Location, CTime: lifecycles[0].CTime},
				{UUID: lifecycles[1].UUID, Location: lifecycles[1].Location, CTime: lifecycles[1].CTime},
			},
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.SelectLifecycleIndex(context.Background(), types.CID(k))
			require.Equal(t, v.err, err)
			require.LessOrEqual(t, 2, len(result))
		})
	}
}

func Test_SelectLifecycle(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     types.UUID
		result types.Lifecycle
		err    error
	}{
		// "happy_path": { // XXX: kinda redundant
		// 	id:     lifecycles[0].UUID,
		// 	result: lifecycles[0],
		// },
		"no_rows_returned": {
			id:     "missing",
			result: types.Lifecycle{UUID: "missing"},
			err:    sql.ErrNoRows,
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.SelectLifecycle(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result.Name, result.Name)
			// require.Equal(t, v.result, result)
		})
	}
}

func Test_InsertLifecycle(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		lc     types.Lifecycle
		result types.Lifecycle
		err    error
	}{
		"happy_path": {
			lc: types.Lifecycle{
				Location:       "inserted record",
				GrainCost:      1,
				BulkCost:       2,
				Yield:          3,
				Count:          4,
				Gross:          5,
				Strain:         strains[1],
				GrainSubstrate: substrates[types.GrainType][0],
				BulkSubstrate:  substrates[types.BulkType][0],
			},
			result: types.Lifecycle{
				Location:       "inserted record",
				GrainCost:      1,
				BulkCost:       2,
				Yield:          3,
				Count:          4,
				Gross:          5,
				Strain:         strains[1],
				GrainSubstrate: substrates[types.GrainType][0],
				BulkSubstrate:  substrates[types.BulkType][0],
			},
		},
		"no_rows_affected_grain": {
			lc: types.Lifecycle{
				Location:       "failed insert",
				Strain:         strains[1],
				GrainSubstrate: types.Substrate{UUID: "foobar"},
				BulkSubstrate:  substrates[types.BulkType][0],
			},
			// result: types.Lifecycle{Name: "failed insert"},
			err: fmt.Errorf("lifecycle was not added: 0"),
		},
		"no_rows_affected_bulk": {
			lc: types.Lifecycle{
				Location:       "failed insert",
				Strain:         strains[1],
				GrainSubstrate: substrates[types.GrainType][0],
				BulkSubstrate:  types.Substrate{UUID: "foobar"},
			},
			// result: types.Lifecycle{Name: "failed insert"},
			err: fmt.Errorf("lifecycle was not added: 0"),
		},
		"no_rows_affected_strain": {
			lc: types.Lifecycle{
				Location:       "failed insert",
				Strain:         types.Strain{UUID: "foobar"},
				GrainSubstrate: substrates[types.GrainType][0],
				BulkSubstrate:  substrates[types.BulkType][0],
			},
			// result: types.Lifecycle{Name: "failed insert"},
			err: fmt.Errorf("lifecycle was not added: 0"),
		},
		"no_rows_affected_check_grain_type": {
			lc: types.Lifecycle{
				Location:       "failed insert",
				Strain:         strains[0],
				GrainSubstrate: substrates[types.BulkType][0],
				BulkSubstrate:  substrates[types.BulkType][0],
			},
			// result: types.Lifecycle{Name: "failed insert"},
			err: fmt.Errorf("lifecycle was not added: 0"),
		},
		"no_rows_affected_check_bulk_type": {
			lc: types.Lifecycle{
				Location:       "failed insert",
				Strain:         strains[0],
				GrainSubstrate: substrates[types.GrainType][0],
				BulkSubstrate:  substrates[types.GrainType][0],
			},
			// result: types.Lifecycle{Name: "failed insert"},
			err: fmt.Errorf("lifecycle was not added: 0"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.InsertLifecycle(context.Background(), v.lc, types.CID(k))
			equalErrorMessages(t, v.err, err)
			if err != nil {
				require.NotEmpty(t, result.UUID)
			}
		})
	}
}

func Test_UpdateLifecycle(t *testing.T) {
	t.Parallel()

	lc, err := db.SelectLifecycle(context.Background(), "update me!", "Test_UpdateLifecycle")
	require.Nil(t, err)

	set := map[string]struct {
		xform func(types.Lifecycle) types.Lifecycle
		err   error
	}{
		"happy_path": {
			xform: func(lc types.Lifecycle) types.Lifecycle {
				lc.Location = "updated"
				return lc
			},
		},
		"no_rows_affected_strain": {
			xform: func(lc types.Lifecycle) types.Lifecycle {
				lc.Strain = types.Strain{UUID: "missing"}
				return lc
			},
			err: fmt.Errorf("one of strain, grain or bulk is not the right type"),
		},
		"no_rows_affected_grain": {
			xform: func(lc types.Lifecycle) types.Lifecycle {
				lc.GrainSubstrate = types.Substrate{UUID: "missing"}
				return lc
			},
			err: fmt.Errorf("one of strain, grain or bulk is not the right type"),
		},
		"no_rows_affected_bulk": {
			xform: func(lc types.Lifecycle) types.Lifecycle {
				lc.BulkSubstrate = types.Substrate{UUID: "missing"}
				return lc
			},
			err: fmt.Errorf("one of strain, grain or bulk is not the right type"),
		},
		"check_grain_type": {
			xform: func(lc types.Lifecycle) types.Lifecycle {
				lc.GrainSubstrate = substrates[types.BulkType][0]
				return lc
			},
			err: fmt.Errorf("one of strain, grain or bulk is not the right type"),
		},
		"check_bulk_type": {
			xform: func(lc types.Lifecycle) types.Lifecycle {
				lc.BulkSubstrate = substrates[types.GrainType][0]
				return lc
			},
			err: fmt.Errorf("one of strain, grain or bulk is not the right type"),
		},
		"unique_key_violation": {
			xform: func(lc types.Lifecycle) types.Lifecycle {
				lc.Location = "reference implementation 2"
				return lc
			},
			err: fmt.Errorf(uniqueKeyViolation, "lifecycles_location_ctime_key"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			lc := v.xform(lc)
			_, err := db.UpdateLifecycle(context.Background(), lc, types.CID(k))
			equalErrorMessages(t, v.err, err)
		})
	}
}

func Test_DeleteLifecycle(t *testing.T) {
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
			err: fmt.Errorf("lifecycle could not be deleted: 'missing'"),
		},
		"referential_violation": {
			id:  "0",
			err: fmt.Errorf("pq: foreign key violation"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.DeleteLifecycle(context.Background(), v.id, types.CID(k))
			equalErrorMessages(t, v.err, err)
		})
	}
}
