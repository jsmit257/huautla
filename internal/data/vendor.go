package data

import (
	"context"
	"fmt"
	"time"

	"github.com/jsmit257/huautla/types"

	log "github.com/sirupsen/logrus"
)

const (
	selectVendor = "select-vendor"
	updateVendor = "update-vendor"
	insertVendor = "insert-vendor"
	deleteVendor = "delete-vendor"
)

func (db *Conn) selectAllVendors(ctx context.Context, cid types.CID) ([]types.Vendor, error) {
	start := time.Now()
	l := db.logger.WithFields(log.Fields{
		"method": "selectAllVendors",
		"cid":    cid,
	})

	var err error

	defer func() {
		duration := time.Since(start)

		l.
			WithField("duration", duration).
			WithError(err).
			Infof("finished work")

		// TODO: metrics
	}()

	result := make([]types.Vendor, 0, 100)

	if rows, err := db.query.ExecContext(ctx, db.sql["select-vendor"]); err != nil {
		return nil, err
	} else if rows == nil {
		return nil, fmt.Errorf("no result returned from selectAllVendor")
	}

	return result, err
}

func (db *Conn) selectVendor(ctx context.Context, id types.UUID, cid types.CID) (types.Vendor, error) {
	start := time.Now()
	l := db.logger.WithFields(log.Fields{
		"method": "selectVendor",
		"cid":    cid,
		"id":     id,
	})

	var err error

	l.Info("starting work")

	defer func() {
		duration := time.Since(start)

		l.
			WithField("duration", duration).
			WithError(err).
			Infof("finished work")

		// TODO: metrics
	}()

	result := types.Vendor{UUID: id}
	err = db.
		QueryRowContext(ctx, db.sql["select-vendor"], id).
		Scan(&result.Name)

	return result, err
}
