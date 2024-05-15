CREATE TABLE IF NOT EXISTS categories (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    version integer NOT NULL DEFAULT 1
);
CREATE TABLE IF NOT EXISTS figures_categories (
    figure_id bigint NOT NULL REFERENCES figures ON DELETE CASCADE,
    category_id bigint NOT NULL REFERENCES categories ON DELETE CASCADE,
    PRIMARY KEY (figure_id, category_id)
);

INSERT INTO categories (name)
VALUES
    ('khan'),
    ('batyr'),
    ('writer'),
    ('poet'),
    ('thinker');