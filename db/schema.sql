CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    phone      VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(255)        NOT NULL,
    last_name  VARCHAR(255)        NOT NULL
);

CREATE TABLE users_avatars
(
    user_id    INTEGER UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    avatar_url TEXT NOT NULL
);

CREATE TABLE sales_managers
(
    id      SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE sale_types
(
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL
);

CREATE TABLE sales_manager_goals
(
    id               SERIAL PRIMARY KEY,
    from_date        TIMESTAMP NOT NULL,
    to_date          TIMESTAMP NOT NULL,
    amount           BIGINT    NOT NULL,
    sales_manager_id INTEGER REFERENCES sales_managers (id) ON DELETE CASCADE,
    UNIQUE (from_date, to_date, sales_manager_id)
);

CREATE TABLE sales
(
    id               SERIAL PRIMARY KEY,
    sales_manager_id INTEGER REFERENCES sales_managers (id) ON DELETE CASCADE NOT NULL,
    date             TIMESTAMP                                                NOT NULL,
    amount           BIGINT                                                   NOT NULL,
    sale_type_id     INTEGER REFERENCES sale_types (id) ON DELETE CASCADE     NOT NULL,
    description      TEXT                                                     NOT NULL
);

CREATE TABLE branches
(
    id          SERIAL PRIMARY KEY,
    title       TEXT NOT NULL,
    description Text
);

CREATE TABLE branch_sales_managers
(
    sales_manager_id INTEGER REFERENCES sales_managers (id) ON DELETE CASCADE,
    branch_id        INTEGER REFERENCES branches (id) ON DELETE CASCADE,
    UNIQUE (sales_manager_id)
);

CREATE INDEX idx_branch_sales_managers_sales_manager_id ON branch_sales_managers (sales_manager_id);
CREATE INDEX idx_branch_sales_managers_branch_id ON branch_sales_managers (branch_id);

CREATE VIEW user_avatar_view AS
SELECT u.id AS id, u.phone AS phone, u.first_name AS first_name, u.last_name AS last_name, COALESCE(ua.avatar_url, '') AS avatar_url
FROM users u
         LEFT JOIN users_avatars ua ON u.id = ua.user_id;

