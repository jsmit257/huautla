package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jsmit257/huautla/types"
)

func Test_GetPhotos(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     types.UUID
		photos int
		err    error
	}{
		"happy_path": {
			id:     "generation photo",
			photos: 2,
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.GetPhotos(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.photos, len(result))
		})
	}
}

func Test_AddPhoto(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id    types.UUID
		p     types.Photo
		count int
		err   error
	}{
		"happy_path": {
			id:    "add photo event 0",
			p:     types.Photo{UUID: "inserted", Filename: "filename.png"},
			count: 1,
		},
		"no_rows_affected_photo": {
			id:    "missing",
			count: 0,
			err:   fmt.Errorf("pq: foreign key violation"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			photos, err := db.AddPhoto(
				context.Background(),
				v.id,
				[]types.Photo{},
				v.p,
				types.CID(k))
			equalErrorMessages(t, v.err, err)
			require.Equalf(t, v.count, len(photos), "actual: %#v", photos)
		})
	}
}

func Test_ChangePhoto(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		n      types.Photo
		result []types.Photo
		err    error
	}{
		"happy_path": { // happy path needs to run first, synchronously
			n: types.Photo{
				UUID:     "gen photo 2",
				Filename: "updated!",
			},
			result: []types.Photo{{UUID: "gen photo 2", Filename: "updated!"}},
		},
		"no_rows_affected": { // dunno how this would happen, but whatever
			n: types.Photo{
				UUID: "missing",
			},
			result: []types.Photo{{UUID: "gen photo 2", Filename: "gen photo 2"}},
			err:    fmt.Errorf("photo was not changed"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			photos, err := db.ChangePhoto(
				context.Background(),
				[]types.Photo{{UUID: "gen photo 2", Filename: "gen photo 2"}},
				v.n,
				types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result[0].Filename, photos[0].Filename)
		})
	}
}

func Test_RemovePhoto(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     types.UUID
		result []types.Photo
		err    error
	}{
		"no_rows_affected": {
			id: "missing",
			result: []types.Photo{
				{UUID: "photo 2"},
				{UUID: "gen photo 2"},
			},
			err: fmt.Errorf("photo could not be removed"),
		},
		"foreign_key": {
			id: "gen photo 2",
			result: []types.Photo{
				{UUID: "photo 2"},
				{UUID: "gen photo 2"},
			},
			err: fmt.Errorf("pq: foreign key violation"),
		},
		"happy_path": {
			id:     "photo 2",
			result: []types.Photo{{UUID: "gen photo 2"}},
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			photos, err := db.RemovePhoto(
				context.Background(),
				[]types.Photo{
					{UUID: "photo 2"},
					{UUID: "gen photo 2"},
				},
				v.id,
				types.CID(k))
			equalErrorMessages(t, v.err, err)
			require.Equal(t, v.result, photos)
		})
	}
}
