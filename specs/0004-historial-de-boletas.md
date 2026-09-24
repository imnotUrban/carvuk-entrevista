# 0004 — Historial de boletas

- **Estado:** implementado
- **Fecha:** 2026-09-24

## Contexto
El usuario necesita consultar las boletas ya generadas (0003) y ver el desglose de cada una. Depende de las tablas `compras`, `compra_items` y `boletas` definidas en 0003.

## Caso de uso
- **Actor:** usuario.
- **Flujo principal:** entra a `/boletas`, ve la lista (más recientes primero) con número, fecha, usuario (nombre y correo), total y cantidad de productos; abre una y ve el detalle: datos del usuario, productos (nombre, precio unitario, cantidad, subtotal) y valores **neto, impuesto y total**.
- **Alternativos:** no hay boletas → mensaje "Aún no hay boletas"; id inexistente → pantalla "Boleta no encontrada"; paginar la lista.

## Requisitos
- Las boletas son de solo lectura (sin update ni delete).
- El detalle muestra los datos **guardados** en la boleta (snapshot de nombre y precio de 0003), no los del catálogo actual.
- Tanto la lista como el detalle muestran la información del usuario dueño de la compra (`id`, `nombre`, `correo`, tomada de `compras.user_id → usuarios`), aunque hoy solo exista el usuario semilla. Es el dato actual del usuario, no un snapshot.
- Total = `valor_bruto`. Debe cumplirse: `valor_neto + impuesto = valor_bruto`, y la suma de subtotales = `valor_bruto`.
- Listado ordenado por `created_at desc` por defecto; `limit` máx 100.
- Errores: `400` parámetros inválidos, `404` boleta no encontrada.

## Diseño
- **Backend:** `repository/boleta` (`list.go`, `find.go` con `Preload` de compra, usuario e ítems; `list.go` con `Preload` de usuario), `service/boleta` (`list.go`, `find.go`), `delivery/boleta` (`list.go`, `find.go`), rutas nuevas; tests en `service/boleta/boleta_test.go`; postman.
- **Frontend:** `app/boletas/page.tsx` (tabla + paginación), `app/boletas/[id]/page.tsx` (detalle), link "Boletas" en `nav-bar.tsx`, formato CLP con `lib/format.ts`.
- **API:**

| Método | Ruta            | Descripción                                   |
|--------|-----------------|-----------------------------------------------|
| GET    | `/boletas`      | Listar (`page`, `limit`, `order` = `asc|desc`) |
| GET    | `/boletas/:id`  | Detalle con ítems                              |

```json
// GET /boletas → 200
{ "data": [ {
    "id": 1, "created_at": "...", "valor_bruto": 4200, "cantidad_items": 3,
    "usuario": { "id": 1, "nombre": "Cliente Demo", "correo": "cliente.demo@example.com" }
  } ],
  "total": 1, "page": 1, "limit": 20 }

// GET /boletas/1 → 200  (misma forma que la respuesta de POST en 0003, incluido "usuario")
```

## Criterios de aceptación
- [x] Listar devuelve las boletas creadas, la más reciente primero.
- [x] Cada elemento de la lista incluye `usuario` con `id`, `nombre` y `correo`.
- [x] El detalle incluye `usuario` con `id`, `nombre` y `correo`.
- [x] Paginación (`page`, `limit`) correcta; `limit` > 100 se limita/rechaza como en `boilerplate`.
- [x] Detalle devuelve ítems con nombre, precio unitario, cantidad y subtotal, más neto, impuesto, bruto y porcentaje.
- [x] Detalle de la boleta del ejemplo: neto 3570, impuesto 630, total 4200.
- [x] Boleta inexistente → `404`.
- [x] Tras cambiar el precio o borrar un producto, el detalle sigue mostrando los valores originales.
- [x] Front: `/boletas` lista (con columna usuario) y enlaza al detalle (que muestra los datos del usuario); estados vacío y no encontrado se muestran.

## Fuera de alcance
- Filtros por fecha/monto/usuario, búsqueda, exportación, impresión/PDF.
- Historial por usuario autenticado.

## Decisiones
- La lista y el detalle **sí** muestran la información del usuario (nombre y correo), aunque hoy haya un único usuario semilla (ver 0003).

## Preguntas abiertas
Ninguna.
