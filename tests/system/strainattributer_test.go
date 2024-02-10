package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

func Test_KnownAttributeNames(t *testing.T) {
	set := map[string]struct {
		result []string
		err    error
	}{
		"happy_path": {
			result: []string{"contamination resistance", "headroom", "color"},
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
			s: types.Strain{UUID: "0"},
			result: []types.StrainAttribute{
				{UUID: "0", Name: "contamination resistance", Value: "high"},
				{UUID: "1", Name: "headroom (cm)", Value: "25"},
			},
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			err := db.GetAllAttributes(context.Background(), &v.s, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.s.Attributes, v.result)
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
			s:      types.Strain{UUID: "0"},
			n:      "new name",
			v:      "new value",
			result: 1,
		},
		"no_rows_affected": {
			s:   types.Strain{UUID: "-2"},
			n:   "new name",
			v:   "new value",
			err: fmt.Errorf("attribute was not added"),
		},
		"unique_key_violation": {
			s:   types.Strain{UUID: "0"},
			n:   "contamination resistance",
			v:   "new value",
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
	set := map[string]struct {
		s      types.Strain
		n, v   string
		result []types.StrainAttribute
		err    error
	}{
		"happy_path": {
			s: types.Strain{UUID: "0", Attributes: []types.StrainAttribute{
				{UUID: "0", Name: "contamination resistance", Value: "high"},
				{UUID: "1", Name: "headroom (cm)", Value: "25"},
			}},
			n: "contamination resistance",
			v: "sterile",
			result: []types.StrainAttribute{
				{UUID: "0", Name: "contamination resistance", Value: "sterile"},
				{UUID: "1", Name: "headroom (cm)", Value: "25"},
			},
		},
		"no_rows_affected": {
			s: types.Strain{UUID: "0", Attributes: []types.StrainAttribute{
				{UUID: "0", Name: "contamination resistance", Value: "high"},
				{UUID: "1", Name: "headroom (cm)", Value: "25"},
			}},
			n: "effervescence",
			v: "fuzzy",
			result: []types.StrainAttribute{
				{UUID: "0", Name: "contamination resistance", Value: "high"},
				{UUID: "1", Name: "headroom (cm)", Value: "25"},
			},
			err: fmt.Errorf("attribute was not changed"),
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			err := db.ChangeAttribute(context.Background(), &v.s, v.n, v.v, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.s.Attributes, v.result)
		})
	}
}

func Test_RemoveAttribute(t *testing.T) {
	set := map[string]struct {
		s      types.Strain
		id     types.UUID
		result []types.StrainAttribute
		err    error
	}{
		"happy_path": {
			s: types.Strain{UUID: "0", Attributes: []types.StrainAttribute{
				{UUID: "0"},
				{UUID: "1"},
				{UUID: "2"},
			}},
			id: "1",
			result: []types.StrainAttribute{
				{UUID: "0"},
				{UUID: "2"},
			},
		},
		"no_rows_affected": {
			s: types.Strain{UUID: "0", Attributes: []types.StrainAttribute{
				{UUID: "0"},
				{UUID: "2"},
			}},
			id: "-2",
			result: []types.StrainAttribute{
				{UUID: "0"},
				{UUID: "2"},
			},
			err: fmt.Errorf("attribute was not removed"),
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			err := db.RemoveAttribute(context.Background(), &v.s, v.id, types.CID(k))
			require.Equal(t, v.err, err)
		})
	}
}
