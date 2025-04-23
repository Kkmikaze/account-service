CREATE TABLE accounts
(
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at     TIMESTAMP WITH TIME ZONE DEFAULT now(),
    deleted_at     TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    id             UUID PRIMARY KEY         DEFAULT uuid_generate_v4() NOT NULL,
    citizen_id     VARCHAR(16)                                        NOT NULL,
    name           VARCHAR(150)                                        NOT NULL,
    phone          VARCHAR(16)                                         NOT NULL,
    account_number VARCHAR(14)                                         NOT NULL
);

-- Basic indexes
CREATE INDEX idx_accounts_citizen_id ON accounts (citizen_id);
CREATE INDEX idx_accounts_phone ON accounts (phone);
CREATE INDEX idx_accounts_deleted_at ON accounts (deleted_at);

-- Composite indexes for filtered queries
CREATE INDEX idx_accounts_account_number_not_deleted ON accounts (account_number, deleted_at);

-- Unique indexes
CREATE UNIQUE INDEX uniq_accounts_account_number_not_deleted ON accounts (account_number) WHERE deleted_at IS NULL;