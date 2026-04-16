package htmlview

import (
	"math"
	"testing"
)

func TestFormatPercentSigned(t *testing.T) {
	t.Parallel()
	if got := FormatPercentSigned(12.3); got != "+12.3%" {
		t.Fatalf("got %q", got)
	}
	if got := FormatPercentSigned(-4.5); got != "-4.5%" {
		t.Fatalf("got %q", got)
	}
	if got := FormatPercentSigned(math.NaN()); got != "—" {
		t.Fatalf("NaN: got %q", got)
	}
	if got := FormatPercentSigned(math.Inf(1)); got != "—" {
		t.Fatalf("+Inf: got %q", got)
	}
}

func TestIsNegFloat(t *testing.T) {
	t.Parallel()
	if !IsNegFloat(-1.0) {
		t.Fatal("want negative")
	}
	if IsNegFloat(0.0) || IsNegFloat(1.0) {
		t.Fatal("non-negative")
	}
	if IsNegFloat(math.NaN()) || IsNegFloat(math.Inf(-1)) {
		t.Fatal("NaN/Inf not negative in our sense")
	}
}
