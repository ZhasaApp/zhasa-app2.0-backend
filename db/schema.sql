CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    email      VARCHAR(255) UNIQUE NOT NULL,
    password   TEXT                NOT NULL,
    first_name VARCHAR(255),
    last_name  VARCHAR(255),
    avatar_url TEXT
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
    description TEXT
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
    UNIQUE (sales_manager_id, date, sale_type_id)
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

CREATE VIEW ranked_sales_managers AS
WITH sales_summary AS (SELECT sm.id         AS sales_manager_id,
                              SUM(s.amount) AS total_sales_amount,
                              u.first_name  AS first_name,
                              u.last_name   AS last_name,
                              u.avatar_url  AS avatar_url
                       FROM sales s
                                INNER JOIN sales_managers sm ON s.sales_manager_id = sm.id
                                INNER JOIN users u ON sm.user_id = u.id
                       GROUP BY sm.id),
     goal_summary AS (SELECT sm.id     AS sales_manager_id,
                             sg.from_date,
                             sg.to_date,
                             sg.amount AS goal_amount
                      FROM sales_manager_goals sg
                               INNER JOIN sales_managers sm ON s.sales_manager_id = sm.id)
SELECT ss.sales_manager_id,
       ss.first_name,
       ss.last_name,
       ss.avatar_url,
       COALESCE(ss.total_sales_amount / NULLIF(smg.goal_amount, 0), 0) ::float AS ratio, RANK() OVER (ORDER BY COALESCE(ss.total_sales_amount / NULLIF(smg.goal_amount, 0), 0) ::float DESC) AS position
FROM sales_summary ss
LEFT JOIN goal_summary smg ON ss.sales_manager_id = smg.sales_manager_id;

