package util_test

import (
	"testing"

	"github.com/tgragnato/magnetico/util"
)

func TestRoundToDecimal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		input         float64
		decimalPlaces int
		want          float64
	}{
		{"round to 1 decimal places", 1.2345, 1, 1.2},
		{"round to 2 decimal places", 1.2345, 2, 1.23},
		{"round to 4 decimal places", 1.2345, 4, 1.2345},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := util.RoundToDecimal(test.input, test.decimalPlaces); got != test.want {
				t.Errorf("RoundToDecimal() = %v, want %v", got, test.want)
			}
		})
	}
}
