# Carvuk — Entrevista

CRUD de **boilerplate** (tabla genérica de ejemplo con 5 campos inventados) con backend en Go y frontend en Next.js. Sirve como plantilla para agregar nuevas entidades.

## Estructura del repositorio

```
entrevista/
├── backend/                    # Go + Gin + GORM + PostgreSQL
│   ├── main.go                 # solo: conecta DB, migra, arma el router
│   ├── pkg/                      # código compartido entre entidades
│   │   ├── apperr/                # errores de dominio compartidos (ErrNotFound, ErrValidation, ErrDuplicateCode)
│   │   ├── dbutil/                # helpers GORM compartidos (ILIKE vs LIKE case-insensitive)
│   │   ├── httpkit/               # tipos de respuesta HTTP + mapeo de error -> status
│   │   └── testutil/              # DB SQLite en memoria para los tests de service
│   ├── entity/boilerplate/boilerplate.go   # struct de dominio (GORM model)
│   ├── repository/boilerplate/    # GORM — solo queries, sin reglas de negocio
│   │   ├── repository.go          # interface Repository, Filter, struct, NewRepository
│   │   ├── create.go / find.go (GetByID + ExistsByCode) / list.go / update.go / delete.go (soft delete)
│   ├── service/boilerplate/       # lógica de negocio y validaciones
│   │   ├── service.go             # interface Service, whitelists (sort_by, status), NewService
│   │   ├── create.go / find.go / list.go / update.go / delete.go
│   │   └── boilerplate_test.go    # tests del service
│   ├── delivery/                  # handlers Gin — solo HTTP (bind, status codes)
│   │   ├── router.go              # arma el *gin.Engine llamando a NewModule de cada entidad
│   │   └── boilerplate/           # handler.go (RegisterRoutes), factory.go (NewModule), create/find/list/update/delete.go
│   ├── migrations/                # SQL versionado (documentación del esquema)
│   │   └── 0001_create_boilerplate.up.sql / .down.sql
│   ├── postman_collection.json
│   ├── .air.toml                 # hot-reload en dev (air)
│   ├── .env.example
│   └── go.mod
├── frontend/                    # Next.js (App Router) + Tailwind + shadcn/ui
│   ├── app/boilerplate/page.tsx # CRUD (filtros, tabla, dialog)
│   ├── components/
│   │   ├── nav-bar.tsx
│   │   ├── boilerplate-form-dialog.tsx
│   │   └── ui/                   # componentes shadcn
│   └── lib/
│       ├── api.ts                # cliente fetch tipado hacia el backend
│       └── types.ts
├── docker-compose.yml            # SOLO levanta PostgreSQL
└── CLAUDE.md
```

## Arquitectura del backend (capas + factory por entidad)

Se mantienen las 4 capas (`entity` → `repository` → `service` → `delivery`), cada una con una **subcarpeta por entidad** (hoy solo `boilerplate/`), y dentro de cada subcarpeta el código está partido **un archivo por operación** (`create.go`, `find.go`, `list.go`, `update.go`, `delete.go`). Cada capa expone su propio `interface` + struct privado + constructor (`NewRepository`, `NewService`, `NewHandler`) en el archivo base de la carpeta.

- **entity/boilerplate**: struct de dominio puro con tags de GORM. Tabla `boilerplate`.
- **repository/boilerplate**: solo arma queries de GORM (Where, Count, Order, Limit/Offset). No valida ni decide reglas de negocio. Usa `dbutil.CaseInsensitiveLike`.
- **service/boilerplate**: valida inputs y aplica reglas de negocio (`code` único, `status` dentro de `draft|active|archived`, `quantity`/`amount` no negativos, whitelist de columnas para `sort_by`), traduce `gorm.ErrRecordNotFound` a errores de `apperr`.
- **delivery/boilerplate**: handlers de Gin. Parsean query params / body, llaman al service y usan `httpkit.HandleServiceError` para mapear errores a status HTTP. `factory.go` expone `NewModule(db *gorm.DB) *Handler` (repo → service → handler). `delivery/router.go` solo llama a `NewModule(db)` y registra rutas.
- **pkg/**: paquetes chicos compartidos, para evitar import cycles entre capas/entidades.

### Campos de `boilerplate`

| Campo      | Tipo            | Notas                                                    |
|------------|-----------------|----------------------------------------------------------|
| `id`       | uint            | PK                                                       |
| `name`     | varchar(160)    | requerido                                                |
| `code`     | varchar(64)     | requerido, único entre registros no borrados             |
| `status`   | varchar(20)     | `draft` (default) \| `active` \| `archived`              |
| `quantity` | int             | >= 0, default 0                                          |
| `amount`   | numeric(12,2)   | >= 0, default 0                                          |
| timestamps | `created_at`, `updated_at`, `deleted_at` (soft delete)                        |

### Soft delete

La entidad usa `gorm.DeletedAt` (columna `deleted_at`). `DELETE` en la API **nunca borra la fila físicamente**: hace `UPDATE ... SET deleted_at = now()`. Los listados/Get filtran automáticamente `deleted_at IS NULL`.

## Endpoints

Base URL: `http://localhost:8080/api/v1`

| Método | Ruta               | Descripción                     |
|--------|--------------------|---------------------------------|
| POST   | `/boilerplate`     | Crear                           |
| GET    | `/boilerplate`     | Listar con filtros y paginación |
| GET    | `/boilerplate/:id` | Obtener uno                     |
| PUT    | `/boilerplate/:id` | Actualizar (parcial)            |
| DELETE | `/boilerplate/:id` | Soft delete                     |

Filtros de `GET /boilerplate`: `name`, `code` (contiene, case-insensitive), `status`, `min_amount`, `max_amount`, `page`, `limit` (máx 100), `sort_by` (`id`,`name`,`amount`,`quantity`,`created_at`,`updated_at`), `order` (`asc`,`desc`).

Respuestas de error: `400` validación, `404` no encontrado, `409` código duplicado.

## Base de datos

- **Migraciones SQL** en `backend/migrations/` (`0001_create_boilerplate`): `CREATE TABLE`, índices y un índice único parcial `idx_boilerplate_code_active` (`code` único solo entre registros no borrados, así se puede reutilizar el código de uno soft-deleteado).
- En desarrollo, `main.go` corre `AutoMigrate` de GORM al arrancar y luego ejecuta explícitamente el `CREATE UNIQUE INDEX ... WHERE deleted_at IS NULL` (GORM no puede expresarlo con `AutoMigrate`), espejando `0001_create_boilerplate.up.sql`. `entity.Boilerplate.Code` **no** tiene `uniqueIndex` de GORM (rompería el soft delete); la unicidad se garantiza con el índice parcial + validación en `service`.
- El SQL en `migrations/` es la referencia canónica del esquema; si se reemplaza `AutoMigrate` por `golang-migrate`, esos archivos sirven tal cual.

## Cómo correr el proyecto

1. Levantar PostgreSQL: `docker compose up -d`
2. Backend:
   ```
   cd backend
   cp .env.example .env
   go run .
   ```
   Corre en `http://localhost:8080`. Health check: `GET /health`.

   **Modo dev con hot-reload:** [air](https://github.com/air-verse/air), configurado en `backend/.air.toml`:
   ```
   go install github.com/air-verse/air@latest   # una sola vez
   cd backend
   air
   ```
3. Frontend:
   ```
   cd frontend
   cp .env.local.example .env.local
   npm install
   npm run dev
   ```
   Corre en `http://localhost:3000`, página `/boilerplate`.

## Tests del backend

```
cd backend
go test ./...
```

`service/boilerplate/boilerplate_test.go` (paquete `boilerplate_test`): create (y estado default `draft`), validaciones (nombre/código vacío, estado inválido, cantidad/monto negativos), código duplicado, get, get-not-found, list con filtros (nombre, código, estado, rango de monto, combinados), paginación, update parcial, update-not-found, update rechaza código duplicado, soft delete (deja de listarse/obtenerse pero la fila sigue en la DB) y delete-not-found.

Usan `testutil.SetupDB(t)` (SQLite en memoria vía `AutoMigrate`) en vez de PostgreSQL, y `dbutil.CaseInsensitiveLike` emite `ILIKE` en Postgres y `LOWER(...) LIKE LOWER(...)` en otros dialectos.
