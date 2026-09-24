# 0002 — Carrito de compras (solo frontend)

- **Estado:** implementado (compra simulada; la boleta real y el 409 dependen de 0003)
- **Fecha:** 2026-09-24

## Contexto

Antes de generar una boleta (0003), el usuario arma un carrito con productos del catálogo (0001). El carrito **vive solo en el frontend**: no hay tabla ni endpoints de carrito.

## Caso de uso

- **Actor:** usuario que compra.
- **Flujo principal:** en `/tienda` ve el catálogo, pulsa "Agregar" en un producto; abre el carrito y ve cada línea (producto, precio unitario, cantidad, subtotal) y el total; quita productos si se arrepiente; pulsa "Comprar", que **simula** la compra (ver Requisitos). La generación real de la boleta llegará con 0003.
- **Alternativos:**
  - Agregar un producto que ya está en el carrito incrementa su cantidad.
  - Cambiar la cantidad con +/−; al llegar a 0 la línea se elimina.
  - Quitar una línea completa con "Quitar".
  - "Vaciar carrito".
  - Recargar la página no pierde el carrito.
  - Carrito vacío: el botón "Comprar" queda deshabilitado y se muestra un mensaje.

## Requisitos

- Estado del carrito: lista de líneas `{ producto_id, nombre, precio, stock, cantidad }`, con `cantidad` entero `>= 1` y `<= min(99, stock)`. `stock` es una copia del catálogo al momento de agregar.
- Productos con `stock = 0` se muestran como "Sin stock" y su botón "Agregar" queda deshabilitado.
- Al llegar la cantidad de una línea al stock, el botón "+" / "Agregar" se deshabilita y se indica "Máximo disponible".
- El stock del front es solo una guía de UX (no se reserva): el backend valida el stock real al generar la boleta y puede responder `409` (0003); en ese caso el front muestra el mensaje, conserva el carrito y refresca el catálogo.
- **Compra simulada:** el botón "Comprar" no llama al backend. Muestra un estado "Procesando..." breve (~800 ms), luego un aviso de éxito con el total pagado, vacía el carrito y lo persiste vacío. No descuenta stock ni crea boleta. Cuando se implemente 0003, este botón pasará a "Generar boleta" y llamará al backend.
- Se guarda en `localStorage` (clave `carvuk_cart`) para sobrevivir recargas; se lee dentro de `try/catch` y, si el JSON es inválido, se parte con carrito vacío.
- El total mostrado es `Σ precio × cantidad`. Es **solo informativo**: el backend recalcula todo con los precios vigentes (0003).
- Los precios del carrito son una copia; si un producto cambia de precio o se elimina, el backend prevalece al generar la boleta (ver 0003, errores `400`/`409`).

## Diseño

- **Backend:** ninguno (usa `GET /productos` de 0001).
- **Frontend:**
  - `app/tienda/page.tsx`: catálogo (usa `GET /productos`) con botón "Agregar" + panel/`Sheet` del carrito.
  - `lib/cart.tsx`: `CartProvider` + hook `useCart()` con `items`, `add(producto)`, `setQuantity(id, n)`, `remove(id)`, `clear()`, `total`, `count`. Montado en `app/layout.tsx`.
  - `components/cart-panel.tsx`, y contador de ítems en `nav-bar.tsx`.
  - Formato de moneda CLP compartido en `lib/format.ts` (`Intl.NumberFormat("es-CL")`).

## Criterios de aceptación

- [x] Agregar un producto nuevo crea una línea con cantidad 1.
- [x] Agregar el mismo producto dos veces deja una línea con cantidad 2.
- [x] Cambiar cantidad recalcula subtotal y total; cantidad 0 elimina la línea.
- [x] Quitar una línea la elimina y actualiza el total.
- [x] Ejemplo del enunciado: 2 Leches ($1.100) + 1 Aceite ($2.000) → total $4.200.
- [x] No se puede agregar ni subir la cantidad por encima del stock del producto.
- [x] Un producto con stock 0 aparece "Sin stock" y no se puede agregar.
- [ ] Si el backend responde `409` por stock insuficiente al generar la boleta, se muestra el error y el carrito se conserva.
- [x] El carrito persiste tras recargar la página.
- [x] `localStorage` corrupto o no disponible no rompe la página (carrito vacío).
- [x] Con carrito vacío, "Comprar" está deshabilitado.
- [x] "Comprar" simula la compra: muestra éxito con el total, vacía el carrito y no llama al backend.

## Fuera de alcance

- Compra real (boleta, descuento de stock): spec 0003.
- Persistir el carrito en el backend o por usuario.
- Reserva de stock mientras el producto está en el carrito.
- Cupones, envío.
- Tests automatizados de frontend (se verifican manualmente; el repo no tiene runner de tests en el front).
