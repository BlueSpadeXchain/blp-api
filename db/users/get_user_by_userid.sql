CREATE OR REPLACE FUNCTION get_user_by_userid(user_id VARCHAR)
RETURNS TABLE(
    id UUID,
    userid VARCHAR,
    wallet_address VARCHAR,
    wallet_type VARCHAR,
    balance NUMERIC(30, 6),
    perp_balance NUMERIC(30, 6),
    escrow_balance NUMERIC(30, 6),
    stake_balance NUMERIC(30, 6),
    frozen_balance NUMERIC(30, 6),
    total_balance NUMERIC(30, 6),
    created_at TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
    SELECT * FROM users WHERE users.userid = user_id;
END;
$$ LANGUAGE plpgsql;