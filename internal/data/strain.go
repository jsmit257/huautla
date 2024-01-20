package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) SelectAllStrains(ctx context.Context, cid types.CID) ([]types.Strain, error) {
	var err error

	deferred, start, l := initVendorFuncs("SelectAllStrains", db.logger, err, "nil", cid)
	defer deferred(start, err, l)

	var rows *sql.Rows

	result := make([]types.Strain, 0, 100)

	rows, err = db.query.QueryContext(ctx, db.sql["strain"]["select-all"])
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		row := types.Strain{}
		rows.Scan(
			&row.UUID,
			&row.Name,
			&row.Vendor.UUID,
			&row.Vendor.Name)
		result = append(result, row)
	}

	return result, err
}

func (db *Conn) SelectStrain(ctx context.Context, id types.UUID, cid types.CID) (types.Strain, error) {
	var err error

	deferred, start, l := initVendorFuncs("SelectStrain", db.logger, err, id, cid)
	defer deferred(start, err, l)

	result := types.Strain{UUID: id}
	err = db.
		QueryRowContext(ctx, db.sql["strain"]["select"], id).
		Scan(
			&result.Name,
			&result.Vendor.UUID,
			&result.Vendor.Name)

	return result, err
}

func (db *Conn) InsertStrain(ctx context.Context, s types.Strain, cid types.CID) (types.Strain, error) {
	var err error

	s.UUID = types.UUID(db.generateUUID().String())

	deferred, start, l := initVendorFuncs("InsertStrain", db.logger, err, s.UUID, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, db.sql["strain"]["insert"], s.UUID, s.Name, s.Vendor.UUID)
	if err != nil {
		// FIXME: choose what to do based on the tupe of error
		duplicatePrimaryKeyErr := false
		if duplicatePrimaryKeyErr {
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

	deferred, start, l := initVendorFuncs("UpdateStrain", db.logger, err, id, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, db.sql["stage"]["update"], s.Name, id)
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
	return db.deleteByUUID(ctx, id, cid, "DeleteStrain", "strain", db.logger)
}
