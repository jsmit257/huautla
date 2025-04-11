package data

import (
	"time"

	"github.com/google/uuid"
	pq "github.com/lib/pq"
)

func mockUUIDGen() uuid.UUID {
	return uuid.Must(uuid.FromBytes([]byte("0123456789abcdef")))
}

func pqError(code, detail, table, field, constraint string) error {
	return &pq.Error{
		Code:       pq.ErrorCode(code),
		Detail:     detail,
		Table:      table,
		Column:     field,
		Constraint: constraint,
	}
}

func pkerr() error {
	return PKeyError(pqError("23505", "unique/primary key", "table", "field", "constraint"))
}

var wwtbn = time.Now() // time.Soon()
