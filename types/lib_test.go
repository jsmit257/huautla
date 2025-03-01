package types

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/jsmit257/huautla/internal/metrics"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func Test_GetContextCID(t *testing.T) {
	t.Parallel()
	tcs := map[string]struct {
		ctx context.Context
		cid CID
	}{
		"happy_path": {
			ctx: context.WithValue(context.TODO(),
				Cid,
				CID("happy_path")),
			cid: "happy_path",
		},
		"null_attr": {
			ctx: context.TODO(),
			cid: "context has no cid attribute: context.todoCtx{emptyCtx:context.emptyCtx{}}",
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			cid := GetContextCID(tc.ctx)
			require.Equal(t, tc.cid, cid)
		})
	}
}

func Test_GetContextLog(t *testing.T) {
	t.Parallel()
	tcs := map[string]struct {
		ctx context.Context
	}{
		"happy_path": {
			ctx: context.WithValue(context.TODO(),
				Log,
				logrus.NewEntry(logrus.New())),
		},
		"null_attr": {
			ctx: context.TODO(),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			l := GetContextLog(tc.ctx)
			require.NotNil(t, l.Error)
		})
	}
}

func Test_GetContextMetrics(t *testing.T) {
	t.Parallel()
	// t.Skip() // need to pull the rest of metrics into this project
	tcs := map[string]struct {
		ctx context.Context
	}{
		"happy_path": {
			ctx: context.WithValue(context.TODO(),
				Metrics,
				metrics.DataMetrics),
		},
		"null_attr": {
			ctx: context.TODO(),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			m := GetContextMetrics(tc.ctx)
			require.NotNil(t, m.MustCurryWith)
		})
	}
}

func Test_Validate(t *testing.T) {
	t.Parallel()

	ref := time.Now().UTC()

	tcs := map[string]struct {
		flds  []string
		facts []byte
		org   *time.Time
		err   error
	}{
		"happy_path": {
			flds: []string{"ctime", "mtime"},
			org:  &ref,
			facts: []byte(`[
				{"delta": 1, "interval": "day"},
				{"delta": 0, "interval": "day"}
			]`),
		},
		"missing_fields": {
			err: fmt.Errorf("no fields specified for update"),
		},
		"missing_origin": {
			flds: []string{"ctime", "mtime"},
			err:  fmt.Errorf("origin date must be specified"),
		},
		"invalid_interval": {
			flds: []string{"ctime", "mtime"},
			org:  &ref,
			facts: []byte(`[
				{"delta": 1, "interval": "derp"}
			]`),
			err: fmt.Errorf("invalid interval: 'derp'"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ts := &Timestamp{
				Fields: tc.flds,
				Origin: tc.org,
			}

			if tc.facts != nil {
				err := json.Unmarshal(tc.facts, &ts.Factor)
				require.Nil(t, err)
			}

			err := ts.Validate()
			require.Equal(t, tc.err, err)
		})
	}
}

func Test_UpdateString(t *testing.T) {
	t.Parallel()

	ref := time.Now().UTC()

	tcs := map[string]struct {
		flds   []string
		facts  []byte
		org    *time.Time
		result string
		err    error
	}{
		"happy_path": {
			flds: []string{"ctime", "mtime", "dtime"},
			org:  &ref,
			facts: []byte(`[
				{"delta": -3, "interval": "week"},
				{"delta": 1, "interval": "day"},
				{"delta": 0, "interval": "day"}
			]`), // this is why we don't nest anonymous structs
			result: func(s string) string {
				return fmt.Sprintf(`ctime = timestamp '%s' + interval '-3 week' + interval '1 day', mtime = timestamp '%s' + interval '-3 week' + interval '1 day', dtime = timestamp '%s' + interval '-3 week' + interval '1 day'`,
					s, s, s) // must be a better way
			}(ref.Format(time.RFC3339)),
		},
		"standard_error": {
			err: fmt.Errorf("no fields specified for update"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ts := &Timestamp{
				Fields: tc.flds,
				Origin: tc.org,
			}

			if tc.facts != nil {
				err := json.Unmarshal(tc.facts, &ts.Factor)
				require.Nil(t, err)
			}

			result, err := ts.UpdateString()
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}
