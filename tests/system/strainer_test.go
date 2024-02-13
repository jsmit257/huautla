package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

var strains = []types.Strain{
	{UUID: "0", Name: "Morel", Vendor: vendor0, Attributes: []types.StrainAttribute{
		strainattributes[0],
		strainattributes[1],
	}},
	{UUID: "1", Name: "Hens o'' the Wood", Vendor: vendor0},
}

func Test_SelectAllStrains(t *testing.T) {
	t.Skip() // selectAll doesn't fetch attributes like select, so strains var won't work for both
	t.Parallel()

	set := map[string]struct {
		result []types.Strain
		err    error
	}{
		"happy_path": {
			result: func() []types.Strain {
				result := []types.Strain{}
				for i, j := 0, len(strains); i < j; i++ {
					tmp := strains[i]
					tmp.Attributes = nil
					result = append(result, tmp)
				}
				return result
			}(),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.SelectAllStrains(context.Background(), types.CID(k))
			require.Equal(t, v.err, err)
			require.Subsetf(t, result, v.result, "wtf?!: \n'%q'\n'%q'", result, v.result)
		})
	}
}

func Test_SelectStrain(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     types.UUID
		result types.Strain
		err    error
	}{
		"happy_path": {
			id:     strains[0].UUID,
			result: strains[0],
		},
		"no_results_found": {
			id:     "missing",
			result: types.Strain{UUID: "missing"},
			err:    noRows,
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.SelectStrain(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}

func Test_InsertStrain(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		s   types.Strain
		err error
	}{
		"happy_path": {
			s: types.Strain{Name: "ubermyc", Vendor: vendor0},
		},
		"no_rows_affected": {
			s:   types.Strain{Name: "ubermyc", Vendor: types.Vendor{UUID: "missing"}},
			err: fmt.Errorf("strain was not added"),
		},
		"unique_key_violation": {
			s:   types.Strain{Name: "Morel", Vendor: vendor0},
			err: fmt.Errorf(uniqueKeyViolation, "strains_name_vendor_uuid_key"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.InsertStrain(context.Background(), v.s, types.CID(k))
			equalErrorMessages(t, v.err, err)
			require.NotEmpty(t, result.UUID)
		})
	}
}

func Test_UpdateStrain(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id  types.UUID
		s   types.Strain
		err error
	}{
		"happy_path": {
			id: "update me!",
			s:  types.Strain{Name: "Chicken o' the Wood", Vendor: vendor0},
		},
		"no_rows_affected_strain": {
			id:  "missing",
			s:   types.Strain{Name: "Chicken o' the Wood", Vendor: vendor0},
			err: fmt.Errorf("strain was not updated: 'missing'"),
		},
		// "no_rows_affected_vendor": {  // vendors aren't part of the update (yet)
		// 	id:  "update me!",
		// 	s:   types.Strain{Name: "Chicken o' the Wood", Vendor: types.Vendor{UUID: "missing"}},
		// 	err: fmt.Errorf("strain was not updated: 'update me!'"),
		// },
		"unique_key_violation": {
			id:  "update me!",
			s:   strains[0],
			err: fmt.Errorf(uniqueKeyViolation, "strains_name_vendor_uuid_key"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			// t.Parallel()
			err := db.UpdateStrain(context.Background(), v.id, v.s, types.CID(k))
			equalErrorMessages(t, v.err, err)
		})
	}
}

func Test_DeleteStrain(t *testing.T) {
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
			err: fmt.Errorf("strain could not be deleted: 'missing'"),
		},
		"referential_violation": {
			id: "0",
			err: fmt.Errorf(
				foreignKeyViolation1toMany,
				"strains",
				"strain_attributes_strain_uuid_fkey",
				"strain_attributes"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.DeleteStrain(context.Background(), v.id, types.CID(k))
			equalErrorMessages(t, v.err, err)
		})
	}
}
