CREATE TABLE IF NOT EXISTS productos (
    id         BIGSERIAL PRIMARY KEY,
    nombre     VARCHAR(160) NOT NULL,
    precio     INTEGER NOT NULL CHECK (precio > 0),
    stock      INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_productos_nombre ON productos (nombre);
CREATE INDEX IF NOT EXISTS idx_productos_deleted_at ON productos (deleted_at);
