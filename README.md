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

GET /v1/healthcheck
POST /v1/categories
GET /v1/categories/:id
PUT /v1/categories/:id
DELETE /v1/categories/:id

GET /v1/categories/:id/figures
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
Table categories {
    id bigserial [primary key]
    created_at timestamp
    name text
    version integer
}
```
```
Table figures_categories {
    figure_id bigint NOT NULL REFERENCES figures
    category_id bigint NOT NULL REFERENCES categories
    PRIMARY KEY (figure_id, category_id)
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
    hash bytea [primary key]
    user_id bigint
    expiry timestamp
    scope text
}
```
```
Table users {
    id bigserial [primary key]
    created_at timestamp
    name text
    email citext
    password_hash bytea
    activated bool
    version integer
}
