package types

import (
	"time"
)

type (
	UUID string

	CID string

	SubstrateType string

	Vendor struct {
		UUID `json:"-"`
		Name string `json:"name"`
	}

	Substrate struct {
		UUID        `json:"-"`
		Name        string        `json:"name"`
		Type        SubstrateType `json:"type"`
		Vendor      `json:"vendor"`
		Ingredients []Ingredient `json:"ingredients,omitempty"`
	}

	Ingredient struct {
		UUID `json:"-"`
		Name string `json:"name"`
	}

	Strain struct {
		UUID       `json:"-"`
		Name       string `json:"name"`
		Vendor     `json:"vendor"`
		Attributes []StrainAttribute `json:"attributes"`
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
)
