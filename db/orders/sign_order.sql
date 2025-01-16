CREATE OR REPLACE FUNCTION sign_order(
    order_id UUID
) RETURNS orders AS $$
DECLARE
    total_balance_ NUMERIC(30, 6);
    signed_order orders;
    order_ orders;
BEGIN
    -- select target order to sign
    SELECT * INTO order_ FROM orders WHERE id = order_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Order with ID % does not exist.', order_id;
    END IF;

    IF order_.status != 'unsigned' THEN
        RAISE EXCEPTION 'Order could not be signed. Status: %', order_.status;
    END IF;

    -- Fetch the user's total balance once
    SELECT total_balance INTO total_balance_
    FROM users
    WHERE userid = order_.userid;

    -- Validate user's total balance and leverage
    IF total_balance_ < order_.collateral THEN
        RAISE EXCEPTION 'Insufficient balance to create order. Required: %, Available: %', order_.collateral, total_balance_;
    END IF;

    IF order_.leverage > 50 AND total_balance_ < 30000 THEN
        RAISE EXCEPTION 'Leverage above 50x is not allowed unless total balance exceeds $30,000.';
    END IF;

    IF order_.leverage > 1250 THEN
        RAISE EXCEPTION 'Leverage cannot exceed 1250x.';
    END IF;

    -- Deduct collateral from user's balance and move it to escrow
    UPDATE users
    SET 
        balance = balance - order_.collateral,
        escrow_balance = escrow_balance + order_.collateral
    WHERE userid = order_.userid;

    -- Update the order status to pending
    UPDATE orders
    SET
        status = 'pending'
    WHERE id = order_id;

    -- Return the updated order
    SELECT * INTO signed_order FROM orders WHERE id = order_id;
    RETURN signed_order;

END;
$$ LANGUAGE plpgsql;
