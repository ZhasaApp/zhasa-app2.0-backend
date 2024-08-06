CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    phone      VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(255)        NOT NULL,
    last_name  VARCHAR(255)        NOT NULL,
    created_at TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    password   VARCHAR(255)
);

CREATE TABLE disabled_users
(
    user_id     INTEGER REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    disabled_at TIMESTAMP                                       NOT NULL DEFAULT CURRENT_TIMESTAMP
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

CREATE TABLE departments
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
    key         VARCHAR(255) UNIQUE NOT NULL,
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

CREATE TABLE branch_users
(
    id        SERIAL PRIMARY KEY,
    user_id   INTEGER REFERENCES users (id)    NOT NULL,
    branch_id INTEGER REFERENCES branches (id) NOT NULL,
    UNIQUE (user_id, branch_id)
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
    sale_id  INTEGER REFERENCES sales (id) ON DELETE CASCADE  NOT NULL,
    brand_id INTEGER REFERENCES brands (id) ON DELETE CASCADE NOT NULL,
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
    id           SERIAL PRIMARY KEY,
    user_id      INTEGER REFERENCES users (id)      NOT NULL,
    brand_id     INTEGER REFERENCES brands (id)     NOT NULL,
    sale_type_id INTEGER REFERENCES sale_types (id) NOT NULL,
    value        BIGINT                             NOT NULL,
    from_date    TIMESTAMP                          NOT NULL,
    to_date      TIMESTAMP                          NOT NULL,
    UNIQUE (user_id, brand_id, sale_type_id, from_date, to_date)
);

CREATE TABLE user_brand_ratio
(
    user_id   INTEGER REFERENCES users (id)  NOT NULL,
    brand_id  INTEGER REFERENCES brands (id) NOT NULL,
    ratio     REAL                           NOT NULL,
    from_date TIMESTAMP                      NOT NULL,
    to_date   TIMESTAMP                      NOT NULL,
    UNIQUE (user_id, brand_id, from_date, to_date)
);

CREATE TABLE branch_brands
(
    id        SERIAL PRIMARY KEY,
    branch_id INTEGER REFERENCES branches (id) NOT NULL,
    brand_id  INTEGER REFERENCES brands (id)   NOT NULL,
    UNIQUE (branch_id, brand_id)
);

CREATE TABLE branch_brand_sale_type_goals
(
    id           SERIAL PRIMARY KEY,
    branch_id    INTEGER REFERENCES branches (id)   NOT NULL,
    brand_id     INTEGER REFERENCES brands (id)     NOT NULL,
    sale_type_id INTEGER REFERENCES sale_types (id) NOT NULL,
    value        BIGINT                             NOT NULL,
    from_date    TIMESTAMP                          NOT NULL,
    to_date      TIMESTAMP                          NOT NULL,
    UNIQUE (branch_id, brand_id, sale_type_id, from_date, to_date)
);

CREATE TABLE brand_overall_sale_type_goals
(
    id           SERIAL PRIMARY KEY,
    brand_id     INTEGER REFERENCES brands (id)     NOT NULL,
    sale_type_id INTEGER REFERENCES sale_types (id) NOT NULL,
    value        BIGINT                             NOT NULL,
    from_date    TIMESTAMP                          NOT NULL,
    to_date      TIMESTAMP                          NOT NULL,
    UNIQUE (brand_id, sale_type_id, from_date, to_date)
);

CREATE TABLE posts
(
    id         SERIAL PRIMARY KEY,
    title      VARCHAR(256)                                NOT NULL,
    body       TEXT                                        NOT NULL,
    user_id    INT REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE                    NOT NULL DEFAULT NOW()
);

CREATE TABLE comments
(
    id         SERIAL PRIMARY KEY,
    body       TEXT                                        NOT NULL,
    user_id    INT REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    post_id    INT REFERENCES posts (id) ON DELETE CASCADE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE                    NOT NULL DEFAULT NOW()
);

CREATE TABLE post_images
(
    id        SERIAL PRIMARY KEY,
    image_url TEXT                                        NOT NULL,
    post_id   INT REFERENCES posts (id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE likes
(
    user_id INT REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    post_id INT REFERENCES posts (id) ON DELETE CASCADE NOT NULL,
    PRIMARY KEY (user_id, post_id)
);

CREATE TABLE users_avatars
(
    user_id    INTEGER UNIQUE REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    avatar_url TEXT                                                   NOT NULL
);

CREATE TABLE users_codes
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER REFERENCES users (id) NOT NULL,
    code       INTEGER                       NOT NULL,
    created_at TIMESTAMP                     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE VIEW user_avatar_view AS
SELECT u.id                        AS id,
       u.phone                     AS phone,
       u.first_name                AS first_name,
       u.last_name                 AS last_name,
       COALESCE(ua.avatar_url, '') AS avatar_url
FROM users u
         LEFT JOIN users_avatars ua ON u.id = ua.user_id;

CREATE TABLE branch_brand_users
(
    id              SERIAL PRIMARY KEY,
    branch_brand_id INTEGER REFERENCES branch_brands (id) NOT NULL,
    user_id         INTEGER REFERENCES users (id)         NOT NULL,
    UNIQUE (branch_brand_id, user_id)
);

CREATE TABLE models
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE brand_models
(
    id       SERIAL PRIMARY KEY,
    brand_id INTEGER REFERENCES brands (id) NOT NULL,
    model_id INTEGER REFERENCES models (id) NOT NULL,
    UNIQUE (brand_id, model_id)
);

CREATE TABLE awards
(
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL,
    icon_url    TEXT         NOT NULL
);

CREATE TABLE user_awards
(
    id            SERIAL PRIMARY KEY,
    user_id       INTEGER REFERENCES users (id)  NOT NULL,
    award_id      INTEGER REFERENCES awards (id) NOT NULL,
    award_details JSONB                          NOT NULL, -- Stores period and scope details
    created_at    TIMESTAMP                      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, award_id, award_details)
);
