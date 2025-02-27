-- long <> short cannot be replaced
-- pair cannot be changed


-- collateral can change
-- entry_price is augemented?
-- close price an be changed
-- liquidation is changed by inference
-- max price is changed by inference
-- max value is changed by inference
-- limit price can only be changed if the limit has not been reached or unsigned state
-- tp price/value/collateral/at is complicated
-- modified at is implict
-- pnl is not reset, so can be complicated
-- open fee is modified

--CREATE TYPE modify_order

CREATE TYPE order_global_update AS (
    current_borrowed NUMERIC,
    current_liquidity NUMERIC,
    current_orders_active NUMERIC,
    current_orders_limit NUMERIC,
    current_orders_pending NUMERIC,
    total_borrowed NUMERIC,
    total_liquidation NUMERIC,
    total_orders_active NUMERIC,
    total_orders_filled NUMERIC,
    total_orders_limit NUMERIC,
    total_orders_liquidated NUMERIC,
    total_orders_stopped NUMERIC,
    total_pnl_losses NUMERIC,
    total_pnl_profits NUMERIC,
    total_revenue NUMERIC,
    treasury_balance NUMERIC,
    total_treasury_profits NUMERIC,
    vault_balance NUMERIC,
    total_vault_profits NUMERIC,
    total_liquidity_rewards NUMERIC,
    total_stake_rewards NUMERIC
);

CREATE TYPE order_update AS (
    order_id UUID,
    user_id VARCHAR(20),
    status VARCHAR(20),
    entry_price NUMERIC(20, 2),
    close_price NUMERIC(20, 2),
    tp_value NUMERIC(20, 2),
    pnl NUMERIC(30,6),
    collateral NUMERIC(30,6),
    tp_at TIMESTAMP,
    balance_change NUMERIC(30,6),
    escrow_balance_change NUMERIC(30,6),
    order_global_update_ order_global_update
);

CREATE TYPE orders_update_batch AS (
    order_id UUID,
    status VARCHAR(20),
    entry_price NUMERIC(20,2),
    close_price NUMERIC(20,2),
    tp_value NUMERIC(20, 2),
    pnl NUMERIC(30,6),
    collateral NUMERIC(30,6),
    modified_at TIMESTAMP,
    ended_at TIMESTAMP,
    tp_at TIMESTAMP
);

CREATE TYPE users_update_batch AS (
    user_id VARCHAR(20),
    balance NUMERIC(30,6),
    escrow_balance NUMERIC(30,6)
);