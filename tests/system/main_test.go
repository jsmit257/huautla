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

const invalidUUID types.UUID = "01234567890123456789012345678901234567890"

func init() {
	var err error

	db, err = huautla.New(
		&types.Config{
			PGHost: os.Getenv("pghost"),
			PGUser: os.Getenv("pguser"),
			PGPass: os.Getenv("pgpass"),
			PGPort: func(s string) uint {
				i, err := strconv.Atoi(s)
				if err != nil {
					panic(fmt.Errorf("pgport(%q) could not be parsed as int: %w", s, err))
				} else if i < 1 {
					panic(fmt.Errorf("pgport(%d) cannot be less than 1", i))
				}
				return uint(i)
			}(os.Getenv("pgport")),
		},
		log.WithField("test-suite", "system"))

	panic(err)
}
