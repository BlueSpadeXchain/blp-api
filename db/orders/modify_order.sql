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

-- upon modifying the order, the collateral must reflect the new values
-- this means if a new tp or change in collateral we need to deduct winnings

-- reduced leverage does NOT get fees back, increased leverage added difference in fee
-- for now, utilization fee will not effect the modified trade

CREATE OR REPLACE FUNCTION unsigned_modify_order(
    p_order_id UUID,
    p_entry_price NUMERIC,
    p_leverage NUMERIC,
    p_collateral NUMERIC,
    p_stop_price NUMERIC,
    p_liq_price NUMERIC,
    p_max_price NUMERIC,
    p_max_value NUMERIC,
    p_lim_price NUMERIC,
    p_tp_price NUMERIC,
    p_tp_value NUMERIC,
    p_tp_collateral NUMERIC,
    -- p_tp_at TIMESTAMP, not an input but needs to be reset upon new tp
    p_pnl NUMERIC,
    p_open_fee NUMERIC,
    p_close_fee NUMERIC,
    p_prev_modified_at TIMESTAMP
) RETURNS JSON AS $$
DECLARE
    v_order_modification order_modifications;
    v_order orders2;
    v_user users;
    v_signature_id UUID;
    v_signature_hash VARCHAR(64);
    v_expiry_time TIMESTAMP WITH TIME ZONE;
    v_collateral_operation TEXT;
BEGIN
    -- Select target order to sign
    SELECT * INTO v_order
    FROM orders2 WHERE orders2.id = order_id;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'Order with ID % does not exist.', order_id;
    END IF;

    SELECT * INTO v_user
    FROM users WHERE users.userid = v_order.userid;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'User with ID % does not exist.', v_order.userid;
    END IF;

    -- Backend must pre-execute this
    IF tp_value = 0 AND tp_collateral != NULL THEN
        v_order.collateral := v_order.collateral - v_order.tp_collateral;
        v_order.max_value := v_order.collateral * 10;
    END IF;

    -- Pairty check correct last modified timestamp
    IF p_prev_modified_at != v_order.modified_at THEN
        RAISE EXCEPTION 'Invalid modify timestamp: expected %, found %', p_prev_modified_at, v_order.modified_at;
    END IF;

    -- invalid order
    IF v_order.ended_at != NULL THEN
        RAISE EXCEPTION 'Invalid order state: ended at %', v_order.ended_at;
    END IF;

    IF p_open_fee < v_order.open_fee THEN 
        RAISE EXCEPTION 'Invalid open fee, open fee cannot be reduced';
    END IF;
    
    IF p_close_fee < v_order.close_fee THEN 
        RAISE EXCEPTION 'Invalid close fee, close fee cannot be reduced';
    END IF;

    IF p_lim_price IS NOT NULL AND v_order.status NOT IN ('limit', 'unsigned') THEN
        RAISE EXCEPTION 'Cannot modify limit price for an order with status %', v_order.status;
    END IF;


    IF p_leverage <= 1.1 OR p_leverage > 1250 THEN
        RAISE EXCEPTION 'Invalid leverage, required: 1.1 < value <= 1250';
    END IF;

    IF p_leverage > v_order.leverage THEN
        IF p_open_fee <= v_order.open_fee THEN 
            RAISE EXCEPTION 'Invalid open fee, increased leverage requires increased open fee';
        END IF;
    ELSIF p_leverage < v_order.leverage THEN
        IF p_close_fee <= v_order.close_fee THEN 
            RAISE EXCEPTION 'Invalid close fee, decreased leverage requires decreased close fee';
        END IF;
    END IF;

    -- when leverage is unchanged
    IF p_tp_collateral >= p_collateral THEN
        RAISE EXCEPTION 'Invalid take profit collateral, collateral % must exceed take profit collateral %', p_collateral, p_tp_collateral;
    END IF;

    IF p_tp_value >= v_order.max_value THEN
        RAISE EXCEPTION 'Invalid take profit value, max value % must exceed take profit value %', v_order.max_value, p_tp_value;
    END IF;


    IF v_order.collateral < p_collateral THEN
        IF p_max_value < v_order.max_value THEN
            RAISE EXCEPTION 'Invalid max price, when adding collateral max price must be increased';
        END IF;
    ELSIF v_order.collateral > p_collateral THEN
        IF p_max_value > v_order.max_value THEN
            RAISE EXCEPTION 'Invalid max price, when removing collateral max price must be reduced';
        END IF;
    ELSE -- v_order.collateral = p_collateral
        IF p_max_price != v_order.max_price THEN
            RAISE EXCEPTION 'Invalid max price, when static collateral max price cannot change';
        END IF;

        IF p_max_value != v_order.max_value THEN
            RAISE EXCEPTION 'Invalid max price, when static collateral max price cannot change';
        END IF;
    END IF;

    IF v_order.order_type = 'long' THEN
        IF p_stop_price <= p_liq_price THEN
            RAISE EXCEPTION 'Stop price cannot exceed liquidation price';
        END IF;

        IF p_tp_price >= v_order.max_price THEN
            RAISE EXCEPTION 'Invalid take profit price, max price % must exceed take profit price %', v_order.max_price, p_tp_price;
        END IF;

        -- p_max_price
        -- p_max_value
        IF v_order.collateral < p_collateral THEN
            IF p_max_price > v_order.max_price THEN
                RAISE EXCEPTION 'Invalid long max price, when adding collateral max price must be increased';
            END IF;
        ELSIF v_order.collateral > p_collateral THEN
            IF p_max_price < v_order.max_price THEN
                RAISE EXCEPTION 'Invalid long max price, when removing collateral max price must be reduced';
            END IF;
        END IF;
    ELSIF v_order.order_type = 'short' THEN
        IF p_stop_price >= p_liq_price THEN
            RAISE EXCEPTION 'Liquidiation price cannot exceed stop price';
        END IF;

        IF p_tp_price <= v_order.max_price THEN
            RAISE EXCEPTION 'Invalid take profit price, take profit price % must exceed max price %', p_tp_price, v_order.max_price;
        END IF;

        -- p_max_price
        -- p_max_value
        IF v_order.collateral < p_collateral THEN
            IF p_max_price < v_order.max_price THEN
                RAISE EXCEPTION 'Invalid short max price, when adding collateral max price must be reduced';
            END IF;
        ELSIF v_order.collateral > p_collateral THEN
            IF p_max_price > v_order.max_price THEN
                RAISE EXCEPTION 'Invalid short max price, when removing collateral max price must be increased';
            END IF;
        END IF;
    ELSE
        RAISE EXCEPTION 'Unknown order type % found', v_order.order_type;
    END IF;


    -- Insert the order modification using the provided parameters
    INSERT INTO order_modifications (
        orderid,
        userid,
        leverage,
        collateral,
        entry_price,
        liq_price,
        max_price,
        max_value,
        lim_price,
        stop_price,
        tp_price,
        tp_value,
        tp_collateral,
        pnl,
        open_fee,
        close_fee
    )
    VALUES (
        p_order_id,
        v_order.userid,
        p_leverage,
        p_collateral,
        p_entry_price,
        p_liq_price,
        p_max_price,
        p_max_value,
        p_lim_price,
        p_stop_price,
        p_tp_price,
        p_tp_value,
        p_tp_collateral,
        p_pnl,
        p_open_fee,
        p_close_fee
    )
    RETURNING * INTO v_order_modification;

    -- Create signature verification request
    SELECT signature_id, signature_hash, expiry_time
    INTO v_signature_id, v_signature_hash, v_expiry_time
    FROM generate_signature_hash(v_user.wallet_address, v_user.wallet_type, 'order_modification', v_order_modification.id, 'cancel');
    

    -- Return the result as JSON
    RETURN json_build_object(
        'order_modification_id', v_order_modification.id,
        'signature_id', v_signature_id,
        'signature_hash', v_signature_hash,
        'expiry_time', v_expiry_time
    );
END;
$$ LANGUAGE plpgsql;

-- need to validate timestap (!closed_at, now > modified_at, created_at > modified_at)
-- need to update_global_state_on_modify function
-- need to apply new tp_at to null when applicable

CREATE OR REPLACE FUNCTION signed_modify_order(
    p_order_modification_id UUID,
    p_signature_id UUID
) RETURNS jsonb AS $$
DECLARE
    total_balance_ NUMERIC(30, 6);
    signed_order orders2;
    v_order orders2;
    proof_ signature_validations;
    v_is_valid BOOLEAN;
    v_error_message TEXT;
BEGIN
    -- select target order to sign
    SELECT o.* INTO v_order FROM orders2 o WHERE o.id = order_id;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'Order with ID % does not exist.', order_id;
    END IF;

    -- check if order can be signed
    IF order_.status != 'unsigned' AND order_.status != 'limit' THEN
        RAISE EXCEPTION 'Only orders with status limit or unsigned can be canceled';
    END IF;

    -- select signature proof
    SELECT * INTO proof_ FROM signature_validations WHERE signature_validations.id = p_signature_id;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'Signature ID % does not exist.', p_signature_id;
    END IF;
    IF proof_.reference_table != 'orders2' OR proof_.reference_id != order_id THEN
        RAISE EXCEPTION 'order id and signature id mismatch';
    END IF;

    -- validate request
    SELECT is_valid, error_message 
    INTO v_is_valid, v_error_message FROM validate_signature(p_signature_id);
    
    -- IF NOT v_is_valid THEN
    --     RAISE EXCEPTION 'Signature validation failed for order %: %', order_id, v_error_message;
    -- END IF;

    

    IF v_is_valid THEN
        -- if a limit order is cancelled we need to return funds to user
        IF order_.status = 'limit' THEN
            UPDATE users
            SET 
                balance = balance + order_.collateral + order_.open_fee,
                escrow_balance = escrow_balance - order_.collateral - order_.open_fee
            WHERE userid = order_.userid;
        END IF;
    
        -- Update the order status to pending
        UPDATE orders2
        SET
            status = 'canceled',
            ended_at = CURRENT_TIMESTAMP
        WHERE orders2.id = order_id;
    END IF;

    UPDATE global_state
    SET value = value + 1,
        updated_at = CURRENT_TIMESTAMP
    WHERE key = 'total_orders_canceled';

    -- Return the updated order
    SELECT * INTO signed_order FROM orders2 WHERE orders2.id = order_id;
    RETURN jsonb_build_object(
        'order', to_jsonb(signed_order),
        'is_valid', v_is_valid,
        'error_massage', v_error_message
    );

END;
$$ LANGUAGE plpgsql;