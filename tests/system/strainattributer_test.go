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
	set := map[string]struct {
		result []string
		err    error
	}{
		"happy_path": {
			result: []string{"color", "contamination resistance", "headroom"}, // XXX: needs work
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			result, err := db.KnownAttributeNames(context.Background(), types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}

func Test_GetAllAttributes(t *testing.T) {
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
			s:   types.Strain{UUID: "missing"},
			err: fmt.Errorf("sql: no rows in result set"),
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			err := db.GetAllAttributes(context.Background(), &v.s, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, v.s.Attributes)
		})
	}
}

func Test_AddAttribute(t *testing.T) {
	set := map[string]struct {
		s      types.Strain
		n, v   string
		result int
		err    error
	}{
		"happy_path": {
			s:      types.Strain{UUID: "add attribute"},
			n:      "new name",
			v:      "new value",
			result: 1,
		},
		"no_rows_affected": {
			s:   types.Strain{UUID: "missing"},
			err: fmt.Errorf("attribute was not added"),
		},
		"unique_key_violation": {
			s:   types.Strain{UUID: "add attribute"},
			n:   "contamination resistance",
			err: fmt.Errorf("attribute was not added"),
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			err := db.AddAttribute(context.Background(), &v.s, v.n, v.v, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, len(v.s.Attributes))
		})
	}
}

func Test_ChangeAttribute(t *testing.T) {
	t.Parallel()

	strain, err := db.SelectStrain(context.Background(), "change attribute", "Test_ChangeAttribute")
	require.Nil(t, err)

	err = db.GetAllAttributes(context.Background(), &strain, "Test_ChangeAttribute")
	require.Nil(t, err)

	set := map[string]struct {
		n, v   string
		result []types.StrainAttribute
		err    error
	}{
		"happy_path": { // run this first, synchronously
			n: strainattributes[1].Name,
			v: "malabar",
			result: []types.StrainAttribute{
				strainattributes[1],
				func() types.StrainAttribute {
					result := strainattributes[1]
					result.Value = "malabar"
					return result
				}(),
			},
		},
		"no_rows_affected": {
			n:      "effervescence",
			v:      "fuzzy",
			result: strain.Attributes[:],
			err:    fmt.Errorf("attribute was not changed"),
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			err := db.ChangeAttribute(context.Background(), &strain, v.n, v.v, types.CID(k))
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
		"no_rows_affected": {
			id:     "missing",
			result: strain.Attributes[:],
			err:    fmt.Errorf("attribute was not removed"),
		},
		"happy_path": {
			id:     strain.Attributes[1].UUID,
			result: []types.StrainAttribute{strain.Attributes[0], strain.Attributes[2]},
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			err := db.RemoveAttribute(context.Background(), &strain, v.id, types.CID(k))
			require.Equal(t, v.err, err)
		})
	}
}
