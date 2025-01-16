CREATE OR REPLACE FUNCTION get_or_create_user(wallet_addr VARCHAR, wallet_t VARCHAR)
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
    -- First, check if the user exists
    RETURN QUERY 
    SELECT * 
    FROM users 
    WHERE users.wallet_address = wallet_addr AND users.wallet_type = wallet_t;

    -- If no user exists, create one
    IF NOT FOUND THEN
        RETURN QUERY 
        INSERT INTO users (wallet_address, wallet_type)
        VALUES (wallet_addr, wallet_t)
        RETURNING *;
    END IF;
END;
$$ LANGUAGE plpgsql;