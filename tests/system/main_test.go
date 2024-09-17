package test

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jsmit257/huautla"
	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"
)

var (
	db    types.DB
	epoch time.Time
)

const (
	foreignKeyViolation1to1    = `pq: insert or update on table "%s" violates foreign key constraint "%s"`
	foreignKeyViolation1toMany = `pq: update or delete on table "%s" violates foreign key constraint "%s" on table "%s"`
	uniqueKeyViolation         = `pq: duplicate key value violates unique constraint "%s"`
	checkConstraintViolation   = `pq: new row for relation "%s" violates check constraint "%s"`
)

func init() {
	var err error
	var location *time.Location
	location, _ = time.LoadLocation("Etc/UTC")
	epoch = time.Date(1970, 1, 1, 0, 0, 0, 0, location)

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

func equalErrorMessages(t *testing.T, expected, actual error) {
	if expected != nil && actual != nil {
		require.Equal(t, expected.Error(), actual.Error())
	} else if expected != nil || actual != nil {
		require.Equal(t, expected, actual)
	}
}

type testAttrs map[string]string

func (ta testAttrs) Get(name string) *types.UUID {
	result, ok := ta[name]
	if !ok {
		return nil
	}
	temp := types.UUID(result)
	return &temp
}
func (ta testAttrs) Contains(names ...string) bool {
	return true
}
func (ta testAttrs) Map(m url.Values) (err error) {
	return nil
}
func (ta testAttrs) Set(name string, value string) error {
	return nil
}
