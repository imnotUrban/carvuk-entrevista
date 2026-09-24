CREATE TABLE IF NOT EXISTS usuarios (
    id         BIGSERIAL PRIMARY KEY,
    nombre     VARCHAR(160) NOT NULL,
    correo     VARCHAR(160) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS compras (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES usuarios (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_compras_user_id ON compras (user_id);

CREATE TABLE IF NOT EXISTS compra_items (
    id              BIGSERIAL PRIMARY KEY,
    compra_id       BIGINT NOT NULL REFERENCES compras (id),
    producto_id     BIGINT NOT NULL REFERENCES productos (id),
    nombre          VARCHAR(160) NOT NULL,
    precio_unitario INTEGER NOT NULL CHECK (precio_unitario > 0),
    cantidad        INTEGER NOT NULL CHECK (cantidad BETWEEN 1 AND 99)
);

CREATE INDEX IF NOT EXISTS idx_compra_items_compra_id ON compra_items (compra_id);

CREATE TABLE IF NOT EXISTS boletas (
    id                  BIGSERIAL PRIMARY KEY,
    id_compra           BIGINT NOT NULL REFERENCES compras (id),
    valor_neto          INTEGER NOT NULL CHECK (valor_neto >= 0),
    valor_bruto         INTEGER NOT NULL CHECK (valor_bruto >= 0),
    impuesto            INTEGER NOT NULL CHECK (impuesto >= 0),
    porcentaje_impuesto INTEGER NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (valor_neto + impuesto = valor_bruto)
);

CREATE INDEX IF NOT EXISTS idx_boletas_id_compra ON boletas (id_compra);
