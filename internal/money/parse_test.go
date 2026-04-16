package money

import "testing"

func TestParseEURToCents(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in      string
		want    int64
		wantErr bool
	}{
		{"0", 0, false},
		{"0.00", 0, false},
		{"12.34", 1234, false},
		{"12.3", 1230, false},
		{"1,234.56", 123456, false},
		{"€12.50", 1250, false},
		{"  €3.00  ", 300, false},
		{"  42.00  ", 4200, false},
		{"-10.00", -1000, false},
		{"", 0, true},
		{"12.345", 0, true},
		{"1.2.3", 0, true},
		{"not", 0, true},
	}
	for _, tc := range tests {
		got, err := ParseEURToCents(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ParseEURToCents(%q) want error, got %d", tc.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseEURToCents(%q): %v", tc.in, err)
			continue
		}
		if got != tc.want {
			t.Errorf("ParseEURToCents(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestAbsCents(t *testing.T) {
	t.Parallel()
	if AbsCents(-5) != 5 || AbsCents(5) != 5 || AbsCents(0) != 0 {
		t.Fatal()
	}
}
