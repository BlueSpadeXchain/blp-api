GRANT EXECUTE ON FUNCTION get_orders_by_userid(VARCHAR) to public;
GRANT EXECUTE ON FUNCTION get_orders_by_address(VARCHAR, VARCHAR) to public;
GRANT EXECUTE ON FUNCTION get_order_by_id(UUID) to public;
GRANT EXECUTE ON FUNCTION sign_order(UUID) to public;
GRANT EXECUTE ON FUNCTION create_order(VARCHAR, VARCHAR, NUMERIC, VARCHAR, NUMERIC, NUMERIC, NUMERIC) to public;