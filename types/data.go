package types

import (
	"time"
)

type (
	UUID string

	EventType struct {
		UUID `json:"-"`
	}

	Lifecycle struct {
		UUID `json:"-"`
		Name string `json:"lifecycle_stage"`
	}

	Event struct {
		UUID        `json:"-"`
		Temperature int8      `json:"temp"` // temp? sounds like temporary instead of temperature
		MTime       time.Time `json:"modified_date"`
		CTime       time.Time `json:"create_date"`
		EventType   EventType `json:"-"`
		Lifecycle   Lifecycle `json:"-"`
	}

	Lifecycler interface{}

	DB interface {
		Lifecycler
	}
)
