CREATE TABLE branches
(
    id          SERIAL PRIMARY KEY,
    title       TEXT               NOT NULL,
    description Text               NOT NULL,
    branch_key  VARCHAR(16) UNIQUE NOT NULL,
    created_at  TIMESTAMP          NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    phone      VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(255)        NOT NULL,
    last_name  VARCHAR(255)        NOT NULL,
    created_at TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users_codes
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER REFERENCES users (id) NOT NULL,
    code       INTEGER                       NOT NULL,
    created_at TIMESTAMP                     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users_avatars
(
    user_id    INTEGER UNIQUE REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    avatar_url TEXT                                                   NOT NULL
);

CREATE TABLE sales_managers
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER UNIQUE REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    branch_id  INTEGER REFERENCES branches (id) ON DELETE CASCADE     NOT NULL,
    created_at TIMESTAMP                                              NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sale_types
(
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    color       VARCHAR(255) NOT NULL,
    gravity     INTEGER      NOT NULL
);

CREATE TABLE sales_manager_goals_by_types
(
    id               SERIAL PRIMARY KEY,
    from_date        TIMESTAMP                                                NOT NULL,
    to_date          TIMESTAMP                                                NOT NULL,
    amount           BIGINT                                                   NOT NULL,
    sales_manager_id INTEGER REFERENCES sales_managers (id) ON DELETE CASCADE NOT NULL,
    type_id          INTEGER REFERENCES sale_types (id) ON DELETE CASCADE     NOT NULL,
    UNIQUE (from_date, to_date, sales_manager_id, type_id)
);

CREATE TABLE sales_manager_goals_ratio_by_period
(
    from_date TIMESTAMP NOT NULL,
    to_date   TIMESTAMP NOT NULL,
    ratio FLOAT NOT NULL,
    sales_manager_id INTEGER REFERENCES sales_managers (id) ON DELETE CASCADE NOT NULL,
    UNIQUE (from_date, to_date, sales_manager_id)
);

CREATE TABLE branch_goals_by_types
(
    id        SERIAL PRIMARY KEY,
    from_date TIMESTAMP                                          NOT NULL,
    to_date   TIMESTAMP                                          NOT NULL,
    amount    BIGINT                                             NOT NULL,
    branch_id INTEGER REFERENCES branches (id) ON DELETE CASCADE NOT NULL,
    type_id          INTEGER REFERENCES sale_types (id) ON DELETE CASCADE     NOT NULL,
    UNIQUE (from_date, to_date, branch_id, type_id)
);

CREATE TABLE sales
(
    id               SERIAL PRIMARY KEY,
    sales_manager_id INTEGER REFERENCES sales_managers (id) ON DELETE CASCADE NOT NULL,
    sale_date        TIMESTAMP                                                NOT NULL,
    amount           BIGINT                                                   NOT NULL,
    sale_type_id     INTEGER REFERENCES sale_types (id) ON DELETE CASCADE     NOT NULL,
    description      TEXT                                                     NOT NULL,
    created_at       TIMESTAMP                                                NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sales_sales_manager_id
    ON sales (sales_manager_id);

CREATE VIEW user_avatar_view AS
SELECT u.id                        AS id,
       u.phone                     AS phone,
       u.first_name                AS first_name,
       u.last_name                 AS last_name,
       COALESCE(ua.avatar_url, '') AS avatar_url
FROM users u
         LEFT JOIN users_avatars ua ON u.id = ua.user_id;

CREATE VIEW sales_managers_view AS
SELECT u.id         AS user_id,
       u.phone      AS phone,
       u.first_name AS first_name,
       u.last_name  AS last_name,
       u.avatar_url AS avatar_url,
       s.id         as sales_manager_id,
       b.id         as branch_id,
       b.title      as branch_title
FROM user_avatar_view u
         JOIN sales_managers s ON u.id = s.user_id
         JOIN branches b ON s.branch_id = b.id;

CREATE TABLE branch_directors
(
    id        SERIAL PRIMARY KEY,
    user_id   INTEGER UNIQUE REFERENCES users (id) ON DELETE CASCADE    NOT NULL,
    branch_id INTEGER UNIQUE REFERENCES branches (id) ON DELETE CASCADE NOT NULL
);

CREATE VIEW branch_directors_view AS
SELECT u.id         AS user_id,
       u.phone      AS phone,
       u.first_name AS first_name,
       u.last_name  AS last_name,
       u.avatar_url AS avatar_url,
       bd.id        as branch_director_id,
       b.id         as branch_id,
       b.title      as branch_title
FROM user_avatar_view u
         JOIN branch_directors bd ON u.id = bd.user_id
         JOIN branches b ON bd.branch_id = b.id;


CREATE TABLE branches_goals_ratio_by_period
(
    from_date TIMESTAMP NOT NULL,
    to_date   TIMESTAMP NOT NULL,
    ratio FLOAT NOT NULL,
    branch_id INTEGER REFERENCES branches (id) ON DELETE CASCADE NOT NULL,
    UNIQUE (from_date, to_date, branch_id)
);

