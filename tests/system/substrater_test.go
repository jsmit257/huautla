package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

var substrates = []types.Substrate{
	{UUID: "0", Name: "Rye", Type: "Grain", Vendor: vendor0, Ingredients: ingredients},
	{UUID: "1", Name: "Millet", Type: "Grain", Vendor: vendor0, Ingredients: ingredients},
	{UUID: "2", Name: "Cedar chips", Type: "Bulk", Vendor: vendor0, Ingredients: ingredients},
}

func Test_SelectAllSubstrates(t *testing.T) {
	set := map[string]struct {
		result []types.Substrate
		err    error
	}{
		"happy_path": {
			result: substrates,
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			result, err := db.SelectAllSubstrates(context.Background(), types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result[0:len(v.result)])
		})
	}
}

func Test_SelectSubstrate(t *testing.T) {
	t.Parallel()
	set := map[string]struct {
		id     types.UUID
		result types.Substrate
		err    error
	}{
		"happy_path": {
			id:     substrates[0].UUID,
			result: substrates[0],
		},
		"no_rows_returned": {
			id:  "missing",
			err: noRows,
		},
		"query_fails": {
			id:  invalidUUID,
			err: noRows,
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.SelectSubstrate(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}

func Test_InsertSubstrate(t *testing.T) {
	t.Parallel()
	set := map[string]struct {
		s   types.Substrate
		err error
	}{
		"happy_path": {
			s: types.Substrate{Name: "Honey Solution", Type: "Bulk", Vendor: vendor0},
		},
		"unique_key_violation": {
			s:   substrates[0],
			err: fmt.Errorf("duplicate key violation"),
		},
		"check_constraint_violation": {
			s:   types.Substrate{Name: "Maltodexterin", Type: "Stardust", Vendor: vendor0},
			err: fmt.Errorf("check constraint"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.InsertSubstrate(context.Background(), v.s, types.CID(k))
			require.Equal(t, v.err, err)
			require.NotEmpty(t, result.UUID)
		})
	}
}

func Test_UpdateSubstrate(t *testing.T) {
	t.Parallel()

	substrate, err := db.SelectSubstrate(context.Background(), "update me!", "Test_UpdateSubstrate")
	require.Nil(t, err)

	set := map[string]struct {
		s   types.Substrate
		err error
	}{
		"happy_path": {
			s: func() types.Substrate {
				result := substrate
				result.Name = "Updated"
				return result
			}(),
		},
		"unique_key_violation_name": {
			s: func() types.Substrate {
				result := substrate
				result.Name = "Millet"
				return result
			}(),
			err: fmt.Errorf("duplicate key violation"),
		},
		// "unique_key_violation_vendor": { // XXX: can't currently update vendpr
		// 	s: func() types.Substrate {
		// 		result := substrate
		// 		result.Vendor = types.Vendor{UUID: "missing"}
		// 		return result
		// 	}(),
		// 	err: fmt.Errorf("duplicate key violation"),
		// },
		// "check_constraint_violation": { // XXX: can't currently update type
		// 	s: func() types.Substrate {
		// 		result := substrate
		// 		result.Type = "Stardust"
		// 		return result
		// 	}(),
		// 	err: fmt.Errorf("check constraint"),
		// },
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.UpdateSubstrate(context.Background(), substrate.UUID, v.s, types.CID(k))
			require.Equal(t, v.err, err)
		})
	}
}

func Test_DeleteSubstrate(t *testing.T) {
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
			err: fmt.Errorf("vendor table was not deleted 'missing'"),
		},
		"query_fails": {
			id:  invalidUUID,
			err: fmt.Errorf("some error"),
		},
		"referential_violation": {
			id:  substrates[0].UUID,
			err: fmt.Errorf("referential constraint"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.DeleteSubstrate(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err)
		})
	}
}
