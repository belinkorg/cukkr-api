CREATE TABLE IF NOT EXISTS user_barbershops
(
    user_id       VARCHAR(50)                         NOT NULL,
    barbershop_id VARCHAR(50)                         NOT NULL,
    role          VARCHAR(30)                         NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, barbershop_id)
    );

-- Create indexes for efficient lookups
CREATE INDEX idx_user_barbershops_user_id ON user_barbershops (user_id);
CREATE INDEX idx_user_barbershops_barbershop_id ON user_barbershops (barbershop_id);
CREATE INDEX idx_user_barbershops_role ON user_barbershops (role);

ALTER TABLE user_barbershops
    ADD CONSTRAINT fk_user_barbershops_user
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

ALTER TABLE user_barbershops
    ADD CONSTRAINT fk_user_barbershops_barbershop
        FOREIGN KEY (barbershop_id) REFERENCES barbershops (id) ON DELETE CASCADE;
