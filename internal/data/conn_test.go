package data

import (
	"time"

	"github.com/google/uuid"
)

func mockUUIDGen() uuid.UUID {
	return uuid.Must(uuid.FromBytes([]byte("0123456789abcdef")))
}

var whenwillthenbenow = time.Now() // time.Soon()
