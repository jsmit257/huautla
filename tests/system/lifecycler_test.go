package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

func Test_SelectLifecycle(t *testing.T) {
	set := map[string]struct {
		id     types.UUID
		result types.Lifecycle
		err    error
	}{
		"happy_path": {
			id: "0",
			result: types.Lifecycle{
				Name:      "reference implementation",
				Location:  "testing",
				GrainCost: 0,
				BulkCost:  0,
				Yield:     0,
				Count:     0,
				Gross:     0,
				MTime:     epoch,
				CTime:     epoch,
				Strain: types.Strain{
					UUID:   "0",
					Name:   "Morel",
					Vendor: vendor0,
				},
				GrainSubstrate: types.Substrate{
					UUID:   "0",
					Name:   "Rye",
					Type:   "Grain",
					Vendor: vendor0,
				},
				BulkSubstrate: types.Substrate{
					UUID:   "2",
					Name:   "Cedar chips",
					Type:   "Bulk",
					Vendor: vendor0,
				},
			},
		},
		"no_rows_returned": {
			id:     "foobar",
			result: types.Lifecycle{UUID: "0"},
			err:    noRows,
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			result, err := db.SelectLifecycle(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result.Name, result.Name)
		})
	}
}
func Test_InsertLifecycle(t *testing.T) {
	set := map[string]struct {
		lc     types.Lifecycle
		result types.Lifecycle
		err    error
	}{
		"happy_path": {
			lc: types.Lifecycle{
				Name:           "inserted record",
				Location:       "testing",
				GrainCost:      1,
				BulkCost:       2,
				Yield:          3,
				Count:          4,
				Gross:          5,
				Strain:         types.Strain{UUID: "1"},
				GrainSubstrate: types.Substrate{UUID: "1"},
				BulkSubstrate:  types.Substrate{UUID: "2"},
			},
			result: types.Lifecycle{
				Name:      "inserted record",
				Location:  "testing",
				GrainCost: 1,
				BulkCost:  2,
				Yield:     3,
				Count:     4,
				Gross:     5,
				Strain: types.Strain{
					UUID:   "1",
					Name:   "Hens o' the Wood",
					Vendor: vendor0,
				},
				GrainSubstrate: types.Substrate{
					UUID:   "1",
					Name:   "Millet",
					Type:   "Grain",
					Vendor: vendor0,
				},
				BulkSubstrate: types.Substrate{
					UUID:   "2",
					Name:   "Cedar chips",
					Type:   "Bulk",
					Vendor: vendor0,
				},
			},
		},
		"no_rows_affected_grain": {
			lc: types.Lifecycle{
				Name:           "failed insert",
				Location:       "testing",
				Strain:         types.Strain{UUID: "1"},
				GrainSubstrate: types.Substrate{UUID: "foobar"},
				BulkSubstrate:  types.Substrate{UUID: "2"},
			},
			err: fmt.Errorf("lifecycle was not added"),
		},
		"no_rows_affected_bulk": {
			lc: types.Lifecycle{
				Name:           "failed insert",
				Location:       "testing",
				Strain:         types.Strain{UUID: "1"},
				GrainSubstrate: types.Substrate{UUID: "1"},
				BulkSubstrate:  types.Substrate{UUID: "foobar"},
			},
			err: fmt.Errorf("lifecycle was not added"),
		},
		"no_rows_affected_strain": {
			lc: types.Lifecycle{
				Name:           "failed insert",
				Location:       "testing",
				Strain:         types.Strain{UUID: "foobar"},
				GrainSubstrate: types.Substrate{UUID: "1"},
				BulkSubstrate:  types.Substrate{UUID: "2"},
			},
			err: fmt.Errorf("lifecycle was not added"),
		},
		"no_rows_affected_check_grain_type": {
			lc: types.Lifecycle{
				Name:           "failed insert",
				Location:       "testing",
				Strain:         types.Strain{UUID: "0"},
				GrainSubstrate: types.Substrate{UUID: "2"},
				BulkSubstrate:  types.Substrate{UUID: "2"},
			},
			err: fmt.Errorf("lifecycle was not added"),
		},
		"no_rows_affected_check_bulk_type": {
			lc: types.Lifecycle{
				Name:           "failed insert",
				Location:       "testing",
				Strain:         types.Strain{UUID: "0"},
				GrainSubstrate: types.Substrate{UUID: "1"},
				BulkSubstrate:  types.Substrate{UUID: "1"},
			},
			err: fmt.Errorf("lifecycle was not added"),
		},
		"unique_key_violation": {
			lc: types.Lifecycle{
				Name:           "reference implementation",
				Location:       "testing",
				Strain:         types.Strain{UUID: "1"},
				GrainSubstrate: types.Substrate{UUID: "1"},
				BulkSubstrate:  types.Substrate{UUID: "2"},
			},
			err: fmt.Errorf("lifecycle was not added"),
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			result, err := db.InsertLifecycle(context.Background(), v.lc, types.CID(k))
			require.Equal(t, v.err, err)
			require.NotEmpty(t, result.UUID)
			require.Equal(t, v.result.Name, result.Name)
			require.Equal(t, v.result.Location, result.Location)
			require.Equal(t, v.result.GrainCost, result.GrainCost)
			require.Equal(t, v.result.BulkCost, result.BulkCost)
			require.Equal(t, v.result.Yield, result.Yield)
			require.Equal(t, v.result.Count, result.Count)
			require.Equal(t, v.result.Gross, result.Gross)
			require.Equal(t, v.result.Strain.Name, result.Strain.Name)
			require.Equal(t, v.result.GrainSubstrate.Name, result.GrainSubstrate.Name)
			require.Equal(t, v.result.BulkSubstrate.Name, result.BulkSubstrate.Name)
			require.LessOrEqual(t, epoch, v.result.MTime)
			require.LessOrEqual(t, epoch, v.result.CTime)
		})
	}
}
func Test_UpdateLifecycle(t *testing.T) {
	set := map[string]struct {
		lc  types.Lifecycle
		err error
	}{
		"happy_path": {
			lc: types.Lifecycle{
				UUID:           "1",
				Name:           "updated record",
				Location:       "testing2",
				GrainCost:      5,
				BulkCost:       4,
				Yield:          3,
				Count:          2,
				Gross:          1,
				Strain:         types.Strain{UUID: "0"},
				GrainSubstrate: types.Substrate{UUID: "0"},
				BulkSubstrate:  types.Substrate{UUID: "2"},
			},
		},
		"no_rows_affected_grain": {
			lc: types.Lifecycle{
				UUID:           "1",
				Name:           "failed update",
				Location:       "testing",
				Strain:         types.Strain{UUID: "1"},
				GrainSubstrate: types.Substrate{UUID: "foobar"},
				BulkSubstrate:  types.Substrate{UUID: "2"},
			},
			err: fmt.Errorf("lifecycle was not added"),
		},
		"no_rows_affected_bulk": {
			lc: types.Lifecycle{
				UUID:           "1",
				Name:           "failed update",
				Location:       "testing",
				Strain:         types.Strain{UUID: "1"},
				GrainSubstrate: types.Substrate{UUID: "1"},
				BulkSubstrate:  types.Substrate{UUID: "foobar"},
			},
			err: fmt.Errorf("lifecycle was not added"),
		},
		"no_rows_affected_strain": {
			lc: types.Lifecycle{
				UUID:           "1",
				Name:           "failed update",
				Location:       "testing",
				Strain:         types.Strain{UUID: "foobar"},
				GrainSubstrate: types.Substrate{UUID: "1"},
				BulkSubstrate:  types.Substrate{UUID: "2"},
			},
			err: fmt.Errorf("lifecycle was not added"),
		},
		"no_rows_affected_check_grain_type": {
			lc: types.Lifecycle{
				UUID:           "1",
				Name:           "failed update",
				Strain:         types.Strain{UUID: "0"},
				GrainSubstrate: types.Substrate{UUID: "2"},
				BulkSubstrate:  types.Substrate{UUID: "2"},
			},
			err: fmt.Errorf("lifecycle was not added"),
		},
		"no_rows_affected_check_bulk_type": {
			lc: types.Lifecycle{
				UUID:           "1",
				Name:           "failed update",
				Strain:         types.Strain{UUID: "0"},
				GrainSubstrate: types.Substrate{UUID: "1"},
				BulkSubstrate:  types.Substrate{UUID: "1"},
			},
			err: fmt.Errorf("lifecycle was not added"),
		},
		"unique_key_violation": {
			lc: types.Lifecycle{
				UUID:           "1",
				Name:           "reference implementation",
				Strain:         types.Strain{UUID: "1"},
				GrainSubstrate: types.Substrate{UUID: "1"},
				BulkSubstrate:  types.Substrate{UUID: "2"},
			},
			err: fmt.Errorf("lifecycle was not added"),
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			err := db.UpdateLifecycle(context.Background(), v.lc, types.CID(k))
			require.Equal(t, v.err, err)
		})
	}
}
func Test_DeleteLifecycle(t *testing.T) {
	set := map[string]struct {
		id  types.UUID
		err error
	}{
		"happy_path": {
			id: "-1",
		},
		"no_rows_affected": {
			id:  "foobar",
			err: fmt.Errorf("lifecycle could not be deleted 'foobar'"),
		},
		"referential_violation": {
			id:  "0",
			err: fmt.Errorf("lifecycle could not be deleted 'foobar'"),
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			err := db.DeleteLifecycle(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err)
		})
	}
}
