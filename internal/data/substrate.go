package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) SelectAllSubstrates(ctx context.Context, cid types.CID) ([]types.Substrate, error) {
	var err error

	deferred, start, l := initVendorFuncs("SelectAllSubstrates", db.logger, err, "nil", cid)
	defer deferred(start, err, l)

	var rows *sql.Rows

	result := make([]types.Substrate, 0, 100)

	rows, err = db.query.QueryContext(ctx, db.sql["substrate"]["select-all"])
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		row := types.Substrate{}
		rows.Scan(
			&row.UUID,
			&row.Name,
			&row.Type,
			&row.Vendor.UUID,
			&row.Vendor.Name)
		result = append(result, row)
	}

	return result, err
}

func (db *Conn) SelectSubstrate(ctx context.Context, id types.UUID, cid types.CID) (types.Substrate, error) {
	var err error

	deferred, start, l := initVendorFuncs("SelectSubstrate", db.logger, err, id, cid)
	defer deferred(start, err, l)

	result := types.Substrate{UUID: id}
	err = db.
		QueryRowContext(ctx, db.sql["substrate"]["select"], id).
		Scan(
			&result.Name,
			&result.Type,
			&result.Vendor.UUID,
			&result.Vendor.Name)

	return result, err
}

func (db *Conn) InsertSubstrate(ctx context.Context, s types.Substrate, cid types.CID) (types.Substrate, error) {
	var err error

	s.UUID = types.UUID(db.generateUUID().String())

	deferred, start, l := initVendorFuncs("InsertSubstrate", db.logger, err, s.UUID, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, db.sql["substrate"]["insert"], s.UUID, s.Name, s.Type, s.Vendor.UUID)
	if err != nil {
		if isUniqueViolation(err) {
			return db.InsertSubstrate(ctx, s, cid) // FIXME: infinite loop?
		}
		return s, err
	} else if rows, err := result.RowsAffected(); err != nil {
		return s, err
	} else if rows != 1 { // most likely cause is a bad vendor.uuid
		return s, fmt.Errorf("substrate was not added")
	}

	return s, err
}

func (db *Conn) UpdateSubstrate(ctx context.Context, id types.UUID, s types.Substrate, cid types.CID) error {
	var err error

	deferred, start, l := initVendorFuncs("UpdateSubstrate", db.logger, err, id, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, db.sql["substrate"]["update"], s.Name, s.Type, s.Vendor.UUID, id)
	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return fmt.Errorf("substrate was not updated: '%s'", id)
	}
	return nil
}

func (db *Conn) DeleteSubstrate(ctx context.Context, id types.UUID, cid types.CID) error {
	return db.deleteByUUID(ctx, id, cid, "DeleteSubstrate", "substrate", db.logger)
}
