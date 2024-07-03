-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE cats (
                      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                      name VARCHAR(255) NOT NULL,
                      years_of_experience INT NOT NULL,
                      breed VARCHAR(255) NOT NULL,
                      salary DECIMAL(10, 2) NOT NULL
);

CREATE TABLE missions (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          name VARCHAR(255) NOT NULL,
                          cat_id UUID,
                          complete BOOLEAN NOT NULL
);

CREATE TABLE targets (
                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                         mission_id UUID NOT NULL,
                         name VARCHAR(255),
                         country VARCHAR(255),
                         notes TEXT,
                         complete BOOLEAN NOT NULL,
                         FOREIGN KEY (mission_id) REFERENCES missions(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE targets cascade;
DROP TABLE missions cascade;
DROP TABLE cats cascade;
-- +goose StatementEnd