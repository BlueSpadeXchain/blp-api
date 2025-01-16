CREATE OR REPLACE FUNCTION get_order_by_id(id UUID)
RETURNS orders AS $$
DECLARE
    order_ orders;
    user_ users;
BEGIN
    SELECT * INTO order_
    FROM orders
    WHERE orders.id = id;

    SELECT * INTO user_
    FROM users
    WHERE users.userid = order_.user_id;

    RETURN json_build_object(
        'order', row_to_json(order_),
        'user', row_to_json(user_)
    );
END;
$$ LANGUAGE plpgsql;



CREATE OR REPLACE FUNCTION get_orders_by_userid(user_id VARCHAR)
RETURNS SETOF orders AS $$
BEGIN
    RETURN QUERY
    SELECT * FROM deposits WHERE deposits.userid = user_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_orders_by_address(wallet_addr VARCHAR, wallet_t VARCHAR)
RETURNS SETOF orders AS $$
DECLARE
    user_data RECORD;
BEGIN
    -- Use get_or_create_user to retrieve or create the user
    SELECT * INTO user_data 
    FROM get_or_create_user(wallet_addr, wallet_t);

    -- Return the orders associated with the user
    RETURN QUERY 
    SELECT * FROM get_orders_by_userid(user_data.userid);
END;
$$ LANGUAGE plpgsql;