# 0003 — Generar boleta a partir del carrito

- **Estado:** implementado
- **Fecha:** 2026-09-24

## Contexto
Al confirmar el carrito (0002) se registra una **compra** y se emite una **boleta** con los valores calculados. El esquema del diagrama es: `usuario 1—n compras`, `compras n—n productos`, `compras 1—n boletas` (`boletas.id_compra`).

El diagrama guarda `compras.productos_id[]`, que **no permite representar cantidades** (el ejemplo tiene 2 Leches) ni conservar el precio con que se vendió. Por eso este spec usa una tabla de detalle `compra_items` en lugar del arreglo (decisión confirmada, ver Decisiones).

## Caso de uso
- **Actor:** usuario que compra.
- **Flujo principal:** desde el carrito pulsa "Generar boleta" → el front envía las líneas `{producto_id, cantidad}` → el backend valida, calcula, guarda compra + ítems + boleta en una transacción → responde con la boleta → el front limpia el carrito y muestra el detalle de la boleta (`/boletas/:id`, spec 0004).
- **Alternativos:** carrito vacío → error; producto inexistente/eliminado → error y el carrito no se limpia; falla de red → el carrito se conserva.

## Requisitos
### Cálculo (regla de 0005)
Con enteros CLP:
- `valor_bruto` = Σ (`precio_unitario` × `cantidad`) — total que paga el cliente.
- `impuesto` = `round(valor_bruto × 15 / 100)`, redondeo al entero más cercano (half up).
- `valor_neto` = `valor_bruto − impuesto`.
- `porcentaje_impuesto` = 15 (se guarda en la boleta).

Ejemplo: 2 Leches ($1.100) + 1 Aceite ($2.000) → bruto 4.200, impuesto 630, neto 3.570.

### Reglas
- Los precios los toma **el backend** desde `productos`; el request nunca trae precios, impuesto ni totales.
- `items` no vacío (máx. 100 líneas); `cantidad` entero `1..99`; `producto_id` sin repetir en el request (si se repite → `400`).
- Todos los productos deben existir y no estar borrados; si alguno no → `404` indicando el id.
- **Snapshot:** cada ítem guarda `nombre` y `precio_unitario` vigentes al comprar, para que la boleta histórica no cambie si luego se edita/borra el producto.
- **Stock (ver 0001):** cada producto debe tener `stock >= cantidad` pedida. Al generar la boleta se descuenta `cantidad` del stock de cada producto, en la **misma transacción** que crea compra, ítems y boleta.
  - El descuento es atómico y evita ventas simultáneas por encima del stock: `UPDATE productos SET stock = stock - :cant WHERE id = :id AND stock >= :cant`; si no afecta ninguna fila, se revierte todo y se responde `409`.
  - El error se modela como `apperr.ErrInsufficientStock` y `httpkit` lo mapea a `409`; el mensaje indica producto y stock disponible.
- Compra + ítems + boleta + descuento de stock se crean en una **única transacción**.
- `user_id`: no hay autenticación en el proyecto. Se crea la tabla `usuarios` (`id`, `nombre`, `correo`) con un usuario semilla inicial "Cliente Demo" (`cliente.demo@example.com`), insertado al arrancar si la tabla está vacía. El request no trae `user_id`: toda compra se asocia a ese usuario semilla. Sin selector ni CRUD de usuarios.
- Errores: `400` validación, `404` producto no encontrado, `409` stock insuficiente.

## Diseño
- **Migración** `0003_create_compras_boletas.up/down.sql`:
  - `usuarios(id, nombre, correo, created_at)`
  - `compras(id, user_id → usuarios, created_at)`
  - `compra_items(id, compra_id → compras, producto_id → productos, nombre, precio_unitario, cantidad)`
  - `boletas(id, id_compra → compras, valor_neto, valor_bruto, impuesto, porcentaje_impuesto, created_at)` — se agrega `impuesto` (el diagrama solo tiene neto, bruto y porcentaje) para poder mostrarlo tal cual se calculó.
- **Backend** (mismo patrón por capas): `pkg/apperr` (nuevo `ErrInsufficientStock`) y `pkg/httpkit` (→ `409`), `entity/{usuario,compra,boleta}`, `repository/boleta` (crea compra+ítems+boleta y descuenta stock con `db.Transaction`), `service/boleta/create.go` (validación + cálculo usando la función de 0005), `delivery/boleta/create.go`, ruta en `delivery/router.go`, seed del usuario demo en `main.go`, `postman_collection.json`.
- **Frontend:** botón "Generar boleta" en `cart-panel.tsx`, `lib/api.ts` (`createBoleta`), redirección a `/boletas/:id` y `cart.clear()` solo si la respuesta fue `201`.
- **API:**

```json
// POST /api/v1/boletas
{ "items": [ { "producto_id": 1, "cantidad": 2 }, { "producto_id": 3, "cantidad": 1 } ] }

// 201
{
  "id": 1, "compra_id": 1,
  "valor_bruto": 4200, "impuesto": 630, "valor_neto": 3570,
  "porcentaje_impuesto": 15,
  "usuario": { "id": 1, "nombre": "Cliente Demo", "correo": "cliente.demo@example.com" },
  "items": [
    { "producto_id": 1, "nombre": "Leche",  "precio_unitario": 1100, "cantidad": 2, "subtotal": 2200 },
    { "producto_id": 3, "nombre": "Aceite", "precio_unitario": 2000, "cantidad": 1, "subtotal": 2000 }
  ],
  "created_at": "..."
}
```

## Criterios de aceptación
- [x] Ejemplo del enunciado produce bruto 4200, impuesto 630, neto 3570.
- [x] Se crean compra, ítems y boleta asociados (`boleta.id_compra` = compra creada).
- [x] Los precios usados son los de la DB aunque el request intente enviar otros.
- [x] `items` vacío → `400`; `cantidad` 0, negativa o > 99 → `400`; `producto_id` repetido → `400`.
- [x] Producto inexistente o soft-deleteado → `404` y no se crea nada (transacción revertida).
- [x] Al generar la boleta, el stock de cada producto disminuye en la cantidad comprada.
- [x] Cantidad mayor al stock disponible → `409`; no se crea nada y el stock de ningún producto cambia (transacción revertida, incluso si otro producto de la boleta sí alcanzaba).
- [x] Comprar exactamente el stock disponible es válido y deja el stock en 0; luego otra compra del producto → `409`.
- [x] Tras editar el precio o borrar un producto, la boleta ya emitida conserva nombre y precio originales.
- [x] Redondeo: bruto 1.001 → impuesto 150, neto 851 (150,15 → 150).
- [x] Front: tras `201` se limpia el carrito y se navega a la boleta; ante error se conserva el carrito y se muestra el mensaje. (La página `/boletas/:id` pertenece a 0004 y aún no existe: la navegación apunta a ella.)

## Fuera de alcance
- Pagos, medios de pago, folio/numeración SII, PDF/impresión de boleta.
- Anular o editar boletas (son inmutables), por lo tanto tampoco se devuelve stock.
- Reservar stock mientras el producto está en el carrito (solo se descuenta al generar la boleta).
- Autenticación real de usuarios.

## Decisiones
1. **`compras.productos_id[]` → `compra_items`:** se reemplaza el arreglo por la tabla de detalle `compra_items` (necesaria para cantidades y snapshot de precio).
2. **Usuario:** un usuario semilla inicial ("Cliente Demo"); no hay selector ni CRUD de usuarios.
3. **Redondeo del impuesto:** half up al entero (regla en 0005).

## Preguntas abiertas
Ninguna.
