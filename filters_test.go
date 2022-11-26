package dsp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFakedFreq(t *testing.T) {
	r := require.New(t)

	cof, sr := fakedFreq(0.0005, 0.016666)
	r.Equal(0.5, cof)
	r.Equal(16, sr)

	cof, sr = fakedFreq(0.0005, 0.0016666)
	r.Equal(0.5, cof)
	r.Equal(1, sr)

	cof, sr = fakedFreq(0.0005, 0.00016666)
	r.Equal(500.0, cof)
	r.Equal(166, sr)
}
