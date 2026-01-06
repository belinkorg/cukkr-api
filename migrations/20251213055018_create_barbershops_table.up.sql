CREATE TABLE IF NOT EXISTS barbershops
(
    id           VARCHAR(50) PRIMARY KEY,
    slug         VARCHAR(50)  NOT NULL UNIQUE,
    name         VARCHAR(255) NOT NULL,
    email        VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(20) UNIQUE,
    description  TEXT,
    address      TEXT,
    is_active    BOOLEAN   DEFAULT TRUE,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP
    );

-- Create indexes
CREATE INDEX idx_barbershops_deleted_at ON barbershops (deleted_at);
CREATE UNIQUE INDEX idx_barbershops_slug ON barbershops (slug);
CREATE UNIQUE INDEX idx_barbershops_email ON barbershops (email);
CREATE UNIQUE INDEX idx_barbershops_phone_number ON barbershops (phone_number) WHERE phone_number IS NOT NULL;