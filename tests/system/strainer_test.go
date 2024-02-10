package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

func Test_SelectAllStrains(t *testing.T) {
	set := map[string]struct {
		result []types.Strain
		err    error
	}{
		"happy_path": {
			result: []types.Strain{
				{UUID: "0", Name: "Morel", Vendor: types.Vendor{}},
				{UUID: "1", Name: "Hens o'' the Wood", Vendor: types.Vendor{}},
			},
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			result, err := db.SelectAllStrains(context.Background(), types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}

func Test_SelectStrain(t *testing.T) {
	set := map[string]struct {
		id     types.UUID
		result types.Strain
		err    error
	}{
		"happy_path": {
			id:     "0",
			result: types.Strain{UUID: "0", Name: "Morel", Vendor: types.Vendor{}},
		},
		"no_results_found": {
			id:  "foobar",
			err: noRows,
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			result, err := db.SelectStrain(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}

func Test_InsertStrain(t *testing.T) {
	set := map[string]struct {
		s   types.Strain
		err error
	}{
		"happy_path": {
			s: types.Strain{Name: "ubermyc", Vendor: types.Vendor{UUID: "0"}},
		},
		"no_rows_affected": {
			s: types.Strain{Name: "ubermyc", Vendor: types.Vendor{UUID: "foobar"}},
		},
		"unique_key_violation": {
			s: types.Strain{Name: "Morel", Vendor: types.Vendor{UUID: "0"}},
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			result, err := db.InsertStrain(context.Background(), v.s, types.CID(k))
			require.Equal(t, v.err, err)
			require.NotEmpty(t, result)
		})
	}
}

func Test_UpdateStrain(t *testing.T) {
	set := map[string]struct {
		id  types.UUID
		s   types.Strain
		err error
	}{
		"happy_path": {
			id: "1",
			s:  types.Strain{Name: "Chicken o' the Wood", Vendor: types.Vendor{UUID: "0"}},
		},
		"no_rows_affected_strain": {
			id: "foobar",
			s:  types.Strain{Name: "Chicken o' the Wood", Vendor: types.Vendor{UUID: "0"}},
		},
		"no_rows_affected_vendor": {
			id: "0",
			s:  types.Strain{Name: "Chicken o' the Wood", Vendor: types.Vendor{UUID: "foobar"}},
		},
		"unique_key_violation": {
			id: "1",
			s:  types.Strain{Name: "Morel", Vendor: types.Vendor{UUID: "0"}},
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			err := db.UpdateStrain(context.Background(), v.id, v.s, types.CID(k))
			require.Equal(t, v.err, err)
		})
	}
}

func Test_DeleteStrain(t *testing.T) {
	set := map[string]struct {
		id  types.UUID
		err error
	}{
		"happy_path": {
			id: "-1",
		},
		"no_rows_affected": {
			id:  "foobar",
			err: fmt.Errorf("strain could not be deleted 'foobar'"),
		},
		"referential_violation": {
			id:  "0",
			err: fmt.Errorf("referential constraint"),
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			err := db.DeleteStrain(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err)
		})
	}
}
