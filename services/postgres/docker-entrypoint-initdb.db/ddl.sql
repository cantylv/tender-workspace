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

CREATE TYPE IF NOT EXISTS tender_type AS ENUM (
    'Construction',
    'Delivery',
    'Manufacture'
);

CREATE TYPE IF NOT EXISTS tender_status AS ENUM (
    'CREATED',
    'PUBLISHED',
    'CLOSED'
);

-- необходимо добавить еще id организации-исполнителя
CREATE TABLE IF NOT EXISTS tender (
    id SERIAL PRIMARY KEY,
    organization_id INT REFERENCES organization(id) ON DELETE CASCADE NOT NULL,
    creator_username INT NOT NULL,
    service_type tender_type NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    status tender_status NOT NULL,
    version INT NOT NULL
);

CREATE TABLE IF NOT EXISTS offer (
    id SERIAL PRIMARY KEY,
    tender_id INT REFERENCES tender(id) ON DELETE CASCADE NOT NULL,
    performer_id INT REFERENCES employee(id) ON DELETE CASCADE NOT NULL,
    message TEXT NOT NULL,
    version INT NOT NULL
);