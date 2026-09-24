package tax

// TaxRatePercent es la tasa de impuesto fija (%) aplicada a toda boleta.
const TaxRatePercent = 15

type Totals struct {
	Bruto    int
	Impuesto int
	Neto     int
}

// Calculate es la única implementación del cálculo de impuesto.
// Impuesto = round_half_up(bruto * 15 / 100); Neto = Bruto - Impuesto.
func Calculate(bruto int) Totals {
	impuesto := (bruto*TaxRatePercent + 50) / 100
	return Totals{Bruto: bruto, Impuesto: impuesto, Neto: bruto - impuesto}
}
