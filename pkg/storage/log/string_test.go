package log

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEntry(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		creator func() Entry
		want    Entry
	}{
		{"empty from string", "", func() Entry { return FromString("") }, entry("")},
		{"empty from bytes", "", func() Entry { return FromBytes([]byte{}) }, entry("")},
		{"empty from JSON", "", func() Entry {
			var e entry
			if err := json.Unmarshal([]byte(`""`), &e); err != nil {
				panic(err)
			}
			return e
		}, entry("")},

		{"qweryuiop from string", "qweryuiop", func() Entry { return FromString("qweryuiop") }, entry("qweryuiop")},
		{"qweryuiop from bytes", "qweryuiop", func() Entry { return FromBytes([]byte(`qweryuiop`)) }, entry("qweryuiop")},
		{"qweryuiop from JSON", "qweryuiop", func() Entry {
			var e entry
			if err := json.Unmarshal([]byte(`"qweryuiop"`), &e); err != nil {
				panic(err)
			}
			return e
		}, entry("qweryuiop")},

		{"ßąś from string", "ßąś", func() Entry { return FromString("ßąś") }, entry("ßąś")},
		{"ßąś from bytes", "ßąś", func() Entry { return FromBytes([]byte(`ßąś`)) }, entry("ßąś")},
		{"ßąś from JSON", "ßąś", func() Entry {
			var e entry
			if err := json.Unmarshal([]byte(`"ßąś"`), &e); err != nil {
				panic(err)
			}
			return e
		}, entry("ßąś")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.creator()
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.val, got.String())
			require.Equal(t, []byte(tt.val), got.Bytes())
		})
	}
}
