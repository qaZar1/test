CREATE SCHEMA IF NOT EXISTS wallets;

CREATE TABLE IF NOT EXISTS wallets.balances (
    wallet_id UUID PRIMARY KEY,
    amount NUMERIC(19, 2) NOT NULL
);

INSERT INTO wallets.balances (wallet_id, amount)
VALUES (gen_random_uuid(), 100.00);
