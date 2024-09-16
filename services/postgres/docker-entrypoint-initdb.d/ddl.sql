CREATE TABLE employee (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);

CREATE TABLE organization (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization_responsible (
    id SERIAL PRIMARY KEY,
    organization_id INT REFERENCES organization(id) ON DELETE CASCADE,
    user_id INT REFERENCES employee(id) ON DELETE CASCADE
);

CREATE TYPE tender_type AS ENUM (
    'Construction',
    'Delivery',
    'Manufacture'
);

CREATE TYPE tender_status AS ENUM (
    'Created',
    'Published',
    'Closed'
);

-- необходимо добавить еще id организации-исполнителя
CREATE TABLE tender (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    type tender_type NOT NULL,
    status tender_status NOT NULL,
    version INT NOT NULL,
    organization_id INT REFERENCES organization(id) ON DELETE CASCADE NOT NULL,
    creator_id INT REFERENCES employee(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE creator_type AS ENUM (
    'User',
    'Responsible'
);

CREATE TYPE bid_status AS ENUM (
    'Created',
    'Published',
    'Canceled',
    'Approved',
    'Rejected'
);

CREATE TABLE bids (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    status bid_status NOT NULL,
    version INT NOT NULL,
    tender_id INT REFERENCES tender(id) ON DELETE CASCADE NOT NULL,
    creator_id INT REFERENCES employee(id) ON DELETE CASCADE NOT NULL,
    author_type creator_type NOT NULL,
    organization_id INT REFERENCES organization(id) ON DELETE CASCADE NULL,  
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE feedback (
    id SERIAL PRIMARY KEY,
    bid_id INT REFERENCES bids(id) ON DELETE CASCADE NOT NULL,
    text TEXT NOT NULL
);