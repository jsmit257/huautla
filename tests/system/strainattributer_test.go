package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

var strainattributes = []types.StrainAttribute{
	{UUID: "0", Name: "contamination resistance", Value: "high"},
	{UUID: "1", Name: "headroom (cm)", Value: "25"},
	{UUID: "2", Name: "color", Value: "purple"},
}

func Test_KnownAttributeNames(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		result []string
		err    error
	}{
		"happy_path": {
			result: []string{"color", "contamination resistance", "headroom (cm)"},
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.KnownAttributeNames(context.Background(), types.CID(k))
			require.Equal(t, v.err, err)
			require.Subset(t, result, v.result, "result: '%q', expected: '%q'", result, v.result)
		})
	}
}

func Test_GetAllAttributes(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		s      types.Strain
		result []types.StrainAttribute
		err    error
	}{
		"happy_path": {
			s:      strains[0],
			result: []types.StrainAttribute{strainattributes[0], strainattributes[1]},
		},
		"no_rows_returned": {
			s:      types.Strain{UUID: "missing"},
			result: []types.StrainAttribute{},
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.GetAllAttributes(context.Background(), &v.s, types.CID(k))
			require.Equal(t, v.err, err)
			require.ElementsMatch(t, v.result, v.s.Attributes)
		})
	}
}

func Test_AddAttribute(t *testing.T) {
	t.Parallel()

	strain, err := db.SelectStrain(context.Background(), "add attribute", "Test_AddAttribute")
	require.Nil(t, err)

	set := map[string]struct {
		s      types.Strain
		a      types.StrainAttribute
		result int
		err    error
	}{
		"happy_path": {
			s:      strain,
			a:      types.StrainAttribute{Name: "new name", Value: "new value"},
			result: 2,
		},
		"no_rows_affected": {
			s:   types.Strain{UUID: "missing"},
			err: fmt.Errorf("attribute was not added"),
		},
		"unique_key_violation": {
			s:      strain,
			a:      types.StrainAttribute{Name: "existing"},
			result: 1,
			err:    fmt.Errorf(uniqueKeyViolation, "strain_attributes_name_strain_uuid_key"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			a, err := db.AddAttribute(context.Background(), &v.s, v.a, types.CID(k))
			equalErrorMessages(t, v.err, err)
			require.Equal(t, v.result, len(v.s.Attributes))
			require.NotEmpty(t, a)
		})
	}
}

func Test_ChangeAttribute(t *testing.T) {
	t.Parallel()

	strain, err := db.SelectStrain(context.Background(), "change attribute", "Test_ChangeAttribute")
	require.Nil(t, err)

	set := map[string]struct {
		a      types.StrainAttribute
		result []types.StrainAttribute
		err    error
	}{
		"happy_path": { // run this first, synchronously
			a: types.StrainAttribute{
				UUID:  strain.Attributes[0].UUID,
				Name:  strain.Attributes[0].Name,
				Value: "malabar",
			},
			result: []types.StrainAttribute{
				func() types.StrainAttribute {
					result := strain.Attributes[0]
					result.Value = "malabar"
					return result
				}(),
			},
		},
		"no_rows_affected": {
			a:      types.StrainAttribute{Name: "effervescence", Value: "fuzzy"},
			result: strain.Attributes[:],
			err:    fmt.Errorf("attribute was not changed"),
		},
	}
	for k, v := range set {
		k, v, strain := k, v, strain
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.ChangeAttribute(context.Background(), &strain, v.a, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, strain.Attributes)
		})
	}
}

func Test_RemoveAttribute(t *testing.T) {
	t.Parallel()

	strain, err := db.SelectStrain(context.Background(), "remove attribute", "Test_RemoveAttribute")
	require.Nil(t, err)

	err = db.GetAllAttributes(context.Background(), &strain, "Test_RemoveAttribute")
	require.Nil(t, err)

	set := map[string]struct {
		id     types.UUID
		result []types.StrainAttribute
		err    error
	}{
		"happy_path": {
			id:     strain.Attributes[1].UUID,
			result: []types.StrainAttribute{strain.Attributes[0], strain.Attributes[2]},
		},
		"no_rows_affected": {
			id:     "missing",
			result: strain.Attributes[:],
			err:    fmt.Errorf("attribute was not removed"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.RemoveAttribute(context.Background(), &strain, v.id, types.CID(k))
			require.Equal(t, v.err, err)
		})
	}
}
