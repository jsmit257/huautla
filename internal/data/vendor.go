package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) SelectAllVendors(ctx context.Context, cid types.CID) ([]types.Vendor, error) {
	var err error

	deferred, start, l := initVendorFuncs("SelectAllVendors", db.logger, "nil", cid)
	defer deferred(start, err, l)

	var rows *sql.Rows

	result := make([]types.Vendor, 0, 100)

	rows, err = db.query.QueryContext(ctx, db.sql["vendor"]["select-all"])
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		row := types.Vendor{}
		rows.Scan(&row.UUID, &row.Name)
		result = append(result, row)
	}

	return result, err
}

func (db *Conn) SelectVendor(ctx context.Context, id types.UUID, cid types.CID) (types.Vendor, error) {
	var err error

	deferred, start, l := initVendorFuncs("SelectVendor", db.logger, id, cid)
	defer deferred(start, err, l)

	result := types.Vendor{UUID: id}
	err = db.
		QueryRowContext(ctx, db.sql["vendor"]["select"], id).
		Scan(&result.Name)

	return result, err
}

func (db *Conn) InsertVendor(ctx context.Context, v types.Vendor, cid types.CID) (types.Vendor, error) {
	var err error

	v.UUID = types.UUID(db.generateUUID().String())

	deferred, start, l := initVendorFuncs("InsertVendor", db.logger, v.UUID, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, db.sql["vendor"]["insert"], v.UUID, v.Name)
	if err != nil {
		if isUniqueViolation(err) {
			return db.InsertVendor(ctx, v, cid) // FIXME: infinite loop?
		}
		return v, err
	} else if rows, err := result.RowsAffected(); err != nil {
		return v, err
	} else if rows != 1 {
		return v, fmt.Errorf("vendor was not added")
	}

	return v, err
}

func (db *Conn) UpdateVendor(ctx context.Context, id types.UUID, v types.Vendor, cid types.CID) error {
	var err error

	deferred, start, l := initVendorFuncs("UpdateVendor", db.logger, id, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, db.sql["vendor"]["update"], v.Name, id)
	if err != nil {
		if isUniqueViolation(err) {
			return db.UpdateVendor(ctx, id, v, cid) // FIXME: infinite loop?
		}
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return fmt.Errorf("vendor was not updated: '%s'", id)
	}
	return nil
}

func (db *Conn) DeleteVendor(ctx context.Context, id types.UUID, cid types.CID) error {
	return db.deleteByUUID(ctx, id, cid, "DeleteVendor", "vendor", db.logger)
}
