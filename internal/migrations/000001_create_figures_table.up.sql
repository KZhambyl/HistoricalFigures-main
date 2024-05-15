CREATE TABLE IF NOT EXISTS figures (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    years_of_life text NOT NULL,
    description text NOT NULL,
    version integer NOT NULL DEFAULT 1
);