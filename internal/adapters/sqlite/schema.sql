CREATE TABLE IF NOT EXISTS payments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    public_id TEXT NOT NULL UNIQUE,

    order_id TEXT NOT NULL,
    payer_id INTEGER NOT NULL,

    amount INTEGER NOT NULL,
    currency TEXT NOT NULL,

    status TEXT NOT NULL,

    provider TEXT NOT NULL,
    method TEXT NOT NULL,

    idempotency_key TEXT NOT NULL,

    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    paid_at DATETIME
);

-- uniqueness guarantees
CREATE UNIQUE INDEX IF NOT EXISTS ux_payments_public_id
    ON payments(public_id);

CREATE UNIQUE INDEX IF NOT EXISTS ux_payments_idempotency
    ON payments(idempotency_key);

CREATE INDEX IF NOT EXISTS idx_payments_order_id
    ON payments(order_id);
