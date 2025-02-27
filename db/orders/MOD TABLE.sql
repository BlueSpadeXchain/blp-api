CREATE TABLE order_modifications (
    -- Identity and References
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    userid VARCHAR(16) REFERENCES users(userid) ON DELETE CASCADE,
    orderid UUID REFERENCES orders(id) ON DELETE CASCADE,
    leverage NUMERIC(7, 2) NOT NULL CHECK (leverage > 0),
    collateral NUMERIC(20, 6) NOT NULL CONSTRAINT positive_collateral CHECK (collateral > 0),
    -- Price Points
    entry_price NUMERIC(20, 6),
    liq_price NUMERIC(20, 6),
    max_price NUMERIC(20, 6),
    max_value NUMERIC(20, 6),
    close_price NUMERIC(20, 6) DEFAULT NULL,
    -- Optional Order Execution Points
    lim_price NUMERIC(20, 6),
    stop_price NUMERIC(20, 6),
    tp_price NUMERIC(20, 6),
    tp_value NUMERIC(20, 6),
    tp_collateral NUMERIC(20, 6) DEFAULT 0,
    -- Protocol fees
    open_fee NUMERIC(20, 6) DEFAULT 0,
    close_fee NUMERIC(20, 6) DEFAULT 0,
    -- PnL
    pnl NUMERIC(30, 6) DEFAULT 0,
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    signed_at TIMESTAMP,
    canceled_at TIMESTAMP
);