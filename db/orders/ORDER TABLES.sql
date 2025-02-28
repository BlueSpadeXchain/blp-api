-- Table for storing synthetic perpetual trading orders
CREATE TABLE orders (
    -- Identity and References
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    userid VARCHAR(16) REFERENCES users(userid) ON DELETE CASCADE,
    
    -- Core Order Properties
    order_type VARCHAR(10) NOT NULL CHECK (order_type IN ('long', 'short')),
    leverage NUMERIC(7, 2) NOT NULL CHECK (leverage > 0),
    pair_id VARCHAR(64) NOT NULL,
    collateral NUMERIC(20, 6) NOT NULL CONSTRAINT positive_collateral CHECK (collateral > 0),
    
    -- Order Status
    status VARCHAR(20) NOT NULL DEFAULT 'unsigned' 
        CHECK (status IN ('unsigned', 'pending', 'limit', 'filled', 'canceled', 'closed', 'liquidated', 'stopped')),
    
    -- Price Points
    entry_price NUMERIC(20, 6),
    liq_price NUMERIC(20, 6),
    max_price NUMERIC(20, 6),
    max_value NUMERIC(20, 6) NOT NULL GENERATED ALWAYS AS (collateral * 10) STORED,
    close_price NUMERIC(20, 6) DEFAULT NULL,
    
    -- Optional Order Execution Points
    lim_price NUMERIC(20, 6),
    stop_price NUMERIC(20, 6),
    tp_price NUMERIC(20, 6),
    tp_value NUMERIC(20, 6),
    tp_collateral NUMERIC(20, 6) DEFAULT 0,

    -- Protocol fees
    open_fee NUMERIC(20, 6) DEFAULT 0;
    close_fee NUMERIC(20, 6) DEFAULT 0;
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    started_at TIMESTAMP,
    signed_at TIMESTAMP,
    modified_at TIMESTAMP,
    tp_at TIMESTAMP,
    ended_at TIMESTAMP,

    -- Complex Constraints
    CONSTRAINT valid_limit_price CHECK (
        lim_price IS NULL OR lim_price != 0
    ),
    CONSTRAINT valid_stop_price CHECK (
        stop_price IS NULL OR
        (stop_price != 0 AND liq_price IS NOT NULL AND stop_price > liq_price AND entry_price IS NOT NULL AND stop_price < entry_price)
    ),
    CONSTRAINT valid_tp CHECK (
        (tp_price IS NULL AND tp_value IS NULL AND tp_collateral = 0) OR
        (
            tp_price IS NOT NULL AND tp_value IS NOT NULL AND 
            max_price IS NOT NULL AND max_price > tp_price AND 
            entry_price IS NOT NULL AND tp_price > entry_price AND
            tp_collateral >= 0 AND tp_collateral < collateral AND
            max_value > tp_value
        )
    ),
    CONSTRAINT valid_timestamps CHECK (
        (started_at IS NULL OR started_at >= created_at) AND
        (signed_at IS NULL OR signed_at >= created_at) AND
        (modified_at IS NULL OR modified_at >= signed_at) AND
        (ended_at IS NULL OR ended_at >= COALESCE(signed_at, created_at, modified_at))
    )
);

-- Table description
COMMENT ON TABLE orders2 IS 'Stores synthetic perpetual trading orders with execution parameters, including limit orders, stop losses, and take profits';

-- Column comments
COMMENT ON COLUMN orders.id IS 'Unique identifier for the order';
COMMENT ON COLUMN orders.userid IS 'Reference to the user who created the order';
COMMENT ON COLUMN orders.order_type IS 'Type of order: long (buy) or short (sell)';
COMMENT ON COLUMN orders.leverage IS 'Trading leverage multiplier, must be positive';
COMMENT ON COLUMN orders.pair IS 'Trading pair identifier (e.g., BTC-USD)';
COMMENT ON COLUMN orders.collateral IS 'Amount of collateral provided for the trade, must be positive';
COMMENT ON COLUMN orders.status IS 'Current order status: unsigned -> pending -> filled/canceled/closed/liquidated';
COMMENT ON COLUMN orders.entry_price IS 'Price at which the order was executed';
COMMENT ON COLUMN orders.liq_price IS 'Price at which the position will be liquidated';
COMMENT ON COLUMN orders.max_price IS 'Maximum price threshold for the order (10x)';
COMMENT ON COLUMN orders.max_value IS '10x collateral input';
COMMENT ON COLUMN orders.close_price IS 'regardless of the source, mark price when closed';
COMMENT ON COLUMN orders.lim_price IS 'Optional limit price for order execution';
COMMENT ON COLUMN orders.stop_price IS 'Optional stop-loss price, must be between liq_price and entry_price';
COMMENT ON COLUMN orders.tp_price IS 'Optional take-profit price target, requires tp_value';
COMMENT ON COLUMN orders.tp_value IS 'Optional take-profit value, required if tp_price is set';
COMMENT ON COLUMN orders.tp_collateral IS 'Optional take-profit collateral, required if tp_price is set';
COMMENT ON COLUMN orders.created_at IS 'Timestamp when the order was created';
COMMENT ON COLUMN orders.started_at IS 'Timestamp when a order begins execution';
COMMENT ON COLUMN orders.signed_at IS 'Timestamp when the order was signed by the user';
COMMENT ON COLUMN orders.modified_at IS 'Timestamp when a order is successfully modified';
COMMENT ON COLUMN orders.tp_at IS 'Timestamp when a order is successfully taken profit';
COMMENT ON COLUMN orders.ended_at IS 'Timestamp when the order was completed, canceled, or liquidated';
COMMENT ON COLUMN orders.open_fee IS 'Fee when a order is successfully opened';
COMMENT ON COLUMN orders.close_fee IS 'Fee when a order is successfully closed';

-- Indexes for common queries
CREATE INDEX idx_orders2_userid ON orders2(userid);
CREATE INDEX idx_orders2_status ON orders2(status);
CREATE INDEX idx_orders2_created_at ON orders2(created_at);
