package huautla

import (
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	tcs := map[string]struct {
		cfg types.Config
		err error
	}{
		// unfortunately, no happy path for this
		"missing_hostname": {
			cfg: types.Config{
				PGUser: "postgres",
				PGPass: "root",
				PGPort: 5432,
			},
			err: fmt.Errorf("postgres connection needs hostname attribute"),
		},
		"missing_username": {
			cfg: types.Config{
				PGHost: "huautla",
				PGPass: "root",
				PGPort: 5432,
			},
			err: fmt.Errorf("postgres connection needs username attribute"),
		},
		"missing_password": {
			cfg: types.Config{
				PGHost: "huautla",
				PGUser: "postgres",
				PGPort: 5432,
			},
			err: fmt.Errorf("postgres connection needs password attribute"),
		},
		"missing_port": {
			cfg: types.Config{
				PGHost: "huautla",
				PGUser: "postgres",
				PGPass: "root",
			},
			err: fmt.Errorf("postgres connection needs port attribute"),
		},
		"missing_ssl": {
			cfg: types.Config{
				PGHost: "huautla",
				PGUser: "postgres",
				PGPass: "root",
				PGPort: 5432,
			},
		},
	}

	for n, v := range tcs {
		n, v := n, v
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			_, err := New(&v.cfg, nil)
			require.Equal(t, v.err, err)
		})
	}
}
