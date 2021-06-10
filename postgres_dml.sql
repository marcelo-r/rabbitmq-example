CREATE TABLE IF NOT EXISTS record (
    id SERIAL PRIMARY KEY,
    seq_id VARCHAR(50) UNIQUE NOT NULL,
    datetime TIMESTAMP NOT NULL, 
    email VARCHAR(256),
    ipv4 VARCHAR(15) NOT NULL,
    mac VARCHAR(17) NOT NULL,
    country_code VARCHAR(4),
    user_agent VARCHAR(256) NOT NULL
);
