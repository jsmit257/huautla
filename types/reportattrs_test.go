package types

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewReportAttrs(t *testing.T) {
	t.Parallel()

	tcs := map[string]struct {
		vals url.Values
		err  error
	}{
		"happy_path": {
			vals: url.Values{"bulk-id": []string{"bulk-id"}},
		},
		"empty_path": {
			// still happy, just empty
		},
		"empty_value": {
			vals: url.Values{"bulk-id": []string{""}},
			err:  fmt.Errorf("failed to find param values in the following fields: [bulk-id]"),
		},
		"bogus_name": {
			vals: url.Values{"bogus-id": []string{"bogus"}},
			err:  fmt.Errorf("failed to find param values in the following fields: [bogus-id]"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ra, err := NewReportAttrs(tc.vals)
			require.NotNil(t, ra)
			require.Equal(t, tc.err, err)
		})
	}
}

func Test_Set(t *testing.T) {
	t.Parallel()

	tcs := map[string]struct {
		name, value string
		err         error
	}{
		"happy_path": {
			name:  "bulk-id",
			value: "bulk-id",
		},
		"empty_value": {
			name:  "bulk-id",
			value: "",
			err:   fmt.Errorf("empty value for key: bulk-id"),
		},
		"bogus_name": {
			name:  "bogus-id",
			value: "bogus-id",
			err:   fmt.Errorf("unknown parameter: bogus-id"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r, _ := NewReportAttrs(url.Values{})
			err := r.Set(tc.name, tc.value)
			require.Equal(t, tc.err, err)
		})
	}
}

func Test_Get(t *testing.T) {
	t.Parallel()

	bulkID := UUID("bulk-id")

	tcs := map[string]struct {
		seed   url.Values
		name   string
		result *UUID
	}{
		"happy_path": {
			seed:   url.Values{"bulk-id": []string{"bulk-id"}},
			name:   "bulk-id",
			result: &bulkID,
		},
		"empty_value": {
			// seed:   url.Values{"bulk-id": []string{"bulk-id"}},
			name:   "bulk-id",
			result: nil,
		},
		"bogus_name": {
			name:   "bogus-id",
			result: nil,
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r, _ := NewReportAttrs(tc.seed)
			result := r.Get(tc.name)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_Contains(t *testing.T) {
	t.Parallel()

	tcs := map[string]struct {
		seed   url.Values
		names  []string
		result bool
	}{
		"happy_path": {
			seed:   url.Values{"bulk-id": []string{"bulk-id"}},
			names:  []string{"bulk-id"},
			result: true,
		},
		"empty_value": {
			// seed:   url.Values{"bulk-id": []string{"bulk-id"}},
			names: []string{"bulk-id"},
		},
		"bogus_name": {
			names: []string{"bogus-id"},
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r, _ := NewReportAttrs(tc.seed)
			result := r.Contains(tc.names...)
			require.Equal(t, tc.result, result)
		})
	}
}
