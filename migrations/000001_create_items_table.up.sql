CREATE TABLE IF NOT EXISTS companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('PRODUCT', 'SERVICE')),
    price DECIMAL(19,4) NOT NULL CHECK (price >= 0),
    category_name VARCHAR(100),
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE', 'ARCHIVED')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT uc_company_code UNIQUE (company_id, code)
);

CREATE INDEX idx_items_company_id ON items(company_id);
CREATE INDEX idx_items_status ON items(status);
CREATE INDEX idx_items_type ON items(type);

-- Seed company for testing
INSERT INTO companies (id, name) VALUES ('c0a80121-7ac0-11d1-898c-00c04fd8d5c1', 'Test Company A') ON CONFLICT DO NOTHING;
