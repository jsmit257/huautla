package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jsmit257/huautla/types"
)

func Test_GetNotes(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id    types.UUID
		count int
		err   error
	}{
		"happy_path": {
			id:    "notable lifecycle",
			count: 2,
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.GetNotes(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.count, len(result))
		})
	}
}

func Test_AddNote(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id    types.UUID
		n     types.Note
		count int
		err   error
	}{
		"happy_path": {
			id:    "insert notable",
			n:     types.Note{Note: "inserted"},
			count: 2,
		},
		"no_rows_affected_note": {
			id:    "missing",
			n:     types.Note{Note: "inserted"},
			count: 1,
			err:   fmt.Errorf("note was not added"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			notes, err := db.AddNote(context.Background(), v.id, []types.Note{{UUID: "notable lifecycle 2"}}, v.n, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equalf(t, v.count, len(notes), "actual: %#v", notes)
		})
	}
}

func Test_ChangeNote(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		n      types.Note
		result []types.Note
		err    error
	}{
		"happy_path": { // happy path needs to run first, synchronously
			n: types.Note{
				UUID: "notable generation 2",
				Note: "updated!",
			},
			result: []types.Note{{UUID: "notable generation 2", Note: "updated!"}},
		},
		"no_rows_affected": { // dunno how this would happen, but whatever
			n: types.Note{
				UUID: "missing",
			},
			result: []types.Note{{UUID: "notable generation 2", Note: "notable generation 2"}},
			err:    fmt.Errorf("note was not changed"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			notes, err := db.ChangeNote(
				context.Background(),
				[]types.Note{{UUID: "notable generation 2", Note: "notable generation 2"}},
				v.n,
				types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result[0].Note, notes[0].Note)
		})
	}
}

func Test_RemoveNote(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     types.UUID
		result []types.Note
		err    error
	}{
		"no_rows_affected": {
			id:     "missing",
			result: []types.Note{{UUID: "notable lifecycle 2"}},
			err:    fmt.Errorf("note could not be removed"),
		},
		"happy_path": {
			id:     "notable lifecycle 2",
			result: []types.Note{},
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			notes, err := db.RemoveNote(
				context.Background(),
				[]types.Note{{UUID: "notable lifecycle 2"}},
				v.id,
				types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, notes)
		})
	}
}
