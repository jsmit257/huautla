package data

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/jsmit257/huautla/types"
)

var (
	_photos = []photo{
		{UUID: "id-0", Filename: "photo 0", MTime: wwtbn, CTime: wwtbn},
		{UUID: "id-1", Filename: "photo 1", MTime: wwtbn, CTime: wwtbn, Notes: []types.Note{
			types.Note(_notes[0]),
			types.Note(_notes[2]),
		}},
		{UUID: "id-2", Filename: "photo 2", MTime: wwtbn, CTime: wwtbn},
		{UUID: "id-3", Filename: "photo 3", MTime: wwtbn, CTime: wwtbn, Owner: &types.PhotoOwner{
			ParentType: "generation",
			OwnerUUID:  "genowner",
			ParentUUID: "genuuid3",
			Label:      "sourcestrain->genowner",
		}},
		{UUID: "id-4", Filename: "photo 4", MTime: wwtbn, CTime: wwtbn, Owner: &types.PhotoOwner{
			ParentType: "lifecycle",
			OwnerUUID:  "lcowner",
			ParentUUID: "lcuuid4",
			Label:      "lcname->lcowner",
		}},
		{UUID: "id-5", Filename: "photo 5", MTime: wwtbn, CTime: wwtbn, Owner: &types.PhotoOwner{
			ParentType: "strain",
			OwnerUUID:  "strainowner",
			Label:      "strain->genowner",
		}},
	}
	photoFields = row{
		"id",
		"filename",
		"mtime",
		"ctime",
		"note_uuid",
		"note",
		"note_mtime",
		"note_ctime",
	}
	photoValues = [][]driver.Value{
		{_photos[0].UUID, _photos[0].Filename, _photos[0].MTime, _photos[0].CTime, nil, nil, nil, nil},
		{_photos[1].UUID, _photos[1].Filename, _photos[1].MTime, _photos[1].CTime, _photos[1].Notes[0].UUID, _photos[1].Notes[0].Note, _photos[1].Notes[0].MTime, _photos[1].Notes[0].CTime},
		{_photos[1].UUID, _photos[1].Filename, _photos[1].MTime, _photos[1].CTime, _photos[1].Notes[1].UUID, _photos[1].Notes[1].Note, _photos[1].Notes[1].MTime, _photos[1].Notes[1].CTime},
		{_photos[2].UUID, _photos[2].Filename, _photos[2].MTime, _photos[2].CTime, nil, nil, nil, nil},
	}
	allPhotoFields = row{
		"id",
		"filename",
		"mtime",
		"ctime",
		"parent_type",
		"owner_uuid",
		"parent_uuid",
		"label",
	}
	allPhotoValues = [][]driver.Value{
		{_photos[3].UUID, _photos[3].Filename, _photos[3].MTime, _photos[3].CTime, _photos[3].Owner.ParentType, _photos[3].Owner.OwnerUUID, _photos[3].Owner.ParentUUID, _photos[3].Owner.Label},
		{_photos[4].UUID, _photos[4].Filename, _photos[4].MTime, _photos[4].CTime, _photos[4].Owner.ParentType, _photos[4].Owner.OwnerUUID, _photos[4].Owner.ParentUUID, _photos[4].Owner.Label},
		{_photos[5].UUID, _photos[5].Filename, _photos[5].MTime, _photos[5].CTime, _photos[5].Owner.ParentType, _photos[5].Owner.OwnerUUID, _photos[5].Owner.ParentUUID, _photos[5].Owner.Label},
	}
)

func Test_AllPhotos(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "Test_AllPhotos")

	set := map[string]struct {
		db     func() *sql.DB
		result []types.Photo
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				allPhotoFields.mock(mock, allPhotoValues...)
				return db
			},
			result: []types.Photo{
				types.Photo(_photos[3]),
				types.Photo(_photos[4]),
				types.Photo(_photos[5]),
			},
		},
		"db_error": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			result: []types.Photo{},
			err:    fmt.Errorf("some error"),
		},
	}

	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			result, err := (&Conn{
				query:        v.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", k),
			}).AllPhotos(context.Background(), "Test_AllPhotos")

			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}

func Test_GetPhotos(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "Test_GetPhotos")

	set := map[string]struct {
		db     func() *sql.DB
		result []types.Photo
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				photoFields.mock(mock, photoValues...)
				return db
			},
			result: []types.Photo{
				types.Photo(_photos[0]),
				types.Photo(_photos[1]),
				types.Photo(_photos[2]),
			},
		},
		"db_error": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			result: []types.Photo{},
			err:    fmt.Errorf("some error"),
		},
	}

	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			result, err := (&Conn{
				query:        v.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", k),
			}).GetPhotos(context.Background(), "0", "Test_GetPhotos")

			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}

func Test_AddPhoto(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "Test_AddPhoto")

	set := map[string]struct {
		db     func() *sql.DB
		result int
		n      types.Photo
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			result: 1,
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("photo was not added"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			err: fmt.Errorf("some error"),
		},
	}

	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			result, err := (&Conn{
				query:        v.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", k),
			}).AddPhoto(context.Background(), "0", []types.Photo{}, v.n, "Test_AddPhoto")

			require.Equal(t, v.err, err)
			require.Equal(t, v.result, len(result))
		})
	}
}

func Test_ChangePhoto(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "Test_ChangePhoto")

	set := map[string]struct {
		db  func() *sql.DB
		p   types.Photo
		id  types.UUID
		err error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			p: types.Photo{UUID: "0", Filename: "photo"},
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("photo was not changed"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			err: fmt.Errorf("some error"),
		},
	}

	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			_, err := (&Conn{
				query:        v.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", k),
			}).ChangePhoto(
				context.Background(),
				[]types.Photo{
					{UUID: "1"},
					{UUID: "0"},
				},
				v.p,
				"Test_ChangePhoto")

			require.Equal(t, v.err, err)
		})
	}
}

func Test_RemovePhoto(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "Test_RemovePhoto")

	set := map[string]struct {
		db     func() *sql.DB
		result []types.Photo
		id     types.UUID
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			result: []types.Photo{{UUID: "1"}},
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			result: []types.Photo{
				{UUID: "1"},
				{UUID: "0"},
			},
			err: fmt.Errorf("photo could not be removed"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			result: []types.Photo{
				{UUID: "1"},
				{UUID: "0"},
			},
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			result: []types.Photo{
				{UUID: "1"},
				{UUID: "0"},
			},
			err: fmt.Errorf("some error"),
		},
	}

	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			result, err := (&Conn{
				query:        v.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", k),
			}).RemovePhoto(
				context.Background(),
				[]types.Photo{
					{UUID: "1"},
					{UUID: "0"},
				},
				"0",
				"Test_RemovePhoto")

			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}
