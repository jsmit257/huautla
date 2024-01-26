package data

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_readSQL(t *testing.T) {
	sql, err := readSQL("./pgsql.yaml")
	require.Nil(t, err)
	require.NotEmpty(t, sql)
	sql, err = readSQL("./bogus.yaml")
	require.NotNil(t, err)
	require.Empty(t, sql)
}

func mockUUIDGen() uuid.UUID {
	return uuid.Must(uuid.FromBytes([]byte("0123456789abcdef")))
}
