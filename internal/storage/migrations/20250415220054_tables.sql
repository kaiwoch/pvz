-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role_name role NOT NULL
);

CREATE TABLE pvz (
    pvz_id UUID PRIMARY KEY,
    registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    city_name city NOT NULL,
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE reception (
    reception_id UUID PRIMARY KEY,
    date_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    pvz_id UUID NOT NULL,
    status_name status NOT NULL,
    FOREIGN KEY (pvz_id) REFERENCES pvz(pvz_id)
);

CREATE TABLE product (
    product_id UUID PRIMARY KEY,
    date_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    type_name product_type NOT NULL,
    reception_id UUID NOT NULL,
    FOREIGN KEY (reception_id) REFERENCES reception(reception_id)
);

INSERT INTO users (user_id, email, password_hash, role_name) VALUES 
    ('00000000-0000-0000-0000-000000000000','dummy_moderator@test.com', '$2a$10$R9Sbv7JPKcjNkcLSTLGhoueS9nP5ieZOkBTWFt0peEMibfRXLIk62', 'moderator'),
    ('00000000-0000-0000-0000-000000000001','dummy_employee@test.com', '$2a$10$ozMElrcQdyqolglq7KZKNOE3re8HtZXIbzfX2cpdq2Cog9ZhzoGDG', 'employee');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product;
DROP TABLE IF EXISTS reception;
DROP TABLE IF EXISTS pvz;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
