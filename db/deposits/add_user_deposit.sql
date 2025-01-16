CREATE OR REPLACE FUNCTION add_user_deposit(
    wallet_addr VARCHAR,
    wallet_t VARCHAR,
    chain TEXT,
    blk TEXT,
    blk_hash VARCHAR,
    tx_hash VARCHAR,
    sndr VARCHAR,
    deposit_nonce TEXT,
    asset_addr VARCHAR,
    amt TEXT,
    val NUMERIC(78, 9)
) RETURNS VOID AS $$
DECLARE
    user_data RECORD;
BEGIN
    -- Use get_or_create_user to retrieve or create the user
    SELECT * INTO user_data 
    FROM get_or_create_user(wallet_addr, wallet_t);

    -- Update the user's balance
    UPDATE users
    SET balance = balance + val
    WHERE id = user_data.id;

    -- Add the deposit record
    INSERT INTO deposits (
        userid, wallet_address, wallet_type, chain_id, block, block_hash, tx_hash,
        sender, deposit_nonce, asset, amount, value
    ) VALUES (
        user_data.userid, wallet_addr, wallet_t, chain, blk, blk_hash, tx_hash,
        sndr, deposit_nonce, asset_addr, amt, val
    );
END;
$$ LANGUAGE plpgsql;
