CREATE TABLE subscriptions (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    service_name VARCHAR(100) NOT NULL,
    price INT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE
);

-- CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions (user_id);
-- CREATE UNIQUE INDEX IF NOT EXISTS idx_user_service ON subscriptions (user_id, service_name);
