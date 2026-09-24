# 0006 — Autenticación con Clerk (Bonus 1)

- **Estado:** draft
- **Fecha:** 2026-09-24

## Contexto
La app (productos, tienda/carrito, boletas) hoy es pública y todas las compras se asocian al usuario semilla "Cliente Demo" (0003). Se protege con autenticación de [Clerk](https://clerk.com/), **manteniendo lo ya construido**: mismas pantallas y mismos endpoints, ahora exigiendo sesión. Cada compra/boleta pasa a pertenecer al usuario autenticado.

## Caso de uso
- **Actor:** persona que usa la app.
- **Flujo principal:** entra a cualquier página → si no tiene sesión, es redirigida a `/sign-in` → inicia sesión (o se registra en `/sign-up`) → vuelve a la página que quería → ve su nombre/avatar en la barra (`UserButton`) y puede cerrar sesión.
- **Alternativos:**
  - Llamar a la API sin token o con token inválido/expirado → `401`.
  - Primera vez que un usuario nuevo llama a la API → se crea automáticamente su fila en `usuarios` (ver Requisitos).
  - Un usuario intenta ver una boleta de otro usuario real → `404`. Las boletas históricas de "Cliente Demo" son la excepción: las ve cualquier usuario autenticado.
  - Cerrar sesión limpia el acceso; el carrito local no se comparte entre usuarios (ver Requisitos).

## Requisitos
### Frontend (Next.js 16, App Router)
- Paquete `@clerk/nextjs`; `<ClerkProvider>` en `app/layout.tsx`.
- `proxy.ts` (en Next 16 reemplaza a `middleware.ts`) con `clerkMiddleware`: rutas públicas solo `/sign-in(.*)` y `/sign-up(.*)`; todo lo demás llama a `auth.protect()`.
- Páginas `app/sign-in/[[...sign-in]]/page.tsx` y `app/sign-up/[[...sign-up]]/page.tsx` con los componentes `<SignIn />` / `<SignUp />`.
- `nav-bar.tsx`: muestra `<UserButton />` con sesión; el `nav-bar` no se muestra en las pantallas de sign-in/sign-up.
- `lib/api.ts`: cada request al backend envía `Authorization: Bearer <token>` (token de sesión de Clerk, obtenido con `getToken()`). Si la respuesta es `401`, redirige a `/sign-in`.
- El carrito (`localStorage`, 0002) se guarda con la clave `carvuk_cart:<clerk_user_id>`, para que otro usuario en el mismo navegador no vea el carrito ajeno; al cerrar sesión no se muestra.

### Backend (Go + Gin)
- SDK `github.com/clerk/clerk-sdk-go/v2`: se configura con `CLERK_SECRET_KEY` y el middleware verifica el JWT de sesión del header `Authorization: Bearer`.
- Nuevo paquete `pkg/auth`:
  - Middleware Gin `auth.Required()` → `401` (`{"error": "..."}` con el formato de `httpkit`) si falta/está inválido el token.
  - `auth.UserID(c)` devuelve el id interno de `usuarios` del usuario autenticado, inyectado por el middleware.
  - El verificador de tokens es una interfaz para poder reemplazarlo por un fake en los tests.
- Rutas protegidas: todo `/api/v1/*` (productos, boletas). Públicas: `GET /health` y, cuando exista, el webhook de 0007 (que se autentica por su propio mecanismo).
- **Usuarios (provisionamiento al primer uso):** la tabla `usuarios` (0003) agrega `clerk_id varchar(64)` **único**, nullable. Si el `sub` del token no existe en `usuarios`, el middleware crea la fila con `nombre` y `correo` obtenidos de la API de Clerk (`user.Get`). Solo se hace la llamada a Clerk la primera vez.
- **Boletas por usuario:** `POST /boletas` asocia la compra al usuario autenticado (se elimina el uso del usuario semilla; la fila "Cliente Demo" se conserva para las boletas históricas, con `clerk_id` nulo). `GET /boletas` lista las boletas del usuario autenticado **más las históricas** (las de usuarios con `clerk_id` nulo, o sea "Cliente Demo"), que son visibles para todos. `GET /boletas/:id` devuelve `404` si la boleta es de otro usuario real (no `403`, para no revelar su existencia). Regla de visibilidad única, reutilizada por lista, detalle y reintento (0007): `usuarios.clerk_id = <usuario autenticado> OR usuarios.clerk_id IS NULL`.
- **Productos:** cualquier usuario autenticado puede usar el CRUD (no hay roles; ver Fuera de alcance).
- CORS: el backend debe permitir el header `Authorization` desde el origen del frontend.

### Configuración
- `frontend/.env.local.example`: `NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY`, `CLERK_SECRET_KEY`, `NEXT_PUBLIC_CLERK_SIGN_IN_URL=/sign-in`, `NEXT_PUBLIC_CLERK_SIGN_UP_URL=/sign-up`.
- `backend/.env.example`: `CLERK_SECRET_KEY`.
- Las claves reales **no** se commitean (`.env` y `.env.local` en `.gitignore`); se usa una instancia de desarrollo de Clerk.

## Diseño
- **Migración** `0004_add_clerk_id_to_usuarios.up/down.sql`: `ALTER TABLE usuarios ADD COLUMN clerk_id VARCHAR(64)` + índice único parcial sobre `clerk_id IS NOT NULL`. En dev, `AutoMigrate` de la entity `Usuario`.
- **Backend:** `pkg/auth/` (middleware, verificador, `UserID`), `repository/usuario/` (`find.go` por `clerk_id`, `create.go`), `service/usuario/` (`FindOrCreateByClerk`), cambios en `delivery/router.go` (grupo `/api/v1` con `auth.Required()`), `service/boleta` (recibe `userID` en `Create`, `List` y `Find`), `postman_collection.json` (variable de token).
- **Frontend:** `proxy.ts`, `app/sign-in/...`, `app/sign-up/...`, `app/layout.tsx`, `components/nav-bar.tsx`, `lib/api.ts`, `lib/cart.tsx` (clave por usuario).
- **API:** los mismos endpoints con header `Authorization: Bearer <jwt>`; nuevo status `401`.

```
GET /api/v1/boletas            (sin token)     → 401 {"error":"unauthorized"}
GET /api/v1/boletas            (token válido)  → 200, solo boletas del usuario
GET /api/v1/boletas/7          (boleta ajena)  → 404
```

## Criterios de aceptación
**Backend (tests con verificador fake):**
- [ ] Request sin `Authorization` a un endpoint de `/api/v1/*` → `401`.
- [ ] Token inválido o expirado → `401`.
- [ ] `GET /health` responde `200` sin token.
- [ ] Primer request de un `sub` nuevo crea la fila en `usuarios` con `clerk_id`, nombre y correo; el segundo request no crea otra.
- [ ] `clerk_id` es único: dos requests concurrentes del mismo usuario nuevo no generan duplicados.
- [ ] `POST /boletas` asocia la compra al usuario autenticado, no al usuario semilla.
- [ ] `GET /boletas` devuelve las boletas del usuario autenticado y las históricas de "Cliente Demo", y no las de otros usuarios reales.
- [ ] `GET /boletas/:id` de una boleta de otro usuario real → `404`.
- [ ] `GET /boletas/:id` de una boleta histórica de "Cliente Demo" → `200` para cualquier usuario autenticado.

**Frontend (verificación manual):**
- [ ] Sin sesión, cualquier ruta (incluida `/`) redirige a `/sign-in`; tras iniciar sesión vuelve a la ruta original.
- [ ] Se puede registrar un usuario nuevo en `/sign-up` y entrar a la app.
- [ ] La barra muestra `UserButton` y permite cerrar sesión; tras cerrar sesión vuelve a `/sign-in`.
- [ ] Todo el flujo existente (productos, tienda/carrito, generar boleta, historial) funciona con sesión iniciada.
- [ ] Si el backend responde `401`, el front redirige a `/sign-in`.
- [ ] Dos usuarios en el mismo navegador no comparten carrito.

## Fuera de alcance
- Roles/permisos (admin vs. cliente) y restricción del CRUD de productos por rol.
- Webhooks de Clerk para sincronizar usuarios (se usa provisionamiento al primer uso).
- MFA, organizaciones, otros proveedores: solo se usan los métodos habilitados en la instancia (email y Google); se configuran en el dashboard de Clerk, sin código extra porque `<SignIn />`/`<SignUp />` los muestran automáticamente.
- Migrar las boletas del usuario "Cliente Demo" a un usuario real.
- Tests automatizados de frontend.

## Decisiones
- **Métodos de login:** email y Google.
- **Boletas históricas de "Cliente Demo":** visibles para todos los usuarios autenticados (no se reasignan ni se ocultan).

## Preguntas abiertas
Ninguna.
