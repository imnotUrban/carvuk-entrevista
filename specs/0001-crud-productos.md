# 0001 — CRUD de productos

- **Estado:** implementado
- **Fecha:** 2026-09-24

## Contexto
El catálogo de productos es la base de las compras y las boletas (specs 0002–0005). Se necesita administrarlo (crear, listar, ver, editar, borrar) y que parta con productos de ejemplo. Reemplaza a la entidad genérica `boilerplate` como caso real; `boilerplate` se mantiene como plantilla.

## Caso de uso
- **Actor:** usuario/administrador del catálogo.
- **Flujo principal:** entra a `/productos`, ve la lista, crea un producto (nombre + precio), lo edita o lo elimina.
- **Alternativos:** filtra por nombre; intenta guardar datos inválidos y ve el error; elimina un producto que ya está en boletas históricas (las boletas no se alteran, ver 0003).

## Requisitos
- Entidad `producto` (tabla `productos`):

| Campo        | Tipo         | Notas                                        |
|--------------|--------------|----------------------------------------------|
| `id`         | uint         | PK                                           |
| `nombre`     | varchar(160) | requerido, se hace trim                      |
| `precio`     | int          | CLP entero, `> 0`. Precio final (con impuesto incluido) |
| `stock`      | int          | unidades disponibles, `>= 0`, default 0      |
| `created_at`, `updated_at`, `deleted_at` | timestamps | soft delete (`gorm.DeletedAt`) |

- Validaciones (service → `apperr.ErrValidation` → `400`): `nombre` no vacío; `precio` entero mayor que 0; `stock` entero mayor o igual a 0.
- **Stock:** se define al crear (opcional, default 0) y se ajusta con `PUT /productos/:id` enviando el nuevo valor absoluto de `stock`. No puede quedar negativo. El descuento de stock al generar una boleta **no** está en este spec (ver Preguntas abiertas).
- Base de datos: columna `stock INTEGER NOT NULL DEFAULT 0` con `CHECK (stock >= 0)` en la migración `0002_create_productos`; en la entity `Stock int` con `gorm:"not null;default:0"`.
- `nombre` **no** es único (no se pide en el esquema).
- `DELETE` es soft delete: no borra la fila, así las boletas históricas siguen pudiendo referenciarlo.
- `PUT` es parcial (solo actualiza los campos enviados), igual que `boilerplate`.
- Errores: `400` validación, `404` no encontrado.
- **Seed:** al arrancar, si la tabla `productos` está vacía, insertar: Leche ($1.100, stock 50), Pan ($1.800, stock 100), Aceite ($2.000, stock 30). (Leche y Aceite vienen del ejemplo del enunciado; el precio del pan y los stocks son inventados.)

## Diseño
- **Backend** (mismo patrón que `boilerplate`, una carpeta por capa):
  - `entity/producto/producto.go`
  - `repository/producto/` (`repository.go`, `create/find/list/update/delete.go`)
  - `service/producto/` (`service.go`, `create/find/list/update/delete.go`, `producto_test.go`)
  - `delivery/producto/` (`handler.go`, `factory.go`, `create/find/list/update/delete.go`) y registro en `delivery/router.go`
  - `migrations/0002_create_productos.up.sql / .down.sql`
  - Seed en `service/producto` (`SeedIfEmpty`) llamado desde `main.go` tras `AutoMigrate`.
  - Agregar la colección al `postman_collection.json`.
- **Frontend:** `app/productos/page.tsx` (tabla + filtro por nombre + paginación), `components/producto-form-dialog.tsx`, tipos y funciones en `lib/types.ts` / `lib/api.ts`, link en `nav-bar.tsx`.
- **API** (`/api/v1`):

| Método | Ruta             | Descripción                      |
|--------|------------------|----------------------------------|
| POST   | `/productos`     | Crear                            |
| GET    | `/productos`     | Listar (`nombre`, `page`, `limit` máx 100, `sort_by` = `id,nombre,precio,stock,created_at`, `order`) |
| GET    | `/productos/:id` | Obtener uno                      |
| PUT    | `/productos/:id` | Actualizar (parcial)             |
| DELETE | `/productos/:id` | Soft delete                      |

```json
// POST /productos
{ "nombre": "Leche", "precio": 1100, "stock": 50 }
// 201
{ "id": 1, "nombre": "Leche", "precio": 1100, "stock": 50, "created_at": "...", "updated_at": "..." }
```

## Criterios de aceptación
- [x] Crear producto válido devuelve `201` con `id`.
- [x] Nombre vacío (o solo espacios) → `400`.
- [x] Precio `0`, negativo o no entero → `400`.
- [x] Stock negativo o no entero → `400`; si se omite al crear, queda en 0.
- [x] Update de `stock` con un valor válido lo persiste; el resto de campos no cambia.
- [x] La migración crea la columna `stock` con default 0 y `CHECK (stock >= 0)`.
- [x] Obtener por id existente → `200`; inexistente → `404`.
- [x] Listar filtra por nombre (contiene, case-insensitive), pagina y ordena por whitelist; `sort_by` inválido → `400`.
- [x] Update parcial modifica solo los campos enviados; con datos inválidos → `400`; inexistente → `404`.
- [x] Delete hace soft delete: deja de listarse/obtenerse pero la fila sigue en la DB; inexistente → `404`.
- [x] Con tabla vacía, el seed crea Leche, Pan y Aceite con su stock; con datos existentes no duplica.
- [x] La página `/productos` muestra la columna stock, permite crear, editar (incluido el stock) y eliminar, y muestra errores del backend.

## Fuera de alcance
- Historial/movimientos de stock, alertas de stock bajo, stock por bodega.
- Categorías, imágenes, descuentos.
- Precios con decimales / otras monedas (solo CLP entero).
- Autenticación/roles para administrar el catálogo.

## Preguntas abiertas
- ¿Precio de Pan? Se propone $1.800 solo como dato de ejemplo.
- El descuento de stock al vender (y el `409` si no alcanza) se especifica en 0003; el límite en el carrito, en 0002.
- ¿El precio ya incluye impuesto? Se asume que sí (es el "bruto" del ejemplo del spec 0003).
