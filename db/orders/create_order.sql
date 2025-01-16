CREATE OR REPLACE FUNCTION create_order(
    user_id VARCHAR,
    order_type VARCHAR,
    leverage NUMERIC(5, 2),
    pair VARCHAR,
    collateral NUMERIC(20, 6),
    entry_price NUMERIC(20, 6),
    liq_price NUMERIC(20, 6)
) RETURNS orders AS $$
DECLARE
    total_balance_ NUMERIC(30, 6);
    new_order orders;
BEGIN
    -- Fetch the user's total balance once
    SELECT total_balance INTO total_balance_
    FROM users
    WHERE userid = user_id;

    -- Validate user's total balance and leverage
    IF total_balance_ < collateral THEN
        RAISE EXCEPTION 'Insufficient balance to create order. Required: %, Available: %', collateral, total_balance;
    END IF;

    IF leverage > 50 AND total_balance_ < 30000 THEN
        RAISE EXCEPTION 'Leverage above 50x is not allowed unless total balance exceeds $30,000.';
    END IF;

    IF leverage > 1250 THEN
        RAISE EXCEPTION 'Leverage cannot exceed 1250x.';
    END IF;

    -- Insert the order
    INSERT INTO orders (
        userid, order_type, leverage, pair, collateral, entry_price, liq_price
    ) VALUES (
        user_id, order_type, leverage, pair, collateral, entry_price, liq_price
    )
    RETURNING * INTO new_order;

    -- Return the inserted order
    RETURN new_order;

END;
$$ LANGUAGE plpgsql;


