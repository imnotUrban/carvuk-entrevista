# Specs

Cada spec describe un caso de uso **antes** de implementarlo. Un archivo por spec: `NNNN-nombre-corto.md` (ej. `0001-exportar-boilerplate-csv.md`).

## Plantilla

```markdown
# NNNN — Título

- **Estado:** draft | aprobado | implementado
- **Fecha:** YYYY-MM-DD

## Contexto
Qué problema resuelve y por qué.

## Caso de uso
Actor, flujo principal y flujos alternativos.

## Requisitos
- Reglas de negocio, validaciones, errores esperados (status HTTP).

## Diseño
- Backend: capas/archivos a tocar (entity, repository, service, delivery, migración).
- Frontend: páginas/componentes.
- API: método, ruta, request/response de ejemplo.

## Criterios de aceptación
- [ ] Casos verificables (idealmente uno por test).

## Fuera de alcance
Lo que explícitamente NO se hace.

## Preguntas abiertas
Decisiones pendientes.
```
