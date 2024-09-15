-- Инициализация БД пользователями/рабочими
INSERT INTO employee (
    username, 
    first_name,
    last_name,
    created_at,
    updated_at
) VALUES (
    'tussan_pussan',
    'Ivan',
    'Lobanov',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
), (
    'igormed',
    'Igor',
    'Medvedev',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
), (
    'futplayer',
    'Miroslav',
    'Kalinin',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
), (
    'niceday55',
    'George',
    'Afanasyev',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);

-- Инициализация БД организациями
INSERT INTO organization (
    name, 
    description,
    type,
    created_at,
    updated_at
) VALUES (
    'T-Bank',
    'T-Bank is a modern financial institution that focuses on providing innovative banking solutions 
    to meet the diverse needs of its customers. Established with the vision of transforming 
    the banking experience, T-Bank leverages technology to offer a wide range of services, including 
    personal and business banking, loans, investment options, and digital banking solutions.',
    'JSC',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
), (
    'Tech Solutions',
    'Tech Solutions is a leading provider of technology consulting and software development services. 
    Our mission is to help businesses leverage technology to improve efficiency and drive growth. 
    We specialize in custom software solutions, IT strategy, and digital transformation.',
    'LLC',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
), (
        'Green Energy Corp',
    'Green Energy Corp is dedicated to providing sustainable energy solutions. 
    We focus on renewable energy sources such as solar and wind power, aiming to reduce carbon footprints 
    and promote environmental sustainability. Our goal is to make clean energy accessible to everyone.',
    'JSC',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);

-- Инициализация БД ответственными за организации

INSERT INTO organization_responsible (
    organization_id,
    user_id
) VALUES (1, 1), (1, 2), (2,3), (3, 3), (3, 1);

