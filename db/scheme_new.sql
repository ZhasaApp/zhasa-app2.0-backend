CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    phone      VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(255)        NOT NULL,
    last_name  VARCHAR(255)        NOT NULL,
    created_at TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_active  BOOLEAN             NOT NULL DEFAULT TRUE
);

CREATE TABLE branches
(
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE brands
(
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) UNIQUE NOT NULL,
    description TEXT                NOT NULL,
    created_at  TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE value_type AS ENUM ('amount', 'count');

CREATE TABLE sale_types
(
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    color       VARCHAR(255) NOT NULL,
    gravity     INTEGER      NOT NULL,
    value_type  value_type   NOT NULL DEFAULT ('count')
);

CREATE TABLE roles
(
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) UNIQUE NOT NULL,
    description TEXT                NOT NULL,
    created_at  TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_roles
(
    id      SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users (id) NOT NULL,
    role_id INTEGER REFERENCES roles (id) NOT NULL,
    UNIQUE (user_id, role_id)
);

CREATE TABLE branch_users_roles(
    id SERIAL PRIMARY KEY,
    user_role_id INTEGER REFERENCES user_roles(id),
    branch_id INTEGER REFERENCES branches(id),
    UNIQUE (user_role_id, branch_id)
);

CREATE TABLE sales
(
    id           SERIAL PRIMARY KEY,
    user_id      INTEGER REFERENCES users (id) ON DELETE CASCADE      NOT NULL,
    sale_date    TIMESTAMP                                            NOT NULL,
    amount       BIGINT                                               NOT NULL,
    sale_type_id INTEGER REFERENCES sale_types (id) ON DELETE CASCADE NOT NULL,
    description  TEXT                                                 NOT NULL,
    created_at   TIMESTAMP                                            NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sales_brands
(
    sale_id  INTEGER REFERENCES sales (id)  NOT NULL,
    brand_id INTEGER REFERENCES brands (id) NOT NULL,
    UNIQUE (sale_id, brand_id)
);

CREATE TABLE user_brands
(
    id       SERIAL PRIMARY KEY,
    user_id  INTEGER REFERENCES users (id)  NOT NULL,
    brand_id INTEGER REFERENCES brands (id) NOT NULL,
    UNIQUE (user_id, brand_id)
);

CREATE TABLE user_brand_sale_type_goals
(
    id SERIAL PRIMARY KEY,
    user_brand   INTEGER REFERENCES users_brands (id),
    sale_type_id INTEGER REFERENCES sale_types (id),
    value        BIGINT    NOT NULL,
    from_date    TIMESTAMP NOT NULL,
    to_date      TIMESTAMP NOT NULL,
    UNIQUE (user_brand, sale_type_id, from_date, to_date)
);

CREATE TABLE branch_brands
(
    id        SERIAL PRIMARY KEY,
    branch_id INTEGER REFERENCES branches (id),
    brand_id  INTEGER REFERENCES brands (id),
    UNIQUE (brand_id, brand_id)
);

CREATE TABLE branch_brand_sale_type_goals
(
    id           SERIAL PRIMARY KEY,
    branch_brand INTEGER REFERENCES branch_brands (id),
    sale_type_id INTEGER REFERENCES sale_types (id),
    value        BIGINT NOT NULL,
    from_date    TIMESTAMP NOT NULL,
    to_date      TIMESTAMP NOT NULL,
    UNIQUE (branch_brand, sale_type_id, from_date, to_date)
);
