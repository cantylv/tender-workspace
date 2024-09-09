CREATE TABLE IF NOT EXISTS employee (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE IF NOT EXISTS organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);

CREATE TABLE IF NOT EXISTS organization (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS organization_responsible (
    id SERIAL PRIMARY KEY,
    organization_id INT REFERENCES organization(id) ON DELETE CASCADE,
    user_id INT REFERENCES employee(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tender (
    id SERIAL PRIMARY KEY
    organization_client_id INT REFERENCES organization(id) NOT NULL ON DELETE CASCADE,
    performer_id INT REFERENCES employee(id),
    organization_performer_id INT REFERENCES organization(id)
);

CREATE TABLE IF NOT EXISTS tender_version (
    id SERIAL PRIMARY KEY,
    organization_id INT REFERENCES organization(id) ON DELETE CASCADE,
    task TEXT NOT NULL,
    price INT NOT NULL,
    version INT NOT NULL
);

CREATE TABLE IF NOT EXISTS offer (
    id SERIAL PRIMARY KEY,
    tender_id INT REFERENCES tender(id),
    performer_id INT REFERENCES employee(id),
    offer_version_id INT REFERENCES offer_version(id)
);

CREATE TABLE IF NOT EXISTS offer_version (
    id SERIAL PRIMARY KEY,
    organization_id INT REFERENCES organization(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    price INT NOT NULL,
    version INT NOT NULL
);