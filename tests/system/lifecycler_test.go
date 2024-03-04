package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

var lifecycles []types.Lifecycle

func init() {
	for _, id := range []types.UUID{"0", "1"} {
		if l, err := db.SelectLifecycle(context.Background(), id, "lifecycle_init"); err != nil {
			panic(err)
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
			err:    noRows,
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
				GrainSubstrate: substrates[1],
				BulkSubstrate:  substrates[2],
			},
			result: types.Lifecycle{
				Location:       "inserted record",
				GrainCost:      1,
				BulkCost:       2,
				Yield:          3,
				Count:          4,
				Gross:          5,
				Strain:         strains[1],
				GrainSubstrate: substrates[1],
				BulkSubstrate:  substrates[2],
			},
		},
		"no_rows_affected_grain": {
			lc: types.Lifecycle{
				Location:       "failed insert",
				Strain:         strains[1],
				GrainSubstrate: types.Substrate{UUID: "foobar"},
				BulkSubstrate:  substrates[2],
			},
			// result: types.Lifecycle{Name: "failed insert"},
			err: fmt.Errorf("lifecycle was not added: 0"),
		},
		"no_rows_affected_bulk": {
			lc: types.Lifecycle{
				Location:       "failed insert",
				Strain:         strains[1],
				GrainSubstrate: substrates[1],
				BulkSubstrate:  types.Substrate{UUID: "foobar"},
			},
			// result: types.Lifecycle{Name: "failed insert"},
			err: fmt.Errorf("lifecycle was not added: 0"),
		},
		"no_rows_affected_strain": {
			lc: types.Lifecycle{
				Location:       "failed insert",
				Strain:         types.Strain{UUID: "foobar"},
				GrainSubstrate: substrates[1],
				BulkSubstrate:  substrates[2],
			},
			// result: types.Lifecycle{Name: "failed insert"},
			err: fmt.Errorf("lifecycle was not added: 0"),
		},
		"no_rows_affected_check_grain_type": {
			lc: types.Lifecycle{
				Location:       "failed insert",
				Strain:         strains[0],
				GrainSubstrate: substrates[2],
				BulkSubstrate:  substrates[2],
			},
			// result: types.Lifecycle{Name: "failed insert"},
			err: fmt.Errorf("lifecycle was not added: 0"),
		},
		"no_rows_affected_check_bulk_type": {
			lc: types.Lifecycle{
				Location:       "failed insert",
				Strain:         strains[0],
				GrainSubstrate: substrates[1],
				BulkSubstrate:  substrates[1],
			},
			// result: types.Lifecycle{Name: "failed insert"},
			err: fmt.Errorf("lifecycle was not added: 0"),
		},
		// // can't really test this anymore since the unique index changed to include
		// // ctime; just leaving it here so nobody tries to reimplement it
		// "unique_key_violation": {
		// 	lc: types.Lifecycle{
		// 		Location:       "reference implementation",
		// 		Strain:         strains[0],
		// 		GrainSubstrate: substrates[1],
		// 		BulkSubstrate:  substrates[2],
		// 	},
		// 	// result: types.Lifecycle{Name: "failed insert"},
		// 	err: fmt.Errorf(uniqueKeyViolation, "lifecycles_name_key"),
		// },
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
			// require.Equal(t, v.result.Name, result.Name)
			// require.Equal(t, v.result.Location, result.Location)
			// require.Equal(t, v.result.GrainCost, result.GrainCost)
			// require.Equal(t, v.result.BulkCost, result.BulkCost)
			// require.Equal(t, v.result.Yield, result.Yield)
			// require.Equal(t, v.result.Count, result.Count)
			// require.Equal(t, v.result.Gross, result.Gross)
			// require.Equal(t, v.result.Strain.Name, result.Strain.Name)
			// require.Equal(t, v.result.GrainSubstrate.Name, result.GrainSubstrate.Name)
			// require.Equal(t, v.result.BulkSubstrate.Name, result.BulkSubstrate.Name)
			// require.Less(t, epoch, v.result.MTime)
			// require.Less(t, epoch, v.result.CTime)
		})
	}
}

func Test_UpdateLifecycle(t *testing.T) {
	t.Parallel()

	updated, err := db.SelectLifecycle(context.Background(), "update me!", "Test_UpdateLifecycle")
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
			err: fmt.Errorf("lifecycle was not updated"),
		},
		"no_rows_affected_grain": {
			xform: func(lc types.Lifecycle) types.Lifecycle {
				lc.GrainSubstrate = types.Substrate{UUID: "missing"}
				return lc
			},
			err: fmt.Errorf("lifecycle was not updated"),
		},
		"no_rows_affected_bulk": {
			xform: func(lc types.Lifecycle) types.Lifecycle {
				lc.BulkSubstrate = types.Substrate{UUID: "missing"}
				return lc
			},
			err: fmt.Errorf("lifecycle was not updated"),
		},
		"check_grain_type": {
			xform: func(lc types.Lifecycle) types.Lifecycle {
				lc.GrainSubstrate = substrates[2]
				return lc
			},
			err: fmt.Errorf("lifecycle was not updated"),
		},
		"check_bulk_type": {
			xform: func(lc types.Lifecycle) types.Lifecycle {
				lc.BulkSubstrate = substrates[1]
				return lc
			},
			err: fmt.Errorf("lifecycle was not updated"),
		},
		"unique_key_violation": {
			xform: func(lc types.Lifecycle) types.Lifecycle {
				lc.Location = "reference implementation"
				return lc
			},
			err: fmt.Errorf(uniqueKeyViolation, "lifecycles_location_ctime_key"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			lc := v.xform(updated)
			err := db.UpdateLifecycle(context.Background(), lc, types.CID(k))
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
			err: fmt.Errorf(foreignKeyViolation1toMany, "lifecycles", "events_lifecycle_uuid_fkey", "events"),
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
