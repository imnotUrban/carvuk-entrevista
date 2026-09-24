CREATE TABLE IF NOT EXISTS boilerplate (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(160) NOT NULL,
    code       VARCHAR(64) NOT NULL,
    status     VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'archived')),
    quantity   INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    amount     NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (amount >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_boilerplate_name ON boilerplate (name);
CREATE INDEX IF NOT EXISTS idx_boilerplate_status ON boilerplate (status);
CREATE INDEX IF NOT EXISTS idx_boilerplate_deleted_at ON boilerplate (deleted_at);

-- Unique code only among non-deleted rows, so a soft-deleted row's code can
-- be reused by a new row.
CREATE UNIQUE INDEX IF NOT EXISTS idx_boilerplate_code_active
    ON boilerplate (code)
    WHERE deleted_at IS NULL;
