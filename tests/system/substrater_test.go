package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

func Test_SelectAllSubstrates(t *testing.T) {
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
}

func Test_SelectSubstrate(t *testing.T) {
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
}

func Test_InsertSubstrate(t *testing.T) {
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
}

func Test_UpdateSubstrate(t *testing.T) {
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
}

func Test_DeleteSubstrate(t *testing.T) {
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
}
