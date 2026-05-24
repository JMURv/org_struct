-- +goose Up
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS department (
    id                BIGINT        PRIMARY KEY,
    name              VARCHAR(200)  NOT NULL CHECK (char_length(trim(name)) > 0),
    parent_id         BIGINT,
    created_at        TIMESTAMPTZ   DEFAULT NOW(),

    CONSTRAINT fk_parent_department FOREIGN KEY (parent_id) REFERENCES department (id),
    CONSTRAINT departments_name_parent_unique UNIQUE(name, parent_id)
);

CREATE INDEX IF NOT EXISTS department_parent_idx ON department (parent_id);

CREATE TABLE IF NOT EXISTS employee (
    id                BIGINT        PRIMARY KEY,
    department_id     BIGINT,
    full_name         VARCHAR(200)  NOT NULL,
    position          VARCHAR(200)  NOT NULL,
    hired_at          TIMESTAMPTZ,
    created_at        TIMESTAMPTZ   DEFAULT NOW(),

    CONSTRAINT fk_department FOREIGN KEY (department_id) REFERENCES department (id)
);


-- +goose Down
SELECT 'down SQL query';

DROP TABLE IF EXISTS department;
DROP INDEX IF EXISTS department_parent_idx;
DROP TABLE IF EXISTS employee;

