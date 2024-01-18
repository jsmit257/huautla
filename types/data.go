package types

import (
	"context"
	"time"
)

type (
	UUID string

	CID string

	Vendor struct {
		UUID `json:"-"`
		Name string `json:"name"`
	}

	Substrate struct {
		UUID   `json:"-"`
		Name   string `json:"name"`
		Vendor `json:"vendor"`
	}

	Ingredient struct {
		UUID      `json:"-"`
		Name      string `json:"name"`
		Substrate `json:"substrate"`
	}

	Strain struct {
		UUID   `json:"-"`
		Name   string `json:"name"`
		Vendor `json:"vendor"`
	}

	StrainAttribute struct {
		UUID   `json:"-"`
		Name   string `json:"name"`
		Value  string `json:"value"`
		Strain `json:"strain"`
	}

	Stage struct {
		UUID `json:"-"`
		Name string `json:"name"`
	}

	EventType struct {
		UUID  `json:"-"`
		Name  string `json:"name"`
		Stage `json:"stage"`
	}

	Lifecycle struct {
		UUID           `json:"-"`
		GrainCost      int16     `json:"grain_cost"`
		BulkCost       int16     `json:"bulk_cost"`
		Yield          int16     `json:"Yield"`
		Count          int16     `json:"count"`
		Name           string    `json:"name"`
		Gross          int16     `json:"gross"`
		MTime          time.Time `json:"modified_date"`
		CTime          time.Time `json:"create_date"`
		Strain         `json:"strain"`
		GrainSubstrate Substrate `json:"grain_substrate"`
		BulkSubstrate  Substrate `json:"bulk_substrate"`
	}

	Event struct {
		UUID        `json:"-"`
		Temperature int8      `json:"temp"` // temp? sounds like temporary instead of temperature
		MTime       time.Time `json:"modified_date"`
		CTime       time.Time `json:"create_date"`
		Lifecycle   Lifecycle `json:"lifecycle"`
		EventType   EventType `json:"event_type"`
	}

	Eventer      interface{}
	EventTyper   interface{}
	Ingredienter interface{}
	Lifecycler   interface{}
	Stager       interface{}
	Substrater   interface{}

	// vendors aren't meant to be a comprehensive list of attributes, really just
	// a name that makes a unique constraint with Strain.Name in the strain table;
	// storing them is just for convenience, and to avoid spelling errors
	Vendorer interface {
		// this doesn't support pagination, because our list of vendors will only
		// range into the hundreds *at most*
		SelectAllVendors(ctx context.Context, cid CID) ([]Vendor, error)
		// mostly useful for filling in complex objects like strains or substrates
		SelectVendor(ctx context.Context, id UUID, cid CID) (Vendor, error)
		// assigns a generated uuid, ctime and mtime to `v` and inserts it into vendors;
		// the modified `v` struct is returned
		InsertVendor(ctx context.Context, v Vendor, cid CID) (Vendor, error)
		// assigns a generated mtime to the `v` arg, then updates every vendor attribute
		// except uuid based on the `id` arg; overwrites the uuid in `v` and returns `v`
		UpdateVendor(ctx context.Context, id UUID, v Vendor, cid CID) error
		// this should really never be used; none of the foreign keys have cascade
		// delete, so any referenced vendors will cause an error here; just don't use it
		DeleteVendor(ctx context.Context, id UUID, cid CID) error
	}

	DB interface {
		Eventer
		EventTyper
		Ingredienter
		Lifecycler
		Stager
		Substrater
		Vendorer
	}
)
