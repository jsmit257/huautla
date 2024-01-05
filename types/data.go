package types

import (
	"time"
)

type (
	EventType struct {
		ID string `json:"id,omitempty"`
	}

	Event struct {
		ID    string `json:"id,omitempty"`
		MTime time.Time
		CTime time.Time
	}

	InnoculationEvent struct {
		Event
		Complete int8
	}

	BinEvent struct {
		Event
	}

	ColonizationEvent struct {
		Event
	}

	HarvestEvent struct {
		Event
	}

	SunsetEvent struct {
		Event
	}
)
