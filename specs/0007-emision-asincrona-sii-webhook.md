# 0007 — Emisión asíncrona de boletas con API externa + webhook (Bonus 2)

- **Estado:** draft
- **Fecha:** 2026-09-24

## Contexto
Hasta 0003 la boleta se genera solo localmente. Ahora el negocio se integra con un servicio externo que **procesa/timbra documentos** (simula al SII). El servicio es **asíncrono**: se le solicita la emisión y, al terminar, avisa por **webhook**. Este spec cubre el flujo de punta a punta (backend y frontend). Reemplaza lo que 0003 dejó fuera de alcance ("folio/numeración SII, PDF") y extiende 0004 con dos campos nuevos en el historial.

Servicio externo (Pipedream): `POST https://eol9cyu9z953tcx.m.pipedream.net/document-sii`

```json
// request que enviamos
{ "callbackUrl": "https://<url-publica>/webhooks/sii?token=<secreto>", "documentId": 12 }

// webhook que recibimos después (POST a callbackUrl)
{ "documentId": 12, "status": "issued", "siiCode": "SII-123456", "pdfUrl": "https://....pdf" }
```

## Caso de uso
- **Actor:** usuario que compra (y el servicio externo como actor secundario).
- **Flujo principal:**
  1. El usuario genera la boleta (0003). El backend la crea con `sii_status = pending` y responde `201` de inmediato, sin esperar al SII.
  2. Después de confirmar la transacción, el backend hace `POST` al servicio externo con `documentId = boleta.id` y la `callbackUrl` pública.
  3. El servicio externo, cuando termina, llama al webhook `POST /webhooks/sii`.
  4. El backend valida el webhook y actualiza la boleta a `issued` con `sii_code` y `pdf_url`.
  5. El front, que estaba mostrando la boleta como "Pendiente", se actualiza sola y muestra **Número SII** y **Link PDF**.
- **Alternativos:**
  - El POST al servicio externo falla (red, timeout, `5xx`) → la boleta pasa a `failed`; el usuario puede pulsar "Reintentar emisión".
  - El webhook llega repetido → se responde `200` sin cambiar nada (idempotente).
  - El webhook llega con `documentId` inexistente, token inválido o body inválido → se rechaza (`404` / `401` / `400`).
  - El servicio informa un `status` distinto de `issued` → la boleta pasa a `failed`.

## Requisitos
### Estados de emisión
`sii_status`: `pending` (default al crear) → `issued` | `failed`. `failed` → `pending` solo mediante reintento. `issued` es final e inmutable.

### Solicitud al servicio externo
- Se hace **después** del commit de la transacción de 0003 y **fuera** de ella (una llamada HTTP no debe mantener abierta una transacción ni revertir la boleta si falla).
- Se ejecuta en una goroutine desde el service (`RequestIssuance(ctx, boletaID)`); el `201` de `POST /boletas` no espera su resultado.
- Cliente HTTP con timeout de 10 s. Respuesta `2xx` = solicitud aceptada (la boleta sigue `pending`, se espera el webhook); cualquier otro resultado o error → `failed` y log del error.
- `documentId` es el `id` de la boleta. `callbackUrl` = `PUBLIC_BASE_URL + "/webhooks/sii?token=" + SII_WEBHOOK_TOKEN`.
- El cliente vive detrás de una interfaz (`siiclient.Client`) para poder fakearlo en tests.

### Webhook
- `POST /webhooks/sii`, **fuera** de `/api/v1` y **sin** autenticación de Clerk (0006): el servicio externo no tiene sesión.
- Autenticación del webhook: query param `token` que debe coincidir con `SII_WEBHOOK_TOKEN` (comparación en tiempo constante); si no → `401`.
- Body: `documentId` (uint), `status` (string), `siiCode` (string), `pdfUrl` (string).
- Validaciones → `400`: `documentId` presente; `status` en `issued` | `failed` (u otro valor de error conocido → `failed`; vacío → `400`); si `status = issued`, `siiCode` no vacío (máx. 64) y `pdfUrl` con esquema `http` o `https` (evita links `javascript:`).
- `documentId` inexistente → `404`.
- **Idempotencia:** si la boleta ya está `issued`, responde `200` sin modificarla (aunque el body difiera).
- Éxito → `200 {"ok": true}` rápido (sin trabajo pesado dentro del handler).

### Reintento
- `POST /api/v1/boletas/:id/reintentar-emision` (autenticado, 0006): solo si la boleta es visible para el usuario según la regla de 0006 (propia o histórica de "Cliente Demo"; si no → `404`) y solo si `sii_status = failed` (si no → `409`). Pasa a `pending` y vuelve a solicitar la emisión.

### Datos y compatibilidad
- Boletas anteriores a este spec: la migración las marca `failed` (nunca se solicitó su emisión) para que puedan reintentarse.
- Los campos nuevos aparecen en la respuesta de `POST /boletas` (0003), en la lista y el detalle (0004).

### Configuración
`backend/.env.example`:
- `SII_API_URL` (default `https://eol9cyu9z953tcx.m.pipedream.net/document-sii`)
- `PUBLIC_BASE_URL` (URL pública del backend, p. ej. la de ngrok: `ngrok http 8080`)
- `SII_WEBHOOK_TOKEN` (secreto aleatorio)

## Diseño
- **Migración** `0005_add_sii_to_boletas.up/down.sql`: `sii_status VARCHAR(20) NOT NULL DEFAULT 'pending'`, `sii_code VARCHAR(64) NULL`, `pdf_url TEXT NULL`, `CHECK (sii_status IN ('pending','issued','failed'))`; `UPDATE boletas SET sii_status = 'failed'` para filas existentes. En dev, `AutoMigrate` de la entity.
- **Backend** (patrón por capas):
  - `pkg/siiclient/` (interfaz `Client` + implementación HTTP + tests con `httptest`).
  - `entity/boleta`: campos `SiiStatus`, `SiiCode *string`, `PdfURL *string`.
  - `repository/boleta`: `update.go` con `MarkIssued`, `MarkFailed`, `MarkPending`; las transiciones son condicionales (`UPDATE ... WHERE id = ? AND sii_status <> 'issued'`) para no pisar una boleta emitida.
  - `service/boleta`: `create.go` dispara `RequestIssuance` tras el commit; `issuance.go` (solicitud, reintento, `HandleWebhook`).
  - `delivery/boleta`: `webhook.go` (handler público) y `retry.go`; `delivery/router.go` registra `/webhooks/sii` fuera del grupo autenticado.
  - `postman_collection.json`: requests de webhook y reintento.
- **Frontend:**
  - `lib/types.ts`: `siiStatus`, `siiCode`, `pdfUrl` en boleta.
  - `app/boletas/page.tsx`: columnas **Número SII** y **PDF** + badge de estado (Pendiente / Emitida / Fallida).
  - `app/boletas/[id]/page.tsx`: mismo bloque, botón "Reintentar emisión" cuando está `failed`.
  - Actualización automática: mientras haya boletas `pending` visibles, consulta cada 3 s (`setInterval` con limpieza) y se detiene cuando ninguna está `pending`.
  - El link PDF abre en pestaña nueva (`target="_blank" rel="noopener noreferrer"`); mientras no exista, muestra "—".
- **API:**

```json
// GET /api/v1/boletas/12 → 200 (campos nuevos)
{ "id": 12, "sii_status": "issued", "sii_code": "SII-123456", "pdf_url": "https://....pdf", "...": "..." }

// POST /webhooks/sii?token=<secreto>
{ "documentId": 12, "status": "issued", "siiCode": "SII-123456", "pdfUrl": "https://....pdf" }
// 200 {"ok": true}

// POST /api/v1/boletas/12/reintentar-emision → 202 (boleta en pending)
```

## Criterios de aceptación
**Backend (tests con `siiclient` fake y `testutil.SetupDB`):**
- [ ] Al crear una boleta queda `pending`, sin `sii_code` ni `pdf_url`, y `POST /boletas` responde sin esperar al servicio externo.
- [ ] Tras crear la boleta se invoca al cliente SII con `documentId` = id de la boleta y la `callbackUrl` configurada (incluye el token).
- [ ] Si el cliente SII falla o da timeout, la boleta queda `failed` y la boleta y el stock de 0003 **no** se revierten.
- [ ] Cliente HTTP real (con `httptest`): envía el JSON `{callbackUrl, documentId}` por POST y trata `2xx` como aceptado y otros status como error.
- [ ] Webhook con token válido y `status: issued` deja la boleta `issued` con `sii_code` y `pdf_url`.
- [ ] Webhook sin token o con token incorrecto → `401` y la boleta no cambia.
- [ ] Webhook con `documentId` inexistente → `404`.
- [ ] Webhook con body inválido (sin `documentId`, `status` vacío, `issued` sin `siiCode`, `pdfUrl` con esquema no http/https) → `400`.
- [ ] Webhook repetido sobre una boleta `issued` → `200` y los datos no cambian.
- [ ] Webhook con estado de error deja la boleta `failed`; un webhook de error no pisa una boleta ya `issued`.
- [ ] El webhook no requiere sesión de Clerk (0006), mientras que `/api/v1/*` sigue exigiéndola.
- [ ] Reintento: boleta `failed` propia → `202` y `pending`; boleta `pending` o `issued` → `409`; boleta de otro usuario real → `404`; boleta histórica de "Cliente Demo" `failed` → `202` para cualquier usuario autenticado.
- [ ] `GET /boletas` y `GET /boletas/:id` incluyen `sii_status`, `sii_code` y `pdf_url`.

**Frontend (verificación manual, con ngrok + servicio real):**
- [ ] Al generar una boleta, el historial la muestra "Pendiente" y, sin recargar, pasa a "Emitida" con **Número SII** y **Link PDF** cuando llega el webhook.
- [ ] El link PDF abre el documento en una pestaña nueva.
- [ ] El detalle de la boleta muestra Número SII y Link PDF.
- [ ] Una boleta `failed` muestra "Fallida" y el botón "Reintentar emisión" la devuelve a "Pendiente".
- [ ] Cuando no quedan boletas pendientes, el auto-refresco se detiene.

## Fuera de alcance
- Firma criptográfica del webhook (HMAC); se usa un token compartido en la URL.
- Cola persistente/reintentos automáticos con backoff; el reintento es manual.
- Reanudar solicitudes `pending` que se perdieron si el backend se reinicia a mitad de camino (queda `pending`; se puede resolver con un job futuro).
- WebSockets/SSE; el front usa polling.
- Descargar/almacenar el PDF localmente; solo se guarda el link.
- Anulación de documentos en el SII.

## Preguntas abiertas
- **`callbackUrl` pública:** el backend corre local, así que se necesita ngrok (o similar) y actualizar `PUBLIC_BASE_URL` cada vez que cambie la URL. ¿Confirmas que se usa ngrok?
- **Token en la query:** Pipedream llama a la `callbackUrl` tal cual se la enviamos (incluida la query); se asume que sí y hay que verificarlo con una prueba real.
- ¿Qué `status` de error devuelve el servicio real además de `issued`? Se asume que cualquier otro valor no vacío significa fallo.
- Boletas previas a este spec: se propone marcarlas `failed` para poder reintentarlas. ¿Prefieres dejarlas sin estado SII?
