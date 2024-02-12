package test

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jsmit257/huautla"
	"github.com/jsmit257/huautla/types"

	log "github.com/sirupsen/logrus"
)

var db types.DB
var noRows error = fmt.Errorf("sql: no rows in result set")
var epoch time.Time = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

const (
	invalidUUID types.UUID = "01234567890123456789012345678901234567890"

	foreignKeyViolation string = `pq: update or delete on table "%s" violates foreign key constraint "%s" on table "%s"`
)

func init() {
	var err error

	if db, err = huautla.New(
		&types.Config{
			PGHost: os.Getenv("POSTGRES_HOST"),
			PGUser: os.Getenv("POSTGRES_USER"),
			PGPass: os.Getenv("POSTGRES_PASSWORD"),
			PGPort: func(s string) uint {
				i, err := strconv.Atoi(s)
				if err != nil {
					panic(fmt.Errorf("pgport(%q) could not be parsed as int: %w", s, err))
				} else if i < 1 {
					panic(fmt.Errorf("pgport(%d) cannot be less than 1", i))
				}
				return uint(i)
			}(os.Getenv("POSTGRES_PORT")),
			PGSSL: os.Getenv("POSTGRES_SSLMODE"),
		},
		log.WithField("test-suite", "system")); err != nil {

		panic(err)
	}
}
