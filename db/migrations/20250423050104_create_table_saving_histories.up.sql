CREATE TABLE saving_histories
(
    id          UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT now(),
    account_id  UUID           NOT NULL REFERENCES accounts (id) ON DELETE CASCADE,
    amount      NUMERIC(20, 2) NOT NULL CHECK (amount > 0),
    CONSTRAINT fk_account FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE
);
