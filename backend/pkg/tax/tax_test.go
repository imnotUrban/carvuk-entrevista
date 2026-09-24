package tax_test

import (
	"testing"

	"backend/pkg/tax"
)

func TestCalculate(t *testing.T) {
	cases := []struct{ bruto, impuesto, neto int }{
		{4200, 630, 3570},
		{0, 0, 0},
		{1001, 150, 851},
		{1004, 151, 853},
	}
	for _, c := range cases {
		got := tax.Calculate(c.bruto)
		if got.Bruto != c.bruto || got.Impuesto != c.impuesto || got.Neto != c.neto {
			t.Errorf("Calculate(%d) = %+v, want impuesto %d neto %d", c.bruto, got, c.impuesto, c.neto)
		}
	}
}

func TestCalculateInvariants(t *testing.T) {
	for b := 0; b <= 10000; b++ {
		got := tax.Calculate(b)
		if got.Neto+got.Impuesto != b || got.Impuesto < 0 || got.Neto < 0 {
			t.Fatalf("invariant broken for %d: %+v", b, got)
		}
	}
}
