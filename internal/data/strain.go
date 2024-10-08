package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) SelectAllStrains(ctx context.Context, cid types.CID) ([]types.Strain, error) {
	var err error

	deferred, start, l := initAccessFuncs("SelectAllStrains", db.logger, "nil", cid)
	defer deferred(start, err, l)

	var rows *sql.Rows

	result := make([]types.Strain, 0, 100)

	rows, err = db.query.QueryContext(ctx, psqls["strain"]["select-all"])
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var generationID *types.UUID

	for rows.Next() {
		row := types.Strain{}

		if err = rows.Scan(
			&row.UUID,
			&row.Species,
			&row.Name,
			&row.CTime,
			&row.DTime,
			&row.Vendor.UUID,
			&row.Vendor.Name,
			&row.Vendor.Website,
			&generationID,
		); err != nil {
			break
		}

		if generationID != nil {
			row.Generation = &types.Generation{UUID: *generationID}
		}

		result = append(result, row)
	}

	return result, err
}

func (db *Conn) SelectStrain(ctx context.Context, id types.UUID, cid types.CID) (types.Strain, error) {
	var err error

	deferred, start, l := initAccessFuncs("SelectStrain", db.logger, id, cid)
	defer deferred(start, err, l)

	result := types.Strain{UUID: id}

	var generationID *types.UUID

	if err = db.
		QueryRowContext(ctx, psqls["strain"]["select"], id).
		Scan(
			&result.Species,
			&result.Name,
			&result.CTime,
			&result.DTime,
			&result.Vendor.UUID,
			&result.Vendor.Name,
			&result.Vendor.Website,
			&generationID,
		); err == nil {
		if generationID != nil {
			result.Generation = &types.Generation{UUID: *generationID}
		}

		err = db.GetAllAttributes(ctx, &result, cid)
	}

	return result, err
}

func (db *Conn) InsertStrain(ctx context.Context, s types.Strain, cid types.CID) (types.Strain, error) {
	var err error

	s.UUID = types.UUID(db.generateUUID().String())
	s.CTime = time.Now().UTC()

	deferred, start, l := initAccessFuncs("InsertStrain", db.logger, s.UUID, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, psqls["strain"]["insert"], s.UUID, s.Species, s.Name, s.CTime, s.Vendor.UUID)
	if err != nil {
		if isPrimaryKeyViolation(err) {
			return db.InsertStrain(ctx, s, cid) // FIXME: infinite loop?
		}
		return s, err
	} else if rows, err := result.RowsAffected(); err != nil {
		return s, err
	} else if rows != 1 { // most likely cause is a bad vendor.uuid
		return s, fmt.Errorf("strain was not added")
	}

	return s, err
}

func (db *Conn) UpdateStrain(ctx context.Context, id types.UUID, s types.Strain, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("UpdateStrain", db.logger, id, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, psqls["strain"]["update"], s.Species, s.Name, s.Vendor.UUID, id)
	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return fmt.Errorf("strain was not updated: '%s'", id)
	}

	return nil
}

func (db *Conn) DeleteStrain(ctx context.Context, id types.UUID, cid types.CID) error {
	// TODO: delete all attributes first
	return db.deleteByUUID(ctx, id, cid, "DeleteStrain", "strain", db.logger)
}

func (db *Conn) GeneratedStrain(ctx context.Context, id types.UUID, cid types.CID) (types.Strain, error) {
	var err error

	deferred, start, l := initAccessFuncs("GeneratedStrains", db.logger, id, cid)
	defer deferred(start, err, l)

	result := types.Strain{}

	return result, db.
		QueryRowContext(ctx, psqls["strain"]["generated-strain"], id).
		Scan(
			&result.UUID,
			&result.Species,
			&result.Name,
			&result.Vendor.UUID,
			&result.Vendor.Name,
			&result.Vendor.Website,
			&result.CTime,
		)
}

func (db *Conn) UpdateGeneratedStrain(ctx context.Context, gid *types.UUID, sid types.UUID, cid types.CID) error {
	var err error

	deferred, start, l := initAccessFuncs("UpdateGeneratedStrain", db.logger, sid, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, psqls["strain"]["update-gen-strain"], gid, sid)
	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return sql.ErrNoRows
	}

	return nil
}
