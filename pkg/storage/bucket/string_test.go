package bucket

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBucket(t *testing.T) {
	tests := []struct {
		name string
		want Bucket
	}{
		{"", bucket{""}},
		{"qweryuiop", bucket{"qweryuiop"}},
		{"ßąś", bucket{"ßąś"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBucket(tt.name)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.name, got.String())
			require.Equal(t, []byte(tt.name), got.Bytes())
		})
	}
}
