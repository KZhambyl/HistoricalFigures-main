# Go_midterm
# Project description
```
Mini site like wiki with short informations about significant figures of Kazakh history: Name, Years of life and description of figure.
```

# Historical Figures REST API
```
GET /v1/healthcheck
POST /v1/figures
GET /v1/figures/:id
PUT /v1/figures/:id
DELETE /v1/figures/:id
```
# DB Structure
```
Table figures {
    id bigserial [primary key]
    created_at timestamp
    name text
    years_of_life text
    description text
    version integer
}
```
```
Table schema_migrations {
    version bigint
    dirty boolean
}
```
```
Table tokens {
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL
}
```
```
Table users {
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
    version integer NOT NULL DEFAULT 1
}
