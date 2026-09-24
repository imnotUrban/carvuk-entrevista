# 0005 — Impuesto fijo automático (15%)

- **Estado:** implementado
- **Fecha:** 2026-09-24

## Contexto
Regla de negocio transversal: el impuesto es **siempre 15%** del valor bruto, sin importar los productos, y debe calcularse igual cada vez que se genera un documento. Este spec fija cómo se garantiza esa consistencia; la usa 0003 al emitir boletas.

## Caso de uso
- **Actor:** negocio (regla, no una pantalla).
- **Flujo:** al generar cualquier boleta, el sistema aplica 15% automáticamente; ni el cliente (front/API) ni el catálogo pueden cambiarlo.
- **Alternativo:** un request que intenta enviar `porcentaje_impuesto`, `impuesto`, `valor_neto` o `valor_bruto` → esos campos se ignoran.

## Requisitos
- Constante única `TaxRatePercent = 15` (en `pkg/tax`), sin configuración por producto, usuario ni variable de entorno.
- Función pura única `tax.Calculate(bruto int) Totals{Bruto, Impuesto, Neto}`:
  - `Impuesto = (bruto × 15 + 50) / 100` (división entera → redondeo half up)
  - `Neto = Bruto − Impuesto`
- Es la **única** implementación del cálculo: `service/boleta` la llama y nadie más recalcula (ni repository, ni handler, ni front para persistir).
- La boleta guarda `porcentaje_impuesto = 15` como dato histórico; si la tasa cambiara en el futuro, las boletas antiguas conservan la suya.
- El front puede mostrar un total estimado (0002), pero el backend es la fuente de verdad.
- Invariantes garantizados: `Neto + Impuesto = Bruto`; `Impuesto >= 0`; `Neto >= 0`.

## Diseño
- **Backend:** `pkg/tax/tax.go` + `pkg/tax/tax_test.go`; uso en `service/boleta/create.go`. El DTO de request de `POST /boletas` no declara esos campos. La consistencia se garantiza solo en `service` (sin `CHECK` en la DB).
- **Frontend:** mostrar "Impuesto (15%)" con el valor devuelto por la API; no calcular ni enviar impuestos.

## Criterios de aceptación
- [x] `Calculate(4200)` → impuesto 630, neto 3570.
- [x] `Calculate(0)` → 0/0/0.
- [x] Redondeo: `Calculate(1001)` → impuesto 150, neto 851; `Calculate(1004)` → impuesto 151, neto 853 (150,6 → 151).
- [x] Para un rango de brutos (ej. 0..10.000) siempre `Neto + Impuesto = Bruto`.
- [x] Dos boletas con el mismo bruto (aunque con productos distintos) tienen el mismo impuesto y porcentaje 15.
- [x] Enviar `porcentaje_impuesto: 19` (u otros campos calculados) en `POST /boletas` no altera el resultado: se guarda 15.
- [x] `porcentaje_impuesto` de toda boleta creada es 15.

## Fuera de alcance
- Tasas distintas por producto/categoría/país, exenciones, configuración dinámica de la tasa.
- Desglose de impuesto por línea (solo se calcula sobre el bruto total).

## Decisiones
- Sin `CHECK` en la DB por ahora: la consistencia (`neto + impuesto = bruto`, tasa 15) se garantiza solo con el cálculo en `service`. Se puede agregar más adelante con una migración.
- El impuesto se calcula sobre el total y no por línea; la suma de impuestos por línea podría diferir en ±1 CLP, y se acepta porque no se muestra por línea.

## Preguntas abiertas
Ninguna.
